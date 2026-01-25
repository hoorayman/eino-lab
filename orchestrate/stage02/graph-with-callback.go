package stage02

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func OrcGraphWithCallback() {
	ctx := context.Background()
	input := map[string]string{"role": "tsundere", "content": "你好"}
	g := compose.NewGraph[map[string]string, *schema.Message](
		compose.WithGenLocalState(genFunc),
	)
	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (
		output map[string]string, err error) {
		// 在节点内部处理state
		_ = compose.ProcessState(ctx, func(_ context.Context, state *State) error {
			state.History["tsundere_action"] = "我喜欢你"
			state.History["cute_action"] = "摸摸头"
			return nil
		})
		switch input["role"] {
		case "tsundere":
			return map[string]string{"role": "tsundere", "content": input["content"]}, nil
		case "cute":
			return map[string]string{"role": "cute", "content": input["content"]}, nil
		default:
			return map[string]string{"role": "user", "content": input["content"]}, nil
		}
	})
	TsundereLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (
		output []*schema.Message, err error) {
		_ = compose.ProcessState(ctx, func(_ context.Context, state *State) error {
			input["content"] = input["content"] + state.History["tsundere_action"].(string)
			return nil
		})
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个高冷傲娇的大小姐，每次都会用傲娇的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})

	CuteLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (
		output []*schema.Message, err error) {
		// _ = compose.ProcessState[*State](ctx, func(_ context.Context, state *State) error {
		// 	input["content"] = input["content"] + state.History["action"].(string)
		// 	return nil
		// })
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个可爱的小女孩，每次都会用可爱的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})

	cutePreHandler := func(ctx context.Context, input map[string]string, state *State) (map[string]string, error) {
		input["content"] = input["content"] + state.History["cute_action"].(string)
		return input, nil
	}

	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  "doubao-seed-1-8-251228",
	})
	if err != nil {
		panic(err)
	}
	// 注册节点
	err = g.AddLambdaNode("lambda", lambda)
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("tsundere", TsundereLambda)
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("cute", CuteLambda, compose.WithStatePreHandler(cutePreHandler))
	if err != nil {
		panic(err)
	}
	err = g.AddChatModelNode("model", model)
	if err != nil {
		panic(err)
	}
	// 加入分支
	g.AddBranch("lambda", compose.NewGraphBranch(func(ctx context.Context, in map[string]string) (
		endNode string, err error) {
		switch in["role"] {
		case "tsundere":
			return "tsundere", nil
		case "cute":
			return "cute", nil
		default:
			return "tsundere", nil
		}
	}, map[string]bool{"tsundere": true, "cute": true}))

	// 加入边
	err = g.AddEdge(compose.START, "lambda")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("tsundere", "model")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("cute", "model")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("model", compose.END)
	if err != nil {
		panic(err)
	}
	// 编译
	r, err := g.Compile(ctx)
	if err != nil {
		panic(err)
	}
	// 执行
	answer, err := r.Invoke(ctx, input, compose.WithCallbacks(genCallback()))
	if err != nil {
		panic(err)
	}
	fmt.Println(answer.Content)
}

func genCallback() callbacks.Handler {
	handler := callbacks.NewHandlerBuilder().OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		fmt.Printf("当前%s节点输入:%s\n", info.Component, input)
		return ctx
	}).OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		fmt.Printf("当前%s节点输出:%s\n", info.Component, output)
		return ctx
	}).Build()
	return handler
}
