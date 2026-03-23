package matcher

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"k8doctor/internal/cleaner"
	"k8doctor/internal/knowledge"
)

const (
	SimilarityThreshold = 0.45
	VectorDim           = 384
	ModelName           = "all-MiniLM-L6-v2.onnx"
	MaxLength           = 256
)

// MatchResult 匹配结果
type MatchResult struct {
	Entry     *knowledge.KBEntry
	Score     float32
	MatchType string
}

// Matcher 向量匹配器
type Matcher struct {
	kb            *knowledge.KnowledgeBase
	indexPath     string
	modelPath     string
	vectors       map[string][]float32
	mu            sync.RWMutex
	useONNX       bool
	pythonPath    string
	scriptPath    string
	tokenizer     *bertTokenizer
}

// inferenceResult ONNX 推理结果
type inferenceResult struct {
	Embedding []float32 `json:"embedding"`
	Success   bool      `json:"success"`
}

func New(kb *knowledge.KnowledgeBase) (*Matcher, error) {
	scriptDir := filepath.Join(filepath.Dir(kb.Path), "scripts")
	modelsDir := filepath.Join(filepath.Dir(kb.Path), "embeddings", "models")
	modelPath := filepath.Join(modelsDir, ModelName)
	indexPath := filepath.Join(filepath.Dir(kb.Path), "embeddings", "index.json")

	m := &Matcher{
		kb:         kb,
		indexPath:  indexPath,
		modelPath:  modelPath,
		vectors:    make(map[string][]float32),
		scriptPath: filepath.Join(scriptDir, "onnx_inference.py"),
		tokenizer:  newBertTokenizer(),
	}

	// 初始化分词器
	m.tokenizer = newBertTokenizer()

	// 确保 Python 和模型可用
	m.pythonPath = m.findPython()
	if m.pythonPath != "" && m.checkModelExists() {
		m.useONNX = true
		fmt.Println("✓ 检测到 Python 环境，使用 all-MiniLM-L6-v2 ONNX 模型")
	} else {
		m.useONNX = false
		if m.pythonPath == "" {
			fmt.Println("⚠ 未找到 Python，使用 TF-IDF Fallback 向量引擎")
		} else {
			fmt.Println("⚠ 模型文件不存在，使用 TF-IDF Fallback 向量引擎")
		}
	}

	// 加载或构建索引
	if err := m.loadIndex(); err != nil {
		if err := m.buildIndex(); err != nil {
			return nil, fmt.Errorf("构建向量索引失败: %w", err)
		}
	}

	return m, nil
}

func (m *Matcher) findPython() string {
	// 查找 Python 解释器
	pythons := []string{"python", "python3", "py"}
	if runtime.GOOS == "windows" {
		pythons = []string{"python", "python3", "py", "python.exe", "python3.exe"}
	}

	for _, py := range pythons {
		if path, err := exec.LookPath(py); err == nil {
			return path
		}
	}
	return ""
}

func (m *Matcher) checkModelExists() bool {
	_, err := os.Stat(m.modelPath)
	return err == nil
}

func (m *Matcher) ensureModel() error {
	if m.checkModelExists() {
		return nil
	}

	// 确保目录存在
	modelsDir := filepath.Dir(m.modelPath)
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return err
	}

	// 使用 Python 脚本下载（更可靠）
	fmt.Println("首次运行，正在下载 all-MiniLM-L6-v2 模型 (~90MB)...")
	fmt.Println("这只需要一次，之后会缓存")

	cmd := exec.Command(m.pythonPath, filepath.Join(m.scriptPath, "..", "download_model.py"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("下载模型失败: %w", err)
	}

	if !m.checkModelExists() {
		return fmt.Errorf("模型文件仍未创建")
	}

	return nil
}

func (m *Matcher) getEmbedding(text string) ([]float32, error) {
	if m.useONNX {
		return m.getONNXEmbedding(text)
	}
	return m.getTFIDFEmbedding(text)
}

