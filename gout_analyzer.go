package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
)

// GoutLabAnalyzer 痛风化验单分析工具
type GoutLabAnalyzer struct {
	CallbacksHandler callbacks.Handler
}

// LabResult 化验结果结构
type LabResult struct {
	Parameter    string  `json:"parameter"`    // 检测项目名称
	Value        float64 `json:"value"`        // 检测值
	Unit         string  `json:"unit"`         // 单位
	ReferenceMin float64 `json:"reference_min"` // 参考值下限
	ReferenceMax float64 `json:"reference_max"` // 参考值上限
	Status       string  `json:"status"`       // 正常/偏高/偏低
}

// GoutAnalysisResult 痛风分析结果
type GoutAnalysisResult struct {
	UricAcidLevel    *LabResult   `json:"uric_acid_level"`    // 尿酸水平
	InflammatoryMarkers []LabResult `json:"inflammatory_markers"` // 炎症指标
	KidneyFunction   []LabResult  `json:"kidney_function"`    // 肾功能指标
	RiskLevel        string       `json:"risk_level"`         // 风险等级: 低风险/中风险/高风险
	Recommendations  []string     `json:"recommendations"`    // 建议
	FollowUpNeeded   bool         `json:"follow_up_needed"`   // 是否需要随访
}

// Name 返回工具名称
func (g GoutLabAnalyzer) Name() string {
	return "gout_lab_analyzer"
}

// Description 返回工具描述
func (g GoutLabAnalyzer) Description() string {
	return `痛风化验单分析工具。用于分析血液生化检查结果，特别是尿酸水平和相关指标。
输入格式应包含化验项目名称、数值、单位和参考范围，例如：
"尿酸 520 umol/L (参考范围: 208-428)"
"C反应蛋白 15.2 mg/L (参考范围: <3.0)"
该工具会分析各项指标，评估痛风风险，并提供相应的医学建议。`
}

// Call 执行化验单分析
func (g GoutLabAnalyzer) Call(ctx context.Context, input string) (string, error) {
	if g.CallbacksHandler != nil {
		g.CallbacksHandler.HandleToolStart(ctx, input)
	}

	// 解析输入的化验单数据
	labResults, err := g.parseLabInput(input)
	if err != nil {
		return fmt.Sprintf("解析化验单数据时出错: %v", err), nil
	}

	// 分析化验结果
	analysis := g.analyzeGoutRisk(labResults)

	// 格式化输出结果
	result, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return fmt.Sprintf("格式化分析结果时出错: %v", err), nil
	}

	if g.CallbacksHandler != nil {
		g.CallbacksHandler.HandleToolEnd(ctx, string(result))
	}

	return string(result), nil
}

