package helloworld

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func HelloWorldAgent() {
	ctx := context.Background()
	timeout := 30 * time.Second
	// 初始化模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   "doubao-seed-1-8-251228",
		Timeout: &timeout,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建 ChatModelAgent
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "hello_agent",
		Description: "A friendly greeting assistant",
		Instruction: "You are a friendly assistant. Please respond to the user in a warm tone.",
		Model:       model,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建 Runner, agent需要runner才能运行
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})

	// 执行对话
	input := []adk.Message{
		schema.UserMessage("请介绍一下你自己"),
	}

	events := runner.Run(ctx, input)
	for {
		event, ok := events.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			log.Printf("错误: %v", event.Err)
			break
		}

		if msg, err := event.Output.MessageOutput.GetMessage(); err == nil {
			fmt.Printf("Agent: %s\n", msg.Content)
		}
	}
}
