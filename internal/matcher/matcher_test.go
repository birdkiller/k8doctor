package matcher

import (
	"testing"

	"k8doctor/internal/cleaner"
)

func TestCosineSimilarity(t *testing.T) {
	// 测试余弦相似度计算
	a := []float32{1, 0, 0}
	b := []float32{1, 0, 0}
	sim := cosineSimilarity(a, b)
	if sim != 1.0 {
		t.Errorf("Expected 1.0 for identical vectors, got %f", sim)
	}

	// 正交向量
	c := []float32{1, 0, 0}
	d := []float32{0, 1, 0}
	sim = cosineSimilarity(c, d)
	if sim != 0.0 {
		t.Errorf("Expected 0.0 for orthogonal vectors, got %f", sim)
	}
}

func TestHashString(t *testing.T) {
	h1 := hashString("pod")
	h2 := hashString("pod")
	if h1 != h2 {
		t.Error("Same string should produce same hash")
	}

	h3 := hashString("node")
	if h1 == h3 {
		t.Error("Different strings should produce different hashes")
	}
}

func TestStringHelperFunctions(t *testing.T) {
	// Test stringsContains
	if !stringsContains("OOMKilled", "oom") {
		t.Error("Expected stringsContains to match substring")
	}

	if stringsContains("Pod", "xyz") {
		t.Error("Expected stringsContains to return false for non-matching")
	}

	// Test containsString
	slice := []string{"abc", "def", "ghi"}
	if !containsString(slice, "def") {
		t.Error("Expected containsString to find 'def'")
	}

	if containsString(slice, "xyz") {
		t.Error("Expected containsString to return false for 'xyz'")
	}
}

func TestTFIDFEmbedding(t *testing.T) {
	m := &Matcher{
		tokenizer: newBertTokenizer(),
	}

	text := "kubernetes pod oomkilled crash"
	vec, err := m.getTFIDFEmbedding(text)

	if err != nil {
		t.Errorf("getTFIDFEmbedding failed: %v", err)
	}

	if len(vec) != VectorDim {
		t.Errorf("Expected vector dim %d, got %d", VectorDim, len(vec))
	}

	// Check if vector is normalized
	norm := float32(0)
	for _, v := range vec {
		norm += v * v
	}
	if norm <= 0 {
		t.Error("Vector should not be zero")
	}

	t.Logf("Text: %s", text)
	t.Logf("Vector (first 10): %v", vec[:10])
}

func TestCleanerIntegration(t *testing.T) {
	// Test the full flow from symptom to embedding
	input := "Pod一直重启，日志报OOM killed"

	symptom := cleaner.Clean(input)

	m := &Matcher{}
	vec, err := m.getTFIDFEmbedding(symptom.Context)

	if err != nil {
		t.Errorf("Embedding failed: %v", err)
	}

	if len(vec) != VectorDim {
		t.Errorf("Expected vector dim %d, got %d", VectorDim, len(vec))
	}

	t.Logf("Symptom: %s", input)
	t.Logf("Cleaned Context: %s", symptom.Context)
	t.Logf("Vector (first 10): %v", vec[:10])
}
