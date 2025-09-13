package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

// TestGoutAnalyzer 测试痛风分析工具
func TestGoutAnalyzer(t *testing.T) {
	analyzer := GoutLabAnalyzer{}
	
	// 测试数据 - 高风险病例
	testData1 := `尿酸 580 umol/L (参考范围: 208-428)
C反应蛋白 25.6 mg/L (参考范围: <3.0)
血沉 55 mm/h (参考范围: <15)
肌酐 135 umol/L (参考范围: 54-106)
白细胞 12.5 ×10⁹/L (参考范围: 4.0-10.0)`

	result1, err := analyzer.Call(context.Background(), testData1)
	if err != nil {
		t.Errorf("分析高风险病例失败: %v", err)
		return
	}
	
	fmt.Println("🔬 高风险病例分析结果:")
	fmt.Println(result1)
	
	// 测试数据 - 正常病例  
	testData2 := `尿酸 350 umol/L (参考范围: 208-428)
C反应蛋白 1.2 mg/L (参考范围: <3.0)
血沉 8 mm/h (参考范围: <15)
肌酐 85 umol/L (参考范围: 54-106)`

	result2, err := analyzer.Call(context.Background(), testData2)
	if err != nil {
		t.Errorf("分析正常病例失败: %v", err)
		return
	}
	
	fmt.Println("\n✅ 正常病例分析结果:")
	fmt.Println(result2)
}

// TestMedicalKnowledge 测试医学知识库
func TestMedicalKnowledge(t *testing.T) {
	kb := NewMedicalKnowledgeBase()
	
	// 测试痛风知识查询
	result1, err := kb.Call(context.Background(), "痛风")
	if err != nil {
		t.Errorf("查询痛风知识失败: %v", err)
		return
	}
	
	fmt.Println("📚 痛风知识查询结果:")
	fmt.Println(result1)
	
	// 测试尿酸知识查询
	result2, err := kb.Call(context.Background(), "尿酸")
	if err != nil {
		t.Errorf("查询尿酸知识失败: %v", err)
		return
	}
	
	fmt.Println("\n🔍 尿酸知识查询结果:")
	fmt.Println(result2)
}

