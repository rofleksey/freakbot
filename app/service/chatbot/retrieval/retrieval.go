package retrieval

import (
	"math"
	"sort"
)

// CosineSimilarity returns the cosine similarity between a and b (same length).
// Values are in [-1, 1]. Higher means more similar.
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return float32(dot / (math.Sqrt(normA) * math.Sqrt(normB)))
}

// TopKIndices returns the indices of the top k vectors in embeddings
// by cosine similarity to query. Ties are broken by index order.
// If k >= len(embeddings), returns all indices in order of similarity.
func TopKIndices(query []float32, embeddings [][]float32, k int) []int {
	type scoreIndex struct {
		score float32
		idx   int
	}
	scores := make([]scoreIndex, 0, len(embeddings))
	for i, emb := range embeddings {
		s := CosineSimilarity(query, emb)
		scores = append(scores, scoreIndex{s, i})
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})
	if k > len(scores) {
		k = len(scores)
	}
	out := make([]int, k)
	for i := 0; i < k; i++ {
		out[i] = scores[i].idx
	}
	return out
}
