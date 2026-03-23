package cleaner

import (
	"regexp"
	"strings"
)

// Symptom 清洗后的症状结构
type Symptom struct {
	Resources   []string  // 资源类型: Pod, Node, Service, Deployment 等
	States      []string  // 状态: OOMKilled, CrashLoopBackOff 等
	Keywords    []string  // 关键词
	ErrorCodes  []string  // 错误码: 137, 1, 404 等
	Context     string    // 清洗后的上下文摘要（用于 Embedding）
	RawInput    string    // 原始用户输入
}

// 已知资源类型
var resourcePatterns = []string{
	"Pod", "pod", "POD",
	"Node", "node", "NODE",
	"Service", "service", "svc",
	"Deployment", "deployment", "deploy",
	"StatefulSet", "statefulset",
	"DaemonSet", "daemonset",
	"ReplicaSet", "replicaset",
	"Job", "job", "CronJob", "cronjob",
	"Ingress", "ingress",
	"ConfigMap", "configmap",
	"Secret", "secret",
	"PVC", "pvc", "PersistentVolumeClaim",
	"StorageClass", "storageclass",
	"Endpoint", "endpoint",
	"Namespace", "namespace",
}

// 已知状态关键词
var statePatterns = []string{
	"OOMKilled", "OOM", "oom",
	"CrashLoopBackOff", "crashloop", "CrashLoop",
	"ImagePullBackOff", "imagepull", "ImagePull",
	"ErrImagePull", "errimagepull",
	"Pending", "pending",
	"Running", "running",
	"NotReady", "notready", "NotReady",
	"Ready", "ready",
	"Evicted", "evicted",
	"Terminating", "terminating",
	"BackOff", "backoff",
	"Failed", "failed",
	"Error", "error",
	"Unknown", "unknown",
}

// 已知错误码
var errorCodePatterns = map[string]string{
	"137": "OOMKilled",
	"143": "GracefulShutdown",
	"1":   "GeneralError",
	"2":   "Misuse",
	"126": "PermissionDenied",
	"127": "CommandNotFound",
	"139": "SegmentationFault",
	"134": "Aborted",
	"136": "ArithmeticError",
	"255": "ExitError",
}

// Clean 对用户输入的症状描述进行清洗
func Clean(input string) *Symptom {
	s := &Symptom{
		RawInput: input,
	}

	// 提取资源类型
	s.Resources = extractMatches(input, resourcePatterns)

	// 提取状态
	s.States = extractMatches(input, statePatterns)

	// 提取错误码
	s.ErrorCodes = extractErrorCodes(input)

	// 提取关键词（去除已匹配的部分）
	s.Keywords = extractKeywords(input)

	// 构建用于 Embedding 的上下文
	s.Context = buildContext(s)

	return s
}

// extractMatches 从文本中提取匹配项
func extractMatches(text string, patterns []string) []string {
	found := make([]string, 0)
	seen := make(map[string]bool)

	for _, pattern := range patterns {
		regex := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(pattern))
		if regex.MatchString(text) {
			// 标准化：首字母大写
			normalized := normalizeState(pattern)
			if !seen[normalized] {
				found = append(found, normalized)
				seen[normalized] = true
			}
		}
	}

	return found
}

// extractErrorCodes 提取错误码
func extractErrorCodes(text string) []string {
	found := make([]string, 0)
	seen := make(map[string]bool)

	// 匹配退出码格式
	re := regexp.MustCompile(`\b(exit\s*)?code[=:\s]*(\d+)\b`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) == 3 {
			code := match[2]
			if _, ok := errorCodePatterns[code]; ok {
				if !seen[code] {
					found = append(found, code)
					seen[code] = true
				}
			}
		}
	}

	return found
}

// extractKeywords 提取剩余关键词
func extractKeywords(text string) []string {
	// 移除已匹配的内容后，提取剩余有意义的词
	cleaned := text

	// 移除错误码
	re := regexp.MustCompile(`(?i)\b(exit\s*)?code[=:\s]*\d+\b`)
	cleaned = re.ReplaceAllString(cleaned, "")

	// 移除状态词
	for _, state := range statePatterns {
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(state))
		cleaned = re.ReplaceAllString(cleaned, "")
	}

	// 移除资源类型
	for _, res := range resourcePatterns {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(res) + `\b`)
		cleaned = re.ReplaceAllString(cleaned, "")
	}

	// 提取剩余词
	words := strings.Fields(cleaned)
	keywords := make([]string, 0)
	for _, word := range words {
		word = strings.Trim(word, ".,;:!?()[]{}")
		if len(word) > 2 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// buildContext 构建用于 Embedding 的上下文
func buildContext(s *Symptom) string {
	parts := make([]string, 0)

	if len(s.Resources) > 0 {
		parts = append(parts, strings.Join(s.Resources, "/"))
	}

	if len(s.States) > 0 {
		parts = append(parts, strings.Join(s.States, ", "))
	}

	if len(s.Keywords) > 0 {
		parts = append(parts, strings.Join(s.Keywords, " "))
	}

	return strings.Join(parts, " - ")
}

// normalizeState 标准化状态名称
func normalizeState(s string) string {
	switch strings.ToLower(s) {
	case "pod":
		return "Pod"
	case "node":
		return "Node"
	case "service", "svc":
		return "Service"
	case "deployment", "deploy":
		return "Deployment"
	case "oom", "oomkilled":
		return "OOMKilled"
	case "crashloop", "crashloopbackoff":
		return "CrashLoopBackOff"
	case "imagepull", "imagepullbackoff":
		return "ImagePullBackOff"
	case "errimagepull":
		return "ErrImagePull"
	case "pending":
		return "Pending"
	case "notready":
		return "NotReady"
	case "evicted":
		return "Evicted"
	case "backoff":
		return "BackOff"
	case "failed":
		return "Failed"
	default:
		return s
	}
}