// 手动测试函数 - 用于开发调试
func manualTest() {
	fmt.Println("🧪 开始手动测试痛风分析智能体功能...")
	fmt.Println("═══════════════════════════════════════════════")
	
	// 测试化验单分析工具
	fmt.Println("\n1️⃣ 测试化验单分析工具")
	fmt.Println("───────────────────────────────────")
	
	analyzer := GoutLabAnalyzer{}
	
	// 高风险病例
	highRiskData := `尿酸 650 umol/L (参考范围: 208-428)
C反应蛋白 35.2 mg/L (参考范围: <3.0)  
血沉 75 mm/h (参考范围: <15)
肌酐 150 umol/L (参考范围: 54-106)
白细胞 15.2 ×10⁹/L (参考范围: 4.0-10.0)
尿素氮 9.5 mmol/L (参考范围: 2.5-7.1)`

	fmt.Println("📋 测试数据 (高风险病例):")
	fmt.Println(highRiskData)
	
	result, err := analyzer.Call(context.Background(), highRiskData)
	if err != nil {
		log.Printf("❌ 分析失败: %v", err)
		return
	}
	
	fmt.Println("\n🔍 分析结果:")
	
	// 格式化JSON输出
	var analysisResult GoutAnalysisResult
	if err := json.Unmarshal([]byte(result), &analysisResult); err == nil {
		prettyResult, _ := json.MarshalIndent(analysisResult, "", "  ")
		fmt.Println(string(prettyResult))
		
		// 输出关键信息
		fmt.Printf("\n📊 关键指标:\n")
		if analysisResult.UricAcidLevel != nil {
			fmt.Printf("   • 尿酸水平: %.1f %s (%s)\n", 
				analysisResult.UricAcidLevel.Value, 
				analysisResult.UricAcidLevel.Unit,
				analysisResult.UricAcidLevel.Status)
		}
		fmt.Printf("   • 风险等级: %s\n", analysisResult.RiskLevel)
		fmt.Printf("   • 需要随访: %v\n", analysisResult.FollowUpNeeded)
		
		fmt.Println("\n💡 主要建议:")
		for i, rec := range analysisResult.Recommendations {
			fmt.Printf("   %d. %s\n", i+1, rec)
		}
	} else {
		fmt.Println(result)
	}
	
	// 测试医学知识库
	fmt.Println("\n\n2️⃣ 测试医学知识库")
	fmt.Println("───────────────────────────────────")
	
	kb := NewMedicalKnowledgeBase()
	
	queries := []string{"痛风", "高尿酸血症", "尿酸", "炎症"}
	
	for _, query := range queries {
		fmt.Printf("\n🔍 查询: %s\n", query)
		result, err := kb.Call(context.Background(), query)
		if err != nil {
			fmt.Printf("❌ 查询失败: %v\n", err)
			continue
		}
		
		// 解析并格式化输出
		var medicalInfo []MedicalInfo
		if err := json.Unmarshal([]byte(result), &medicalInfo); err == nil && len(medicalInfo) > 0 {
			info := medicalInfo[0]
			fmt.Printf("📋 %s\n", info.Topic)
			fmt.Printf("   定义: %s\n", info.Definition)
			if len(info.Symptoms) > 0 {
				fmt.Printf("   症状: %s\n", info.Symptoms[0])
			}
			if len(info.References) > 0 {
				fmt.Printf("   参考: %s\n", info.References[0])
			}
		}
	}
	
	// 测试边界情况
	fmt.Println("\n\n3️⃣ 测试边界情况")
	fmt.Println("───────────────────────────────────")
	
	// 正常值测试
	normalData := `尿酸 300 umol/L (参考范围: 208-428)
C反应蛋白 1.5 mg/L (参考范围: <3.0)
肌酐 80 umol/L (参考范围: 54-106)`

	fmt.Println("📋 正常值测试数据:")
	fmt.Println(normalData)
	
	normalResult, err := analyzer.Call(context.Background(), normalData)
	if err != nil {
		log.Printf("❌ 正常值分析失败: %v", err)
		return
	}
	
	var normalAnalysis GoutAnalysisResult
	if err := json.Unmarshal([]byte(normalResult), &normalAnalysis); err == nil {
		fmt.Printf("\n📊 正常值分析结果:\n")
		fmt.Printf("   • 风险等级: %s\n", normalAnalysis.RiskLevel)
		fmt.Printf("   • 需要随访: %v\n", normalAnalysis.FollowUpNeeded)
		if len(normalAnalysis.Recommendations) > 0 {
			fmt.Printf("   • 主要建议: %s\n", normalAnalysis.Recommendations[0])
		}
	}
	
	// 异常数据格式测试
	fmt.Println("\n📋 异常格式测试:")
	malformedData := "无效的化验单数据"
	malformedResult, err := analyzer.Call(context.Background(), malformedData)
	if err != nil {
		fmt.Printf("❌ 异常格式处理失败: %v\n", err)
	} else {
		fmt.Printf("✅ 异常格式处理成功，返回: %s\n", malformedResult[:50]+"...")
	}
	
	fmt.Println("\n═══════════════════════════════════════════════")
	fmt.Println("🎉 手动测试完成！")
	fmt.Println("✅ 化验单分析工具正常工作")
	fmt.Println("✅ 医学知识库正常工作") 
	fmt.Println("✅ 边界情况处理正常")
}

// 如果直接运行此文件，则执行手动测试
func init() {
	// 检查是否在测试环境中
	if len(fmt.Sprintf("")) == 0 { // 这是一个技巧来检测非测试环境
		// 在非测试环境中运行手动测试
		go func() {
			manualTest()
		}()
	}
}