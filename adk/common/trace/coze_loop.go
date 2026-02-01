package trace

import (
	"context"
	"log"
	"os"

	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/coze-dev/cozeloop-go"
)

// **************************************************************
// *** Span 是分布式链路追踪中的基本单位，代表执行链路中的一个操作节点
// *** Span 通过context来构建父子关系，通过context来传递链路追踪信息
// **************************************************************

// CloseFn 关闭 CozeLoop 客户端的函数
type CloseFn func(ctx context.Context)

// EndSpanFn 结束一个 Span 的函数，用于设置输出并完成 Span
type EndSpanFn func(ctx context.Context, output any)

// StartSpanFn 启动一个 Span 的函数，返回新的 context 和结束 Span 的函数
type StartSpanFn func(ctx context.Context, name string, input any) (nCtx context.Context, endFn EndSpanFn)

// AppendCozeLoopCallbackIfConfigured 如果配置了 CozeLoop 环境变量，则初始化并注册 CozeLoop 回调
//
// 这是一个用于链路追踪的配置函数：
// 1. 检查环境变量 COZELOOP_WORKSPACE_ID 和 COZELOOP_API_TOKEN
// 2. 如果配置了，创建 CozeLoop 客户端并注册全局回调处理器
// 3. 返回关闭函数和启动 Span 的函数
//
// 环境变量配置：
//
//	COZELOOP_WORKSPACE_ID=your workspace id
//	COZELOOP_API_TOKEN=your token
//
// 文档地址: https://loop.coze.cn/open/docs/cozeloop/go-sdk
//
// 返回值：
//   - closeFn: 关闭 CozeLoop 客户端的函数
//   - startSpanFn: 启动链路追踪 Span 的函数
func AppendCozeLoopCallbackIfConfigured(_ context.Context) (closeFn CloseFn, startSpanFn StartSpanFn) {
	// 从环境变量读取 CozeLoop 配置
	wsID := os.Getenv("COZELOOP_WORKSPACE_ID")
	apiKey := os.Getenv("COZELOOP_API_TOKEN")

	// 如果没有配置环境变量，返回空函数
	if wsID == "" || apiKey == "" {
		return func(ctx context.Context) {
			// 空 close 函数
		}, buildStartSpanFn(nil) // 传入 nil 表示未配置 CozeLoop
	}

	// 创建 CozeLoop 客户端
	client, err := cozeloop.NewClient(
		cozeloop.WithWorkspaceID(wsID), // 工作空间 ID
		cozeloop.WithAPIToken(apiKey),  // API Token
	)
	if err != nil {
		log.Fatalf("cozeloop.NewClient failed, err: %v", err)
	}

	// 创建 CozeLoop 回调处理器
	handler := ccb.NewLoopHandler(client)

	// 注册为全局回调，这样所有 Eino 的组件（Model、Tool、Graph 等）的执行都会被追踪
	callbacks.AppendGlobalHandlers(handler)

	// 返回关闭函数和启动 Span 的函数
	return client.Close, buildStartSpanFn(client)
}

// buildStartSpanFn 构建 StartSpanFn 函数
//
// 参数：
//   - client: CozeLoop 客户端，如果为 nil 表示未配置 CozeLoop，返回空操作函数
//
// 返回一个函数，该函数用于启动一个链路追踪 Span：
//   - name: Span 的名称，通常标识操作的类型（如 "chat_model", "tool_call"）
//   - input: 操作的输入数据
//   - 返回: 新的 context（包含 span 信息）和结束 Span 的函数
func buildStartSpanFn(client cozeloop.Client) StartSpanFn {
	return func(ctx context.Context, name string, input any) (nCtx context.Context, endFn EndSpanFn) {
		// 如果客户端为 nil，返回空操作函数
		if client == nil {
			return ctx, func(ctx context.Context, output any) {
				// 空 end 函数
			}
		}

		// 启动一个新的 Span
		// name: Span 名称，用于标识操作类型
		// "custom": 表示这是自定义类型的 Span
		nCtx, span := client.StartSpan(ctx, name, "custom")

		// 设置 Span 的输入数据，便于追踪查看
		span.SetInput(ctx, input)

		// 返回包含 span 的 context 和结束 Span 的函数
		return nCtx, buildEndSpanFn(span)
	}
}

// buildEndSpanFn 构建 EndSpanFn 函数
//
// 参数：
//   - span: CozeLoop 的 Span 对象，如果为 nil 表示未配置 CozeLoop，返回空操作函数
//
// 返回一个函数，该函数用于结束一个链路追踪 Span：
//   - output: 操作的输出数据
func buildEndSpanFn(span cozeloop.Span) EndSpanFn {
	return func(ctx context.Context, output any) {
		// 如果 span 为 nil，直接返回
		if span == nil {
			return
		}

		// 设置 Span 的输出数据
		span.SetOutput(ctx, output)

		// 完成 Span，将追踪数据发送到 CozeLoop 平台
		span.Finish(ctx)
	}
}
