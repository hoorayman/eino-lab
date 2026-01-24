package stage04

import (
	"context"
	"log"

	cli "github.com/milvus-io/milvus-sdk-go/v2/client"
)

var MilvusCli cli.Client

func init() {
	ctx := context.Background()
	client, err := cli.NewClient(ctx, cli.Config{
		Address: "192.168.233.128:19530",
		DBName:  "AwesomeEino",
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	MilvusCli = client
}