// parseLabInput 解析化验单输入数据
func (g *GoutLabAnalyzer) parseLabInput(input string) ([]LabResult, error) {
	var results []LabResult
	lines := strings.Split(input, "\n")

	// 正则表达式匹配化验项目格式
	// 匹配格式如: "尿酸 520 umol/L (参考范围: 208-428)" 或 "C反应蛋白 15.2 mg/L (<3.0)"
	re := regexp.MustCompile(`([^0-9]+?)\s*([0-9]+\.?[0-9]*)\s*([a-zA-Z/μmol]+).*?(?:参考范围?[：:]?\s*([<>]?)([0-9]+\.?[0-9]*)\s*[-~至]\s*([0-9]+\.?[0-9]*)|[（(<]\s*([<>]?)([0-9]+\.?[0-9]*)\s*[）)>]?)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 4 {
			parameter := strings.TrimSpace(matches[1])
			value, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				continue
			}
			unit := strings.TrimSpace(matches[3])

			result := LabResult{
				Parameter: parameter,
				Value:     value,
				Unit:      unit,
			}

			// 解析参考范围
			if len(matches) >= 7 && matches[5] != "" && matches[6] != "" {
				// 范围格式: 208-428
				min, err1 := strconv.ParseFloat(matches[5], 64)
				max, err2 := strconv.ParseFloat(matches[6], 64)
				if err1 == nil && err2 == nil {
					result.ReferenceMin = min
					result.ReferenceMax = max
				}
			} else if len(matches) >= 9 && matches[8] != "" {
				// 单值格式: <3.0 或 >100
				refValue, err := strconv.ParseFloat(matches[8], 64)
				if err == nil {
					operator := matches[7]
					if operator == "<" {
						result.ReferenceMax = refValue
						result.ReferenceMin = 0
					} else if operator == ">" {
						result.ReferenceMin = refValue
						result.ReferenceMax = 999999
					}
				}
			}

			// 判断状态
			result.Status = g.determineStatus(result)
			results = append(results, result)
		}
	}

	return results, nil
}

// determineStatus 判断检测结果状态
func (g *GoutLabAnalyzer) determineStatus(result LabResult) string {
	if result.ReferenceMax > 0 && result.Value > result.ReferenceMax {
		return "偏高"
	}
	if result.ReferenceMin > 0 && result.Value < result.ReferenceMin {
		return "偏低"
	}
	return "正常"
}

// analyzeGoutRisk 分析痛风风险
func (g *GoutLabAnalyzer) analyzeGoutRisk(results []LabResult) GoutAnalysisResult {
	analysis := GoutAnalysisResult{
		InflammatoryMarkers: []LabResult{},
		KidneyFunction:      []LabResult{},
		Recommendations:     []string{},
	}

	var uricAcidHigh bool
	var inflammationPresent bool
	var kidneyIssues bool

	// 分析各项指标
	for _, result := range results {
		parameterLower := strings.ToLower(result.Parameter)
		
		// 尿酸分析
		if strings.Contains(parameterLower, "尿酸") || strings.Contains(parameterLower, "uric") {
			analysis.UricAcidLevel = &result
			if result.Status == "偏高" {
				uricAcidHigh = true
				if result.Value > 500 {
					analysis.Recommendations = append(analysis.Recommendations, "尿酸水平显著升高，建议立即就医，考虑药物治疗")
				} else if result.Value > 450 {
					analysis.Recommendations = append(analysis.Recommendations, "尿酸水平偏高，建议调整饮食，限制高嘌呤食物摄入")
				}
			}
		}

		// 炎症指标分析
		if strings.Contains(parameterLower, "c反应蛋白") || strings.Contains(parameterLower, "crp") ||
		   strings.Contains(parameterLower, "血沉") || strings.Contains(parameterLower, "esr") ||
		   strings.Contains(parameterLower, "白细胞") || strings.Contains(parameterLower, "wbc") {
			analysis.InflammatoryMarkers = append(analysis.InflammatoryMarkers, result)
			if result.Status == "偏高" {
				inflammationPresent = true
			}
		}

		// 肾功能指标分析
		if strings.Contains(parameterLower, "肌酐") || strings.Contains(parameterLower, "creatinine") ||
		   strings.Contains(parameterLower, "尿素") || strings.Contains(parameterLower, "urea") ||
		   strings.Contains(parameterLower, "肾小球") || strings.Contains(parameterLower, "gfr") {
			analysis.KidneyFunction = append(analysis.KidneyFunction, result)
			if result.Status == "偏高" || (strings.Contains(parameterLower, "gfr") && result.Status == "偏低") {
				kidneyIssues = true
			}
		}
	}

	// 综合评估风险等级
	if uricAcidHigh && inflammationPresent && kidneyIssues {
		analysis.RiskLevel = "高风险"
		analysis.FollowUpNeeded = true
		analysis.Recommendations = append(analysis.Recommendations, "存在多项异常指标，强烈建议立即就医，需要专业医生制定治疗方案")
	} else if uricAcidHigh && (inflammationPresent || kidneyIssues) {
		analysis.RiskLevel = "中风险"
		analysis.FollowUpNeeded = true
		analysis.Recommendations = append(analysis.Recommendations, "建议尽快就医，进行进一步检查和评估")
	} else if uricAcidHigh {
		analysis.RiskLevel = "低风险"
		analysis.FollowUpNeeded = true
		analysis.Recommendations = append(analysis.Recommendations, "建议调整生活方式，定期复查")
	} else {
		analysis.RiskLevel = "低风险"
		analysis.FollowUpNeeded = false
		analysis.Recommendations = append(analysis.Recommendations, "各项指标基本正常，保持健康的生活方式")
	}

	// 通用建议
	if uricAcidHigh || inflammationPresent {
		analysis.Recommendations = append(analysis.Recommendations,
			"建议低嘌呤饮食：避免内脏、海鲜、浓汤等高嘌呤食物",
			"增加饮水量，每日至少2000ml",
			"限制酒精摄入，特别是啤酒",
			"适量运动，避免剧烈运动",
			"控制体重，避免肥胖")
	}

	if kidneyIssues {
		analysis.Recommendations = append(analysis.Recommendations,
			"注意保护肾功能，避免使用肾毒性药物",
			"控制血压和血糖",
			"定期监测肾功能指标")
	}

	return analysis
}