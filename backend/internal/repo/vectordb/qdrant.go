package vectordb

import (
	"context"
	"github.com/qdrant/go-client/qdrant"
)

type Repo interface {
	StoreQuery(id uint64, text string) error
	SuggestSimilarQueries(input string, limit *uint64) ([]string, error)
}

type repo struct {
	client *qdrant.Client
}

func New(client *qdrant.Client) Repo {
	return &repo{client: client}
}

func (r *repo) StoreQuery(id uint64, text string) error {
	embedding, err := GetQueryEmbedding(text)
	if err != nil {
		return err
	}

	_, err = r.client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "queries",
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(id),
				Vectors: qdrant.NewVectors(embedding...),
				Payload: qdrant.NewValueMap(map[string]any{"text": text}),
			},
		},
	})
	return err
}

func (r *repo) SuggestSimilarQueries(input string, limit *uint64) ([]string, error) {
	embedding, err := GetQueryEmbedding(input)
	if err != nil {
		return nil, err
	}

	searchResult, err := r.client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "queries",
		Query:          qdrant.NewQuery(embedding...),
		Limit:          limit,
		WithPayload:    qdrant.NewWithPayloadInclude("text"),
	})

	var suggestions = make([]string, 0)
	for _, hit := range searchResult {
		if val, ok := hit.Payload["text"]; ok {
			if strVal := val.GetStringValue(); strVal != "" {
				suggestions = append(suggestions, strVal)
			}
		}
	}

	return suggestions, nil
}
