package db

import (
	"context"
	"github.com/qdrant/go-client/qdrant"
)

// CreateQueryCollection creates the qdrant colelction
func CreateQueryCollection() error {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		return err
	}

	err = client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "queries",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     384,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	return err
}
