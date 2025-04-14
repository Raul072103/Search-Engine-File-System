package vectordb

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type EmbeddingRequest struct {
	Text string `json:"text"`
}
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func GetQueryEmbedding(text string) ([]float32, error) {
	reqData := EmbeddingRequest{Text: text}
	body, _ := json.Marshal(reqData)

	resp, err := http.Post("http://localhost:8000/embed", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EmbeddingResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Embedding, err
}
