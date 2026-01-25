package main

import (
	"context"
	"eino-learn/orchestrate/stage02"
	"log"

	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/coze-dev/cozeloop-go"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	client, err := cozeloop.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close(ctx)
	// 在服务 init 时 once 调用
	handler := ccb.NewLoopHandler(client)
	callbacks.AppendGlobalHandlers(handler)
	// stage01.ChatGenerate()
	// stage01.ChatStream()
	// stage02.TemplateChat()
	// stage03.EmbedText()
	// 示例: 使用 IndexerRAG 将文档索引到 Milvus
	// docs := []*schema.Document{
	// 	{
	// 		ID:      "doc_001",
	// 		Content: "人工智能（AI）是计算机科学的一个分支，旨在创建能够执行通常需要人类智能的任务的系统。",
	// 		MetaData: map[string]interface{}{
	// 			"source":   "ai_intro",
	// 			"category": "technology",
	// 		},
	// 	},
	// 	{
	// 		ID:      "doc_002",
	// 		Content: "机器学习是人工智能的子集，它使系统能够从数据中学习并改进，而无需进行显式编程。",
	// 		MetaData: map[string]interface{}{
	// 			"source":   "ml_intro",
	// 			"category": "technology",
	// 		},
	// 	},
	// 	{
	// 		ID:      "doc_003",
	// 		Content: "深度学习是机器学习的一种方法，它使用多层神经网络来模拟人脑的工作方式。",
	// 		MetaData: map[string]interface{}{
	// 			"source":   "dl_intro",
	// 			"category": "technology",
	// 		},
	// 	},
	// }

	// stage04.IndexerRAG(docs)
	// docs := stage05.RetrieverRAG("人工智能")
	// for _, doc := range docs {
	// 	println(fmt.Sprintf("Search result %v", *doc))
	// }
	// stage06.TransDoc()
	// stage07.ToolExample()
	// stage01.OrcChain()
	// stage01.SimpleAgent()
	// stage02.OrcGraph()
	// stage02.OrcGraphWithModel()
	// stage02.OrcGraphWithState()
	// stage02.OrcGraphWithCallback()
	stage02.OutSideOrcGraph()
}
