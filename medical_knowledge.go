package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
)

// MedicalKnowledgeBase 医学知识库工具
type MedicalKnowledgeBase struct {
	CallbacksHandler callbacks.Handler
	knowledge        map[string]MedicalInfo
}

// MedicalInfo 医学信息结构
type MedicalInfo struct {
	Topic       string   `json:"topic"`       // 主题
	Definition  string   `json:"definition"`  // 定义
	Symptoms    []string `json:"symptoms"`    // 症状
	Causes      []string `json:"causes"`      // 病因
	RiskFactors []string `json:"risk_factors"` // 危险因素
	Diagnosis   []string `json:"diagnosis"`   // 诊断标准
	Treatment   []string `json:"treatment"`   // 治疗方法
	Prevention  []string `json:"prevention"`  // 预防措施
	References  []string `json:"references"`  // 参考值/标准
}

// NewMedicalKnowledgeBase 创建医学知识库实例
func NewMedicalKnowledgeBase() *MedicalKnowledgeBase {
	mkb := &MedicalKnowledgeBase{
		knowledge: make(map[string]MedicalInfo),
	}
	mkb.initializeKnowledge()
	return mkb
}

// Name 返回工具名称
func (mkb MedicalKnowledgeBase) Name() string {
	return "medical_knowledge_base"
}

// Description 返回工具描述
func (mkb MedicalKnowledgeBase) Description() string {
	return `医学知识库工具。提供痛风及相关疾病的专业医学知识，包括定义、症状、诊断标准、治疗方法等。
支持查询的主题包括：
- 痛风 (gout)
- 高尿酸血症 (hyperuricemia)  
- 尿酸 (uric_acid)
- 痛风性关节炎 (gouty_arthritis)
- 痛风石 (tophi)
- 肾功能 (kidney_function)
- 炎症指标 (inflammation)
输入查询主题的关键词即可获得相关医学知识。`
}

// Call 执行知识查询
func (mkb MedicalKnowledgeBase) Call(ctx context.Context, input string) (string, error) {
	if mkb.CallbacksHandler != nil {
		mkb.CallbacksHandler.HandleToolStart(ctx, input)
	}

	// 查找相关知识
	results := mkb.searchKnowledge(strings.ToLower(strings.TrimSpace(input)))
	
	if len(results) == 0 {
		return "未找到相关医学知识。请尝试使用以下关键词：痛风、高尿酸血症、尿酸、关节炎、痛风石、肾功能、炎症等。", nil
	}

	// 格式化输出
	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Sprintf("格式化医学知识时出错: %v", err), nil
	}

	if mkb.CallbacksHandler != nil {
		mkb.CallbacksHandler.HandleToolEnd(ctx, string(output))
	}

	return string(output), nil
}

// searchKnowledge 搜索相关医学知识
func (mkb *MedicalKnowledgeBase) searchKnowledge(query string) []MedicalInfo {
	var results []MedicalInfo
	
	for key, info := range mkb.knowledge {
		// 检查查询词是否匹配主题关键词
		if strings.Contains(query, key) || strings.Contains(key, query) ||
		   strings.Contains(strings.ToLower(info.Topic), query) {
			results = append(results, info)
		}
	}
	
	return results
}

