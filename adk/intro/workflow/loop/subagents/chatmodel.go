package subagents

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"

	"eino-learn/adk/common/model"
)

func NewMainAgent() adk.Agent {
	// 创建命令行工具
	shellTool, err := utils.InferTool("execute_command", "执行系统命令",
		func(ctx context.Context, req *ExecuteCommandRequest) (string, error) {
			// 安全检查：只允许特定命令
			// safeCommands := map[string]bool{
			// 	"ls":     true,
			// 	"pwd":    true,
			// 	"cat":    true,
			// 	"echo":   true,
			// 	"date":   true,
			// 	"whoami": true,
			// 	"uname":  true,
			// 	"uptime": true,
			// 	"df":     true,
			// 	"free":   true,
			// 	"ps":     true,
			// 	"grep":   true,
			// 	"wc":     true,
			// 	"head":   true,
			// 	"tail":   true,
			// 	"sort":   true,
			// 	"uniq":   true,
			// }

			// 使用 shell 执行命令，以正确处理引号和特殊字符
			// Windows: cmd /c, Unix: sh -c
			var cmd *exec.Cmd
			if os.Getenv("OS") == "Windows_NT" {
				cmd = exec.CommandContext(ctx, "cmd", "/c", req.Command)
			} else {
				cmd = exec.CommandContext(ctx, "sh", "-c", req.Command)
			}

			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Sprintf("命令执行失败: %v\n输出: %s", err, string(output)), nil
			}

			return string(output), nil
		})
	if err != nil {
		log.Fatalf("创建命令行工具失败，name=%v, err=%v", "execute_command", err)
	}

	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "main_agent",
		Description: "主智能体，负责尝试解决用户的任务",
		Instruction: `你是负责解决用户任务的主智能体。

你的任务：
1. 仔细理解用户的原始问题
2. 使用适当的工具（如 execute_command）完成任务
3. 根据反馈智能体的建议改进你的方案

重要：
- 如果收到反馈智能体的改进建议，请认真对待并在下一轮中改进
- 不断优化你的答案，直到提供完整、准确的解决方案
- 使用命令工具时，确保命令格式正确，特别是引号和特殊字符`,
		Model: model.NewChatModel(),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					shellTool,
				},
			},
			ReturnDirectly: map[string]bool{},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}

func NewCritiqueAgent() adk.Agent {
	exitAndSummarizeTool, err := utils.InferTool("exit_and_summarize", "退出循环并提供最终总结",
		func(ctx context.Context, req *exitAndSummarize) (string, error) {
			_ = adk.SendToolGenAction(ctx, "exit_and_summarize", adk.NewBreakLoopAction("critique_agent"))
			return req.Summary, nil
		})
	if err != nil {
		log.Fatalf("create tool failed, name=%v, err=%v", "exit_and_summarize", err)
	}
	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "critique_agent",
		Description: "反馈智能体，负责对主智能体的工作提出补充改进",
		Instruction: `你是负责反馈主智能体工作的反馈智能体。

你的任务：
1. 审查主智能体的方案和执行结果
2. 如果发现问题或需要改进，提供具体、明确的反馈意见
3. 如果主智能体的回答令人满意，调用 'exit_and_summarize' 工具并总结最终结果

重要：
- 你输出的反馈会直接传递给主智能体，用于下一轮改进
- 反馈要具体明确，指出问题和改进方向
- 不要只是重复问题，要给出建设性建议`,
		Model: model.NewChatModel(),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					exitAndSummarizeTool,
				},
			},
			ReturnDirectly: map[string]bool{
				"exit_and_summarize": true,
			},
		},
	})
	if err != nil {
		log.Fatalf("create agent failed, name=%v, err=%v", "critique_agent", err)
	}
	return a
}

type exitAndSummarize struct {
	Summary string `json:"summary" jsonschema_description:"解决方案的最终总结"`
}

// ExecuteCommandRequest 执行命令的请求参数
type ExecuteCommandRequest struct {
	Command string `json:"command" jsonschema_description:"要执行的命令（如：ls -la, pwd, cat file.txt）"`
}

// getSafeCommandList 获取允许的安全命令列表
func getSafeCommandList() string {
	return "ls, pwd, cat, echo, date, whoami, uname, uptime, df, free, ps, grep, wc, head, tail, sort, uniq"
}
