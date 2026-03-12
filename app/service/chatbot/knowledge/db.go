package knowledge

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"freakbot/app/service/chatbot/retrieval"
	"math"
	"os"
	"path/filepath"
)

// DB holds messages and embeddings for retrieval.
// Messages and Embeddings are aligned by index.
type DB struct {
	Messages   []Message
	Embeddings [][]float32
}

type Message struct {
	ID               int    `json:"id"`
	From             string `json:"from"`
	Date             string `json:"date"`
	Text             string `json:"text"`
	ReplyToMessageID int    `json:"reply_to_message_id,omitempty"`

	// The original query message that led to this response.
	// These fields are optional and may be empty for older data.
	QueryID               int    `json:"query_id,omitempty"`
	QueryFrom             string `json:"query_from,omitempty"`
	QueryDate             string `json:"query_date,omitempty"`
	QueryText             string `json:"query_text,omitempty"`
	QueryReplyToMessageID int    `json:"query_reply_to_message_id,omitempty"`
}

// dbJSON is the on-disk representation of messages and embeddings.
// Embeddings are stored as base64-encoded float32 arrays (little-endian).
type dbJSON struct {
	Messages   []Message `json:"messages"`
	Embeddings []string  `json:"embeddings"`
}

const (
	systemPromptFile = "system_prompt.txt"
	dbFile           = "db.json"
)

// Save writes the DB and system prompt to dataDir.
// Embeddings and messages are stored together in a single JSON file, with
// each embedding row encoded as base64 of little-endian float32 values.
func Save(dataDir string, systemPrompt string, db *DB) error {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("mkdir data dir: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dataDir, systemPromptFile), []byte(systemPrompt), 0644); err != nil {
		return fmt.Errorf("write system prompt: %w", err)
	}

	jsonPath := filepath.Join(dataDir, dbFile)

	wire := dbJSON{
		Messages:   db.Messages,
		Embeddings: make([]string, len(db.Embeddings)),
	}
	for i, row := range db.Embeddings {
		b, err := encodeEmbedding(row)
		if err != nil {
			return fmt.Errorf("encode embedding %d: %w", i, err)
		}
		wire.Embeddings[i] = b
	}

	data, err := json.MarshalIndent(&wire, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal db json: %w", err)
	}
	if err := os.WriteFile(jsonPath, data, 0644); err != nil {
		return fmt.Errorf("write db json: %w", err)
	}
	return nil
}

// Load reads the knowledge base from dataDir.
// Returns (nil, error) if any required file is missing.
func Load(dataDir string) (systemPrompt string, db *DB, err error) {
	promptPath := filepath.Join(dataDir, systemPromptFile)
	data, err := os.ReadFile(promptPath)
	if err != nil {
		return "", nil, fmt.Errorf("read system prompt: %w", err)
	}
	systemPrompt = string(data)

	jsonPath := filepath.Join(dataDir, dbFile)
	wireBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return "", nil, fmt.Errorf("read db json: %w", err)
	}

	var wire dbJSON
	if err := json.Unmarshal(wireBytes, &wire); err != nil {
		return "", nil, fmt.Errorf("unmarshal db json: %w", err)
	}
	if len(wire.Messages) != len(wire.Embeddings) {
		return "", nil, fmt.Errorf("messages count %d != embeddings count %d", len(wire.Messages), len(wire.Embeddings))
	}

	embeddings := make([][]float32, len(wire.Embeddings))
	var dim int
	for i, encoded := range wire.Embeddings {
		row, err := decodeEmbedding(encoded)
		if err != nil {
			return "", nil, fmt.Errorf("decode embedding %d: %w", i, err)
		}
		if dim == 0 {
			dim = len(row)
		} else if len(row) != dim {
			return "", nil, fmt.Errorf("embedding dimension mismatch at %d: %d != %d", i, len(row), dim)
		}
		embeddings[i] = row
	}

	return systemPrompt, &DB{Messages: wire.Messages, Embeddings: embeddings}, nil
}

// TopKSimilar returns the indices of the k messages most similar to the query embedding.
func (db *DB) TopKSimilar(query []float32, k int) []int {
	return retrieval.TopKIndices(query, db.Embeddings, k)
}

// encodeEmbedding converts a slice of float32 to base64-encoded bytes (little-endian).
func encodeEmbedding(row []float32) (string, error) {
	if len(row) == 0 {
		return "", nil
	}
	buf := make([]byte, 4*len(row))
	for i, v := range row {
		binary.LittleEndian.PutUint32(buf[i*4:], mathFloat32bits(v))
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

// decodeEmbedding converts a base64-encoded little-endian float32 array back to []float32.
func decodeEmbedding(encoded string) ([]float32, error) {
	if encoded == "" {
		return nil, nil
	}
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}
	if len(b)%4 != 0 {
		return nil, fmt.Errorf("byte length %d is not a multiple of 4", len(b))
	}
	n := len(b) / 4
	out := make([]float32, n)
	for i := 0; i < n; i++ {
		u := binary.LittleEndian.Uint32(b[i*4:])
		out[i] = math.Float32frombits(u)
	}
	return out, nil
}

// mathFloat32bits is a small wrapper so we can keep the encode logic readable.
func mathFloat32bits(f float32) uint32 {
	return math.Float32bits(f)
}
