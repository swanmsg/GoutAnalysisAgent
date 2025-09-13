package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/tools"
)

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "demo":
			// 运行演示模式
			runDemo()
			return
		case "test":
			// 运行测试模式
			runTestsMain()
			return
		case "example":
			// 运行示例模式
			if err := runExample(); err != nil {
				fmt.Fprintf(os.Stderr, "示例运行错误: %v\n", err)
				os.Exit(1)
			}
			return
		case "help", "-h", "--help":
			printUsage()
			return
		}
	}

	// 默认运行交互模式
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("🩺 痛风化验单分析智能体 (Powered by 阿里百炼 Qwen-plus)")
	fmt.Println("使用方法:")
	fmt.Println("  go run *.go          - 交互式对话模式 (默认)")
	fmt.Println("  go run *.go demo     - 运行演示模式")
	fmt.Println("  go run *.go test     - 运行测试模式")
	fmt.Println("  go run *.go example  - 运行简单示例")
	fmt.Println("  go run *.go help     - 显示此帮助信息")
	fmt.Println("")
	fmt.Println("环境变量:")
	fmt.Println("  DASHSCOPE_API_KEY    - 阿里百炼 API 密钥 (必需)")
}

func run() error {
	// 检查环境变量
	if os.Getenv("DASHSCOPE_API_KEY") == "" {
		fmt.Println("⚠️  请设置 DASHSCOPE_API_KEY 环境变量")
		fmt.Println("   export DASHSCOPE_API_KEY=\"your-dashscope-api-key\"")
		fmt.Println("   获取 API Key: https://dashscope.aliyun.com/")
		return fmt.Errorf("缺少必需的环境变量")
	}
	
	fmt.Println("🩺 痛风化验单分析智能体启动中... (Powered by 阿里百炼 Qwen-plus)")
	fmt.Println("═══════════════════════════════════════")
	
	// 初始化阿里百炼 Qwen LLM (通过OpenAI兼容接口)
	llm, err := openai.New(
		openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
		openai.WithModel("qwen-plus"),
		openai.WithToken(os.Getenv("DASHSCOPE_API_KEY")),
	)
	if err != nil {
		return fmt.Errorf("初始化阿里百炼 Qwen LLM失败: %w", err)
	}

	// 创建专用工具
	goutAnalyzer := GoutLabAnalyzer{}
	medicalKnowledge := NewMedicalKnowledgeBase()

	// 工具列表
	agentTools := []tools.Tool{
		goutAnalyzer,
		medicalKnowledge,
		tools.Calculator{}, // 添加计算器工具用于数值计算
	}

	// 创建对话记忆
	conversationMemory := memory.NewConversationBuffer()

	// 创建对话型智能体
	agent := agents.NewConversationalAgent(
		llm,
		agentTools,
		agents.WithMaxIterations(5),
	)

	// 创建执行器
	executor := agents.NewExecutor(
		agent,
		agents.WithMemory(conversationMemory),
	)

	fmt.Println("✅ 智能体初始化完成！")
	fmt.Println("\n🔬 我是您的痛风化验单分析助手，可以帮您：")
	fmt.Println("   • 分析化验单数据，评估痛风风险")
	fmt.Println("   • 提供痛风相关医学知识")
	fmt.Println("   • 给出个性化的健康建议")
	fmt.Println("   • 解答痛风相关疑问")
	fmt.Println("\n💡 使用示例：")
	fmt.Println("   1. 输入化验单数据进行分析")
	fmt.Println("   2. 询问痛风相关医学知识")
	fmt.Println("   3. 咨询治疗和预防建议")
	fmt.Println("\n输入 'exit' 退出程序")
	fmt.Println("═══════════════════════════════════════")

	// 交互式对话循环
	for {
		fmt.Print("\n👨‍⚕️ 您的问题: ")
		
		var input string
		fmt.Scanln(&input)
		
		if strings.ToLower(strings.TrimSpace(input)) == "exit" {
			fmt.Println("\n👋 感谢使用痛风分析智能体，祝您身体健康！")
			break
		}

		if strings.TrimSpace(input) == "" {
			fmt.Println("请输入您的问题或化验单数据。")
			continue
		}

		// 执行智能体处理
		fmt.Println("\n🔍 分析中...")
		result, err := chains.Run(context.Background(), executor, input)
		if err != nil {
			fmt.Printf("❌ 处理过程中出现错误: %v\n", err)
			continue
		}

		fmt.Println("\n📋 分析结果:")
		fmt.Println("───────────────────────────────────────")
		fmt.Println(result)
		fmt.Println("───────────────────────────────────────")
		
		// 提示继续对话
		fmt.Println("\n💬 如有其他问题，请继续输入，或输入 'exit' 退出")
	}

	return nil
}

// 示例用法函数
func runExample() error {
	// 初始化阿里百炼 Qwen LLM
	llm, err := openai.New(
		openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
		openai.WithModel("qwen-plus"),
		openai.WithToken(os.Getenv("DASHSCOPE_API_KEY")),
	)
	if err != nil {
		return err
	}

	// 创建工具
	goutAnalyzer := GoutLabAnalyzer{}
	medicalKnowledge := NewMedicalKnowledgeBase()
	
	agentTools := []tools.Tool{
		goutAnalyzer,
		medicalKnowledge,
	}

	// 创建智能体
	agent := agents.NewConversationalAgent(llm, agentTools)
	executor := agents.NewExecutor(agent)

	// 示例化验单数据
	labData := `尿酸 520 umol/L (参考范围: 208-428)
C反应蛋白 15.2 mg/L (参考范围: <3.0)
血沉 45 mm/h (参考范围: <15)
肌酐 95 umol/L (参考范围: 54-106)`

	fmt.Println("🔬 分析示例化验单数据:")
	fmt.Println(labData)

	result, err := chains.Run(context.Background(), executor, 
		"请分析这份化验单，评估痛风风险并给出建议：\n"+labData)
	
	if err != nil {
		return err
	}

	fmt.Println("\n📋 分析结果:")
	fmt.Println(result)

	// 知识查询示例
	fmt.Println("\n\n📚 查询痛风相关知识:")
	result2, err := chains.Run(context.Background(), executor, 
		"什么是痛风？有哪些症状和治疗方法？")
	
	if err != nil {
		return err
	}

	fmt.Println(result2)

	return nil
}