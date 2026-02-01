package loop

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"

	"eino-learn/adk/common/prints"
	"eino-learn/adk/intro/workflow/loop/subagents"
)

func LoopAgent() {
	ctx := context.Background()

	fmt.Println("=== 人类参与的 Agent Loop ===")
	fmt.Println("类似 Claude Code 的交互式智能体迭代模式")

	// 获取用户初始问题
	fmt.Println("\n========== 查询示例 ==========")
	fmt.Println("1. 帮我查看当前目录下有哪些文件")
	fmt.Println("2. 帮我查看 main.go 文件的内容")
	fmt.Println("3. 帮我查看系统的内存使用情况")
	fmt.Println("4. 帮我搜索 main.go 中包含 'func' 的行")
	fmt.Println("5. 帮我做系统健康检查")
	fmt.Println("================================")

	query := getUserInput("请输入您的问题或任务（留空使用默认示例）：")

	// 默认示例查询
	if query == "" {
		query = "帮我查看当前目录下有哪些文件，并查看 main.go 的前 10 行"
		fmt.Printf("\n使用默认查询：%s\n\n", query)
	}

	if query == "" {
		fmt.Println("未输入问题，退出")
		return
	}

	// 创建 LoopAgent
	a, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
		Name:          "reflection_agent",
		Description:   "反思型智能体，包含主智能体和改进智能体，用于迭代式任务解决",
		SubAgents:     []adk.Agent{subagents.NewMainAgent(), subagents.NewCritiqueAgent()},
		MaxIterations: 5,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建 Runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		EnableStreaming: true,
		Agent:           a,
	})

	// 初始运行
	messages := []adk.Message{schema.UserMessage(query)}

	// 开始迭代循环
	iteration := 0
	for {
		iteration++

		// 运行智能体
		iter := runner.Run(ctx, messages)

		var currentResult string
		var hasToolCall bool

		// 收集智能体的响应
		for {
			event, ok := iter.Next()
			if !ok {
				break
			}
			if event.Err != nil {
				fmt.Printf("❌ 错误: %v\n", event.Err)
				break
			}

			prints.Event(event)
			if event.Output != nil {
				msg, _, err := adk.GetMessage(event)
				if err != nil {
					fmt.Printf("❌ 获取消息错误: %v\n", err)
					break
				}

				if msg.Content != "" {
					currentResult = msg.Content
				}

				// 检查是否有工具调用
				if len(msg.ToolCalls) > 0 {
					hasToolCall = true
				}

				// 检查是否有特殊动作（如退出循环）
				if event.Action != nil {
					if event.Action.Exit {
						fmt.Println("\n✓ 智能体认为已完成任务，退出循环")
						fmt.Println("╔═══════════════════════════════════════╗")
						fmt.Println("║              最终结果                    ║")
						fmt.Println("╚═══════════════════════════════════════╝")
						fmt.Printf("%s\n", currentResult)
						return
					}
				}
			}
		}

		// 显示人类交互选项
		fmt.Println()
		fmt.Println("1. 继续下一轮迭代 (让智能体继续改进)")
		fmt.Println("2. 提供反馈 (给出您的改进意见)")
		fmt.Println("3. 修改问题 (调整原始需求)")
		fmt.Println("4. 查看详情 (展开完整输出)")
		fmt.Println("5. 退出循环 (接受当前结果)")
		fmt.Println()

		// 获取用户选择
		choice := getUserInput("请选择操作 [1-5]（默认1）：")

		switch choice {
		case "1", "":
			// 继续下一轮迭代
			fmt.Println("\n✓ 继续下一轮迭代...")

		case "2":
			// 提供反馈
			feedback := getUserInput("请输入您的反馈意见：")
			if feedback != "" {
				messages = append(messages,
					schema.UserMessage(fmt.Sprintf("用户反馈：%s\n请根据这个反馈继续改进您的方案。", feedback)),
				)
				fmt.Println("\n✓ 已加入反馈，继续下一轮迭代...")
			}

		case "3":
			// 修改问题
			newQuery := getUserInput("请输入新的问题：")
			if newQuery != "" {
				messages = []adk.Message{schema.UserMessage(newQuery)}
				iteration = 0
				fmt.Println("\n✓ 已更新问题，重新开始...")
			}

		case "4":
			// 查看详情
			if currentResult != "" {
				fmt.Println("\n╔═══════════════════════════════════════╗")
				fmt.Println("║              完整输出                    ║")
				fmt.Println("╚═══════════════════════════════════════╝")
				fmt.Println(currentResult)
				fmt.Println()
				return
			} else {
				fmt.Println("当前无输出内容")
			}

		case "5":
			// 退出循环
			fmt.Println("\n✓ 用户选择退出")
			fmt.Println("╔═══════════════════════════════════════╗")
			fmt.Println("║              当前结果                    ║")
			fmt.Println("╚═══════════════════════════════════════╝")
			if currentResult != "" {
				fmt.Printf("%s\n", currentResult)
			}
			return

		default:
			fmt.Println("❌ 无效选项，退出")
			return
		}

		// 如果有工具调用被中断，需要将工具响应加入对话
		if hasToolCall {
			fmt.Println("⚠️ 检测到工具调用，请等待工具执行完成...")
		}

		// 检查是否达到最大迭代次数
		if iteration >= 5 {
			fmt.Println("\n⚠️ 已达到最大迭代次数（5轮）")
			fmt.Println("╔═══════════════════════════════════════╗")
			fmt.Println("║              最终结果                    ║")
			fmt.Println("╚═══════════════════════════════════════╝")
			if currentResult != "" {
				fmt.Printf("%s\n", currentResult)
			}
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// getUserInput 获取用户输入
func getUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + " ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取输入错误: %v\n", err)
		return ""
	}
	return strings.TrimSpace(input)
}
