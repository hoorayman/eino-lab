package stage01

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func OrcChain() {
	ctx := context.Background()
	timeout := 30 * time.Second
	// 初始化模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   "doubao-seed-1-8-251228",
		Timeout: &timeout,
	})
	if err != nil {
		panic(err)
	}
	// 创建一个Lambda节点，用于处理输入的文本
	// 注意节点之间输入和输出的类型要匹配
	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output []*schema.Message, err error) {
		finalMessage := input + "回答结尾加上: 你的朋友，有什么问题可以问我"
		output = []*schema.Message{
			{
				Role:    schema.User,
				Content: finalMessage,
			},
		}
		return output, nil
	})

	// 链的入口和出口固定为在 NewChain[In, Out] 里写的这两个类型
	// 注册一条chain
	chain := compose.NewChain[string, *schema.Message]()
	// 给chain添加各个节点
	chain.AppendLambda(lambda).AppendChatModel(model)
	// 编译chain
	r, err := chain.Compile(ctx)
	if err != nil {
		panic(fmt.Sprintf("编译链失败: %v", err))
	}
	// 执行chain
	answer, err := r.Invoke(ctx, "你好，请告诉我你的名字")
	if err != nil {
		panic(fmt.Sprintf("执行链失败: %v", err))
	}
	fmt.Println(answer.Content)
}