func (m *Matcher) getONNXEmbedding(text string) ([]float32, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 调用 Python 脚本
	cmd := exec.Command(m.pythonPath, m.scriptPath, text)
	output, err := cmd.Output()
	if err != nil {
		// ONNX 失败，降级到 TF-IDF
		return m.getTFIDFEmbedding(text)
	}

	var result inferenceResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	if !result.Success || len(result.Embedding) == 0 {
		return m.getTFIDFEmbedding(text)
	}

	return result.Embedding, nil
}

func (m *Matcher) getTFIDFEmbedding(text string) ([]float32, error) {
	// 使用 TF-IDF 生成伪向量
	vec := make([]float32, VectorDim)

	words := strings.Fields(strings.ToLower(text))
	for i, word := range words {
		hash := hashString(word)
		for j := 0; j < VectorDim; j++ {
			vec[j] += float32((hash >> uint(j)) & 1)
		}
		vec[hash%VectorDim] += float32(1) / float32(len(words)+1)
	}

	// 归一化
	norm := float32(0)
	for _, v := range vec {
		norm += v * v
	}
	norm = float32(math.Sqrt(float64(norm)))
	if norm > 0 {
		for i := range vec {
			vec[i] /= norm
		}
	}

	return vec, nil
}

func (m *Matcher) Match(s *cleaner.Symptom, topK int) ([]*MatchResult, error) {
	results := make([]*MatchResult, 0)

	// 1. 关键词匹配（兜底）
	keywordMatches := m.matchByKeywords(s)
	for _, entry := range keywordMatches {
		results = append(results, &MatchResult{
			Entry:     entry,
			Score:     0.95,
			MatchType: "keyword",
		})
	}

	// 2. Embedding 语义匹配
	queryVec, err := m.getEmbedding(s.Context)
	if err == nil {
		semanticMatches := m.matchByVector(queryVec)
		for _, r := range semanticMatches {
			if !containsEntry(results, r.Entry) {
				results = append(results, r)
			}
		}
	}

	// 3. 排序并返回 Top-K
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > topK {
		results = results[:topK]
	}

	return results, nil
}

func (m *Matcher) matchByKeywords(s *cleaner.Symptom) []*knowledge.KBEntry {
	matched := make([]*knowledge.KBEntry, 0)

	for _, entry := range m.kb.Entries {
		score := float32(0)

		for _, state := range s.States {
			for _, rs := range entry.RelatedStates {
				if stringsContains(state, rs) || stringsContains(rs, state) {
					score += 0.5
				}
			}
		}

		for _, tag := range s.Keywords {
			for _, et := range entry.Tags {
				if stringsContains(tag, et) || stringsContains(et, tag) {
					score += 0.3
				}
			}
		}

		for _, code := range s.ErrorCodes {
			if containsString(entry.RelatedStates, code) {
				score += 0.4
			}
		}

		if score > 0 {
			matched = append(matched, entry)
		}
	}

	return matched
}

