package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/tools"
)

// 完整的使用示例
func runDemo() {
	// 检查是否有阿里百炼 API Key
	if os.Getenv("DASHSCOPE_API_KEY") == "" {
		fmt.Println("⚠️  请设置 DASHSCOPE_API_KEY 环境变量")
		fmt.Println("   export DASHSCOPE_API_KEY=\"your-dashscope-api-key\"")
		return
	}

	// 运行演示
	if err := runGoutAgentDemo(); err != nil {
		log.Fatalf("演示运行失败: %v", err)
	}
}

func runGoutAgentDemo() error {
	fmt.Println("🩺 痛风化验单分析智能体演示 (Powered by 阿里百炼 Qwen-plus)")
	fmt.Println("═══════════════════════════════════════")

	// 1. 初始化阿里百炼 Qwen LLM
	llm, err := openai.New(
		openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
		openai.WithModel("qwen-plus"),
		openai.WithToken(os.Getenv("DASHSCOPE_API_KEY")),
	)
	if err != nil {
		return fmt.Errorf("初始化阿里百炼 Qwen LLM失败: %w", err)
	}

	// 2. 创建专用工具
	goutAnalyzer := GoutLabAnalyzer{}
	medicalKnowledge := NewMedicalKnowledgeBase()

	agentTools := []tools.Tool{
		goutAnalyzer,
		medicalKnowledge,
		tools.Calculator{},
	}

	// 3. 创建对话记忆
	conversationMemory := memory.NewConversationBuffer()

	// 4. 创建智能体
	agent := agents.NewConversationalAgent(
		llm,
		agentTools,
		agents.WithMaxIterations(3),
	)

	executor := agents.NewExecutor(
		agent,
		agents.WithMemory(conversationMemory),
	)

	fmt.Println("✅ 智能体初始化完成")

	// 5. 演示场景1: 化验单分析
	fmt.Println("\n📋 场景1: 化验单分析")
	fmt.Println("─────────────────────────────────")
	
	labData := `患者化验单数据:
尿酸 520 umol/L (参考范围: 208-428)
C反应蛋白 15.2 mg/L (参考范围: <3.0)
血沉 45 mm/h (参考范围: <15)
肌酐 95 umol/L (参考范围: 54-106)
白细胞 11.2 ×10⁹/L (参考范围: 4.0-10.0)`

	fmt.Println("输入数据:")
	fmt.Println(labData)

	question1 := "请分析这份化验单，评估患者的痛风风险，并给出详细的医学建议：\n" + labData
	
	fmt.Println("\n🤖 智能体分析中...")
	result1, err := chains.Run(context.Background(), executor, question1)
	if err != nil {
		return fmt.Errorf("化验单分析失败: %w", err)
	}

	fmt.Println("\n📊 分析结果:")
	fmt.Println(result1)

	// 6. 演示场景2: 医学知识咨询
	fmt.Println("\n\n📚 场景2: 医学知识咨询")
	fmt.Println("─────────────────────────────────")

	question2 := "什么是痛风？它的主要症状有哪些？如何预防和治疗？"
	
	fmt.Println("问题:", question2)
	fmt.Println("\n🤖 智能体回答中...")
	
	result2, err := chains.Run(context.Background(), executor, question2)
	if err != nil {
		return fmt.Errorf("知识咨询失败: %w", err)
	}

	fmt.Println("\n💡 专业解答:")
	fmt.Println(result2)

	// 7. 演示场景3: 后续咨询
	fmt.Println("\n\n🔄 场景3: 后续咨询 (测试记忆功能)")
	fmt.Println("─────────────────────────────────")

	question3 := "根据前面分析的化验单结果，这位患者需要多久复查一次？"
	
	fmt.Println("问题:", question3)
	fmt.Println("\n🤖 智能体回答中...")
	
	result3, err := chains.Run(context.Background(), executor, question3)
	if err != nil {
		return fmt.Errorf("后续咨询失败: %w", err)
	}

	fmt.Println("\n🔍 智能建议:")
	fmt.Println(result3)

	// 8. 演示场景4: 数值计算
	fmt.Println("\n\n🧮 场景4: 数值计算")
	fmt.Println("─────────────────────────────────")

	question4 := "如果一个患者的尿酸是520 umol/L，正常上限是428 umol/L，请计算超出正常值的百分比。"
	
	fmt.Println("问题:", question4)
	fmt.Println("\n🤖 智能体计算中...")
	
	result4, err := chains.Run(context.Background(), executor, question4)
	if err != nil {
		return fmt.Errorf("数值计算失败: %w", err)
	}

	fmt.Println("\n📐 计算结果:")
	fmt.Println(result4)

	fmt.Println("\n═══════════════════════════════════════")
	fmt.Println("🎉 痛风分析智能体演示完成！")
	fmt.Println("\n🌟 主要特性展示:")
	fmt.Println("✅ 化验单数据智能解析")
	fmt.Println("✅ 风险评估和医学建议")  
	fmt.Println("✅ 专业医学知识咨询")
	fmt.Println("✅ 对话记忆和上下文理解")
	fmt.Println("✅ 数值计算和分析能力")
	fmt.Println("\n💼 适用场景:")
	fmt.Println("• 医疗辅助诊断")
	fmt.Println("• 健康咨询服务")
	fmt.Println("• 医学教育培训")
	fmt.Println("• 健康管理应用")

	return nil
}