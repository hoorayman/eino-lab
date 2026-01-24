package stage05

import (
	"context"
	"eino-learn/compose/stage04"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/schema"
)

func RetrieverRAG(query string) []*schema.Document {
	ctx := context.Background()
	timeout := 30 * time.Second
	apiType := ark.APITypeMultiModal
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   "doubao-embedding-vision-250615",
		APIType: &apiType,
		Timeout: &timeout,
	})
	if err != nil {
		panic(err)
	}
	retriever, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      stage04.MilvusCli,
		Collection:  "tt",
		Partition:   nil,
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:      1,
		Embedding: embedder,
	})
	if err != nil {
		panic(err)
	}

	results, err := retriever.Retrieve(ctx, query)
	if err != nil {
		panic(err)
	}

	return results
}
