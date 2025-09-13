package main

import (
	"context"
	"encoding/json"
	"fmt"
)

func runTestsMain() {
	runTests()
}

// 独立测试运行器
func runTests() {
	fmt.Println("🧪 痛风分析智能体功能测试")
	fmt.Println("═══════════════════════════════════════")
	
	// 1. 测试化验单分析工具
	testGoutAnalyzer()
	
	// 2. 测试医学知识库
	testMedicalKnowledge()
	
	fmt.Println("\n🎉 所有测试完成！")
}

func testGoutAnalyzer() {
	fmt.Println("\n1️⃣ 测试化验单分析工具")
	fmt.Println("─────────────────────────────────")
	
	analyzer := GoutLabAnalyzer{}
	
	// 高风险病例测试
	fmt.Println("📋 高风险病例测试:")
	highRiskData := `尿酸 580 umol/L (参考范围: 208-428)
C反应蛋白 25.6 mg/L (参考范围: <3.0)
血沉 55 mm/h (参考范围: <15)
肌酐 135 umol/L (参考范围: 54-106)`
	
	result, err := analyzer.Call(context.Background(), highRiskData)
	if err != nil {
		fmt.Printf("❌ 分析失败: %v\n", err)
		return
	}
	
	// 解析并显示关键信息
	var analysis GoutAnalysisResult
	if err := json.Unmarshal([]byte(result), &analysis); err == nil {
		fmt.Printf("✅ 分析成功\n")
		fmt.Printf("   🔍 风险等级: %s\n", analysis.RiskLevel)
		fmt.Printf("   📊 需要随访: %v\n", analysis.FollowUpNeeded)
		if analysis.UricAcidLevel != nil {
			fmt.Printf("   🩸 尿酸水平: %.1f %s (%s)\n", 
				analysis.UricAcidLevel.Value, 
				analysis.UricAcidLevel.Unit,
				analysis.UricAcidLevel.Status)
		}
		fmt.Printf("   💡 建议数量: %d 条\n", len(analysis.Recommendations))
	} else {
		fmt.Printf("⚠️  解析结果格式异常: %v\n", err)
	}
	
	// 正常值测试
	fmt.Println("\n📋 正常值测试:")
	normalData := `尿酸 320 umol/L (参考范围: 208-428)
C反应蛋白 2.1 mg/L (参考范围: <3.0)
肌酐 85 umol/L (参考范围: 54-106)`
	
	normalResult, err := analyzer.Call(context.Background(), normalData)
	if err != nil {
		fmt.Printf("❌ 分析失败: %v\n", err)
		return
	}
	
	var normalAnalysis GoutAnalysisResult
	if err := json.Unmarshal([]byte(normalResult), &normalAnalysis); err == nil {
		fmt.Printf("✅ 分析成功\n")
		fmt.Printf("   🔍 风险等级: %s\n", normalAnalysis.RiskLevel)
		fmt.Printf("   📊 需要随访: %v\n", normalAnalysis.FollowUpNeeded)
	}
}

func testMedicalKnowledge() {
	fmt.Println("\n2️⃣ 测试医学知识库")
	fmt.Println("─────────────────────────────────")
	
	kb := NewMedicalKnowledgeBase()
	
	testQueries := []string{"痛风", "尿酸", "高尿酸血症", "炎症"}
	
	for _, query := range testQueries {
		fmt.Printf("🔍 查询: %s\n", query)
		
		result, err := kb.Call(context.Background(), query)
		if err != nil {
			fmt.Printf("❌ 查询失败: %v\n", err)
			continue
		}
		
		// 验证返回结果
		var medicalInfo []MedicalInfo
		if err := json.Unmarshal([]byte(result), &medicalInfo); err == nil {
			if len(medicalInfo) > 0 {
				fmt.Printf("✅ 找到相关知识: %s\n", medicalInfo[0].Topic)
				if len(medicalInfo[0].Definition) > 50 {
					fmt.Printf("   📋 定义: %s...\n", medicalInfo[0].Definition[:50])
				}
			} else {
				fmt.Printf("⚠️  未找到相关知识\n")
			}
		} else {
			// 可能是错误信息
			if len(result) > 100 {
				fmt.Printf("⚠️  返回信息: %s...\n", result[:100])
			} else {
				fmt.Printf("⚠️  返回信息: %s\n", result)
			}
		}
	}
	
	// 测试未知查询
	fmt.Printf("\n🔍 未知查询测试: 糖尿病\n")
	unknownResult, err := kb.Call(context.Background(), "糖尿病")
	if err != nil {
		fmt.Printf("❌ 查询失败: %v\n", err)
	} else {
		if len(unknownResult) > 100 {
			fmt.Printf("✅ 正确处理未知查询: %s...\n", unknownResult[:100])
		} else {
			fmt.Printf("✅ 正确处理未知查询: %s\n", unknownResult)
		}
	}
}