# 痛风化验单分析智能体 🩺

基于 LangChain Go 框架开发的专业痛风化验单分析智能体，使用阿里百炼 Qwen-plus 大模型，能够智能分析血液生化检查结果，评估痛风风险，并提供个性化的医学建议。

## ✨ 功能特性

### 🔬 智能化验单分析
- **自动解析** 多种格式的化验单数据
- **风险评估** 基于尿酸水平、炎症指标、肾功能等综合评估
- **分级诊断** 提供低风险、中风险、高风险的分级评估
- **个性化建议** 根据检查结果给出针对性的治疗和生活建议

### 📚 专业医学知识库
- **痛风知识** 完整的痛风疾病知识体系
- **诊断标准** 权威的医学诊断标准和参考值
- **治疗指南** 基于循证医学的治疗建议
- **预防措施** 科学的预防和生活方式指导

### 💬 智能对话交互
- **自然语言** 支持中文自然语言交互
- **上下文理解** 具备对话记忆和上下文理解能力
- **多轮对话** 支持连续提问和深入讨论
- **实时分析** 即时响应和分析处理

## 🛠️ 技术架构

### 核心组件
```
痛风分析智能体
├── GoutLabAnalyzer (化验单分析工具)
│   ├── 数据解析引擎
│   ├── 风险评估算法
│   └── 建议生成系统
├── MedicalKnowledgeBase (医学知识库)
│   ├── 痛风疾病知识
│   ├── 诊断标准库
│   └── 治疗指南库
└── ConversationalAgent (对话智能体)
    ├── LLM 语言模型
    ├── 工具调度器
    └── 记忆管理器
```

### 技术栈
- **框架**: LangChain Go v0.1.13
- **语言模型**: 阿里百炼 Qwen-plus
- **编程语言**: Go 1.24.3
- **架构模式**: Agent + Tools + Memory

## 📋 支持的检查项目

### 🩸 核心指标
- **尿酸 (Uric Acid)** - 痛风诊断的关键指标
- **C反应蛋白 (CRP)** - 炎症反应指标
- **血沉 (ESR)** - 炎症活动度
- **白细胞计数 (WBC)** - 感染和炎症指标

### 🫘 肾功能指标  
- **血肌酐 (Creatinine)** - 肾功能评估
- **血尿素氮 (BUN)** - 肾脏代谢功能
- **肾小球滤过率 (GFR)** - 肾功能综合评估

### 📊 参考标准
| 检查项目 | 正常范围 | 异常提示 |
|---------|---------|---------|
| 尿酸 (男性) | 208-428 μmol/L | >420 μmol/L 高尿酸血症 |
| 尿酸 (女性) | 155-357 μmol/L | >360 μmol/L 高尿酸血症 |
| CRP | <3.0 mg/L | >3.0 mg/L 炎症反应 |
| 血沉 (男性) | <15 mm/h | >15 mm/h 炎症活动 |
| 血沉 (女性) | <20 mm/h | >20 mm/h 炎症活动 |

## 🚀 快速开始

### 环境要求
- Go 1.23.0+
- 阿里百炼 API Key (从阿里云百炼平台获取)

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/tmc/langchaingo
   cd langchaingo/examples/gout-analysis-agent
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **配置阿里百炼 API Key**
   ```bash
   export DASHSCOPE_API_KEY="your-dashscope-api-key"
   ```
   
   > 📝 **获取 API Key 步骤**:
   > 1. 访问 [阿里云百炼平台](https://dashscope.aliyun.com/)
   > 2. 登录或注册阿里云账号
   > 3. 在控制台中创建 API Key
   > 4. 复制 API Key 并设置到环境变量

4. **运行智能体**
   ```bash
   go run *.go
   ```

## 💡 使用示例

### 化验单分析
```
输入化验单数据：
尿酸 520 umol/L (参考范围: 208-428)
C反应蛋白 15.2 mg/L (参考范围: <3.0)
血沉 45 mm/h (参考范围: <15)
肌酐 95 umol/L (参考范围: 54-106)
```

### 智能体分析结果
```json
{
  "uric_acid_level": {
    "parameter": "尿酸",
    "value": 520,
    "unit": "umol/L",
    "reference_max": 428,
    "status": "偏高"
  },
  "risk_level": "中风险",
  "recommendations": [
    "尿酸水平显著升高，建议立即就医，考虑药物治疗",
    "建议低嘌呤饮食：避免内脏、海鲜、浓汤等高嘌呤食物",
    "增加饮水量，每日至少2000ml"
  ],
  "follow_up_needed": true
}
```

### 医学知识查询
```
Q: 什么是痛风？
A: 痛风是一种由于嘌呤代谢紊乱和/或尿酸排泄减少所致的高尿酸血症直接相关的代谢性疾病，以反复发作的急性关节炎、痛风石形成、慢性关节炎和关节畸形为特征...
```

## 📝 输入数据格式

### 支持的格式

1. **标准格式**
   ```
   尿酸 520 umol/L (参考范围: 208-428)
   C反应蛋白 15.2 mg/L (<3.0)
   ```

2. **简化格式**  
   ```
   尿酸: 520 umol/L
   CRP: 15.2 mg/L
   ```

3. **多行格式**
   ```
   检查项目: 尿酸
   检查结果: 520 umol/L
   参考范围: 208-428 umol/L
   ```

### 识别关键词
- **尿酸**: uric acid, 尿酸, UA
- **C反应蛋白**: CRP, C-reactive protein, C反应蛋白
- **血沉**: ESR, 红细胞沉降率, 血沉
- **肌酐**: creatinine, Cr, 肌酐

## ⚠️ 重要声明

> **医疗免责声明**
> 
> 本智能体仅供教育和参考用途，不能替代专业医疗诊断和治疗建议。
> - ❌ 不可作为临床诊断依据
> - ❌ 不可替代医生专业判断  
> - ❌ 不可用于紧急医疗决策
> - ✅ 请务必咨询专业医生
> - ✅ 定期进行正规体检
> - ✅ 遵循医生治疗方案

## 🧩 扩展开发

### 添加新的检查项目
```go
// 在 gout_analyzer.go 中添加新的检测项目
if strings.Contains(parameterLower, "新检查项目") {
    // 添加分析逻辑
}
```

### 扩展医学知识
```go
// 在 medical_knowledge.go 中添加新知识
mkb.knowledge["新疾病"] = MedicalInfo{
    Topic: "新疾病名称",
    Definition: "疾病定义...",
    // 其他医学信息
}
```

### 自定义风险评估算法
```go
// 在 analyzeGoutRisk 方法中自定义评估逻辑
func (g *GoutLabAnalyzer) analyzeGoutRisk(results []LabResult) GoutAnalysisResult {
    // 自定义风险评估算法
}
```

## 🤝 贡献指南

欢迎贡献代码和改进建议！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系支持

- **项目主页**: [LangChain Go](https://github.com/tmc/langchaingo)
- **问题反馈**: [GitHub Issues](https://github.com/tmc/langchaingo/issues)
- **技术文档**: [LangChain Go Docs](https://pkg.go.dev/github.com/tmc/langchaingo)

---

🩺 **让AI助力健康管理，科技赋能医疗服务！**