func (m *Matcher) matchByVector(queryVec []float32) []*MatchResult {
	results := make([]*MatchResult, 0)

	for _, entry := range m.kb.Entries {
		entryVec, ok := m.vectors[entry.ID]
		if !ok {
			continue
		}

		sim := cosineSimilarity(queryVec, entryVec)
		if sim >= SimilarityThreshold {
			results = append(results, &MatchResult{
				Entry:     entry,
				Score:     sim,
				MatchType: "embedding",
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func (m *Matcher) buildIndex() error {
	fmt.Println("正在构建向量索引...")

	vectors := make(map[string][]float32)

	for i, entry := range m.kb.Entries {
		text := entry.Title + ". " + entry.Content
		vec, err := m.getEmbedding(text)
		if err != nil {
			fmt.Printf("警告: 为条目 %s 生成向量失败: %v\n", entry.ID, err)
			continue
		}
		vectors[entry.ID] = vec
		fmt.Printf("进度: [%d/%d] %s\n", i+1, len(m.kb.Entries), entry.Title)
	}

	m.vectors = vectors
	return m.saveIndex()
}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	dot := float32(0)
	for i := range a {
		dot += a[i] * b[i]
	}
	return dot
}

func (m *Matcher) loadIndex() error {
	data, err := os.ReadFile(m.indexPath)
	if err != nil {
		return err
	}
	var vectors map[string][]float32
	if err := json.Unmarshal(data, &vectors); err != nil {
		return err
	}
	m.vectors = vectors
	return nil
}

func (m *Matcher) saveIndex() error {
	dir := filepath.Dir(m.indexPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(m.vectors)
	if err != nil {
		return err
	}
	return os.WriteFile(m.indexPath, data, 0644)
}

func containsEntry(results []*MatchResult, entry *knowledge.KBEntry) bool {
	for _, r := range results {
		if r.Entry.ID == entry.ID {
			return true
		}
	}
	return false
}

func stringsContains(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func hashString(s string) uint64 {
	h := uint64(5381)
	for _, c := range s {
		h = ((h << 5) + h) + uint64(c)
	}
	return h
}

// ============================================================
// BERT 分词器
// ============================================================

type bertTokenizer struct {
	vocab map[string]int
}

func newBertTokenizer() *bertTokenizer {
	return &bertTokenizer{
		vocab: initVocab(),
	}
}

func initVocab() map[string]int {
	vocab := make(map[string]int)

	// BERT 基础词汇
	commonWords := []string{
		"[PAD]", "[UNK]", "[CLS]", "[SEP]", "[MASK]",
		"the", "a", "an", "is", "are", "was", "were", "be", "been", "being",
		"have", "has", "had", "having", "do", "does", "did", "doing",
		"will", "would", "could", "should", "may", "might", "must", "can", "need",
		"to", "of", "in", "for", "on", "with", "at", "by", "from", "up", "about",
		"into", "over", "after", "under", "above", "and", "but", "or", "nor", "so",
		"yet", "both", "either", "neither", "not", "only", "own", "same", "than",
		"too", "very", "just", "also", "now", "here", "there", "when", "where",
		"why", "how", "all", "each", "every", "few", "more", "most", "other",
		"some", "such", "no", "any", "this", "that", "these", "those",
		// K8s 词汇
		"pod", "pods", "node", "nodes", "service", "services", "deployment", "deployments",
		"container", "containers", "cluster", "namespace", "namespaces",
		"volume", "volumes", "configmap", "configmaps", "secret", "secrets",
		"ingress", "ingresses", "statefulset", "daemonset", "replicaset",
		"job", "cronjob", "error", "fail", "failed", "failure", "pending",
		"running", "terminated", "crash", "crashed", "restart", "restarted",
		"restarting", "memory", "cpu", "disk", "network", "storage",
		"timeout", "connection", "refused", "unavailable", "killed", "oom",
		"image", "images", "pull", "pulling", "registry", "repository",
		"ready", "notready", "available", "evicted", "eviction", "replica", "replicas",
		"kubectl", "get", "describe", "logs", "apply", "create", "delete", "edit",
		"kubernetes", "kube", "api", "server", "etcd", "scheduler", "controller",
		"evicted", "backoff", "imagepullbackoff", "crashloopbackoff", "oomkilled",
		"notready", "memorypressure", "diskpressure", "terminating",
	}

	for i, w := range commonWords {
		vocab[w] = i
	}

	return vocab
}

func (t *bertTokenizer) tokenize(text string) []int {
	words := strings.Fields(strings.ToLower(text))
	ids := []int{t.vocab["[CLS]"]}

	unk := t.vocab["[UNK]"]
	for _, w := range words {
		if id, ok := t.vocab[w]; ok {
			ids = append(ids, id)
		} else {
			ids = append(ids, unk)
		}
	}

	ids = append(ids, t.vocab["[SEP]"])
	return ids
}