// initializeKnowledge 初始化医学知识库
func (mkb *MedicalKnowledgeBase) initializeKnowledge() {
	// 痛风知识
	mkb.knowledge["痛风"] = MedicalInfo{
		Topic: "痛风 (Gout)",
		Definition: "痛风是一种由于嘌呤代谢紊乱和/或尿酸排泄减少所致的高尿酸血症直接相关的代谢性疾病，以反复发作的急性关节炎、痛风石形成、慢性关节炎和关节畸形为特征。",
		Symptoms: []string{
			"急性关节疼痛，多在夜间突然发作",
			"关节红肿热痛，触痛明显",
			"常首发于第一跖趾关节（大脚趾）",
			"可累及踝关节、膝关节、手指关节等",
			"发热、寒战",
			"疼痛呈刀割样、撕裂样",
		},
		Causes: []string{
			"嘌呤代谢紊乱",
			"尿酸生成过多",
			"尿酸排泄减少",
			"遗传因素",
			"饮食因素（高嘌呤饮食）",
			"肥胖",
			"酒精摄入",
		},
		RiskFactors: []string{
			"男性，40岁以上",
			"绝经后女性",
			"肥胖",
			"高血压",
			"糖尿病",
			"肾功能不全",
			"家族史",
			"长期饮酒",
			"高嘌呤饮食",
		},
		Diagnosis: []string{
			"血尿酸>420μmol/L（男性）或>360μmol/L（女性）",
			"关节滑液中发现尿酸盐结晶",
			"急性关节炎典型临床表现",
			"秋水仙碱治疗有效",
			"影像学检查显示痛风石或骨质破坏",
		},
		Treatment: []string{
			"急性期：秋水仙碱、NSAIDs、糖皮质激素",
			"缓解期：别嘌醇、非布司他降尿酸治疗",
			"目标：血尿酸<360μmol/L",
			"有痛风石者：血尿酸<300μmol/L",
			"急性期避免使用降尿酸药物",
		},
		Prevention: []string{
			"低嘌呤饮食",
			"控制体重",
			"限制酒精摄入",
			"多饮水（每日>2000ml）",
			"避免剧烈运动",
			"规律用药",
			"定期监测血尿酸",
		},
	}

	// 高尿酸血症知识
	mkb.knowledge["高尿酸血症"] = MedicalInfo{
		Topic: "高尿酸血症 (Hyperuricemia)",
		Definition: "高尿酸血症是指在正常嘌呤饮食状态下，非同日两次空腹血尿酸水平男性>420μmol/L，女性>360μmol/L。",
		Symptoms: []string{
			"多数患者无明显症状",
			"可能出现疲劳",
			"关节不适",
			"部分患者可发展为痛风",
		},
		Causes: []string{
			"嘌呤合成过多",
			"嘌呤摄入过多",
			"尿酸排泄减少",
			"遗传性酶缺陷",
			"药物影响（利尿剂、阿司匹林等）",
		},
		RiskFactors: []string{
			"遗传因素",
			"高嘌呤饮食",
			"肥胖",
			"饮酒",
			"肾功能减退",
			"某些药物使用",
		},
		Diagnosis: []string{
			"男性血尿酸>420μmol/L",
			"女性血尿酸>360μmol/L",
			"需排除继发性因素",
		},
		Treatment: []string{
			"生活方式干预",
			"必要时药物治疗",
			"治疗目标：血尿酸<360μmol/L",
			"有并发症时：<300μmol/L",
		},
		Prevention: []string{
			"控制饮食",
			"适量运动",
			"控制体重",
			"限制酒精",
			"多饮水",
		},
	}

	// 尿酸参考值
	mkb.knowledge["尿酸"] = MedicalInfo{
		Topic: "尿酸 (Uric Acid)",
		Definition: "尿酸是嘌呤代谢的最终产物，主要通过肾脏排泄。",
		References: []string{
			"正常参考值：",
			"男性：208-428 μmol/L (3.5-7.2 mg/dL)",
			"女性：155-357 μmol/L (2.6-6.0 mg/dL)",
			"高尿酸血症诊断标准：",
			"男性：>420 μmol/L (7.0 mg/dL)",
			"女性：>360 μmol/L (6.0 mg/dL)",
			"痛风治疗目标：",
			"一般患者：<360 μmol/L (6.0 mg/dL)",
			"有痛风石患者：<300 μmol/L (5.0 mg/dL)",
		},
	}

	// 炎症指标知识
	mkb.knowledge["炎症"] = MedicalInfo{
		Topic: "炎症指标 (Inflammatory Markers)",
		Definition: "炎症指标是反映机体炎症反应程度的实验室检查指标。",
		References: []string{
			"C反应蛋白(CRP)：<3.0 mg/L",
			"血沉(ESR)：男性<15mm/h，女性<20mm/h",
			"白细胞计数(WBC)：4.0-10.0×10⁹/L",
			"中性粒细胞百分比：50-70%",
		},
		Diagnosis: []string{
			"急性炎症：CRP显著升高",
			"慢性炎症：轻度升高",
			"感染性疾病：白细胞升高",
			"痛风急性发作：CRP、ESR升高",
		},
	}

	// 肾功能知识
	mkb.knowledge["肾功能"] = MedicalInfo{
		Topic: "肾功能 (Kidney Function)",
		Definition: "肾功能是指肾脏清除代谢产物、维持水电解质平衡和酸碱平衡的能力。",
		References: []string{
			"血肌酐(Cr)：男性54-106μmol/L，女性44-97μmol/L",
			"血尿素氮(BUN)：2.5-7.1mmol/L",
			"估算肾小球滤过率(eGFR)：>90ml/min/1.73m²",
			"尿酸清除率：6.2-17.2ml/min",
		},
		Diagnosis: []string{
			"慢性肾病：eGFR<60ml/min/1.73m²持续3个月",
			"急性肾损伤：血肌酐升高>26.5μmol/L/48h",
			"肾功能不全：eGFR<60ml/min/1.73m²",
		},
		Treatment: []string{
			"保护肾功能",
			"控制血压、血糖",
			"避免肾毒性药物",
			"适当蛋白质摄入限制",
		},
	}

	// 痛风性关节炎
	mkb.knowledge["关节炎"] = MedicalInfo{
		Topic: "痛风性关节炎 (Gouty Arthritis)",
		Definition: "痛风性关节炎是由于尿酸盐结晶沉积在关节滑膜、软骨和其他组织中引起的炎症性关节病。",
		Symptoms: []string{
			"急性发作：关节剧烈疼痛",
			"红肿热痛",
			"活动受限",
			"多在夜间发作",
			"单关节受累多见",
		},
		Diagnosis: []string{
			"典型临床表现",
			"血尿酸升高",
			"关节液尿酸盐结晶",
			"秋水仙碱试验性治疗有效",
			"影像学检查",
		},
		Treatment: []string{
			"急性期抗炎治疗",
			"秋水仙碱",
			"NSAIDs",
			"糖皮质激素",
			"避免降尿酸治疗",
		},
	}

	// 痛风石
	mkb.knowledge["痛风石"] = MedicalInfo{
		Topic: "痛风石 (Tophi)",
		Definition: "痛风石是尿酸盐结晶在软组织中的沉积物，是慢性痛风的特征性表现。",
		Symptoms: []string{
			"关节周围结节",
			"皮下结节",
			"可破溃流出白色物质",
			"关节变形",
			"功能障碍",
		},
		Treatment: []string{
			"积极降尿酸治疗",
			"目标血尿酸<300μmol/L",
			"外科手术切除",
			"物理治疗",
		},
		Prevention: []string{
			"长期降尿酸治疗",
			"定期监测",
			"避免诱发因素",
		},
	}
}