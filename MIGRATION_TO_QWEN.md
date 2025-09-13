# 迁移到阿里百炼 Qwen-plus 模型

## 概述

本项目已成功从 OpenAI GPT 模型迁移到阿里百炼 Qwen-plus 模型。这次迁移通过阿里百炼提供的 OpenAI 兼容 API 实现，保持了代码的最小化修改。

## 主要变更

### 1. 环境变量更改
- **旧配置**: `OPENAI_API_KEY`
- **新配置**: `DASHSCOPE_API_KEY`

### 2. LLM 初始化配置
**修改前**:
```go
llm, err := openai.New()
```

**修改后**:
```go
llm, err := openai.New(
    openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
    openai.WithModel("qwen-plus"),
    openai.WithToken(os.Getenv("DASHSCOPE_API_KEY")),
)
```

### 3. 修改的文件

#### 核心代码文件
- `main.go`: 
  - 更新 LLM 初始化配置
  - 添加环境变量检查
  - 增加命令行参数支持
  - 更新错误信息
  
- `demo.go`:
  - 更新 LLM 初始化配置
  - 修改环境变量检查
  - 更新标题说明

#### 测试和工具文件
- `test_runner.go`: 重命名 main 函数避免冲突
- `test.go`: 无需修改（仅测试工具函数）

#### 文档文件
- `README.md`: 
  - 更新项目描述
  - 修改技术栈说明
  - 更新环境要求
  - 重写 API Key 配置步骤

### 4. 新增功能

#### 命令行参数支持
```bash
# 交互式对话模式 (默认)
go run *.go

# 运行演示模式
go run *.go demo

# 运行测试模式
go run *.go test

# 运行简单示例
go run *.go example

# 显示帮助信息
go run *.go help
```

## 技术实现细节

### OpenAI 兼容性
项目通过阿里百炼的 OpenAI 兼容 API 实现迁移，具有以下优势：
- 最小化代码修改
- 保持原有的 LangChain Go 框架兼容性
- 无需重写核心业务逻辑

### API 端点配置
- **Base URL**: `https://dashscope.aliyuncs.com/compatible-mode/v1`
- **模型**: `qwen-plus`
- **认证**: 通过 `DASHSCOPE_API_KEY` 环境变量

## 使用指南

### 1. 获取 API Key
1. 访问 [阿里云百炼平台](https://dashscope.aliyun.com/)
2. 登录或注册阿里云账号
3. 在控制台中创建 API Key
4. 复制 API Key 并设置到环境变量

### 2. 配置环境变量
```bash
export DASHSCOPE_API_KEY="your-dashscope-api-key"
```

### 3. 运行程序
```bash
# 安装依赖
go mod tidy

# 编译程序
go build -o gout-analysis-agent *.go

# 运行交互式模式
./gout-analysis-agent

# 运行演示模式
./gout-analysis-agent demo

# 查看帮助
./gout-analysis-agent help
```

## 功能验证

### 测试结果
所有核心功能已通过测试验证：
- ✅ 化验单数据智能解析
- ✅ 风险评估和医学建议
- ✅ 专业医学知识咨询
- ✅ 对话记忆和上下文理解
- ✅ 数值计算和分析能力

### 测试方式
```bash
# 运行功能测试
./gout-analysis-agent test

# 运行演示模式
./gout-analysis-agent demo
```

## 性能比较

### 模型特点
- **Qwen-plus**: 阿里巴巴自研的千亿参数级大语言模型
- **中文优化**: 对中文理解和生成能力更强
- **医学知识**: 具备较好的医学领域知识
- **成本效益**: 相比 OpenAI 具有更好的性价比

### 兼容性
- 完全兼容原有的 LangChain Go 框架
- 保持所有工具函数不变
- 支持所有原有功能

## 注意事项

### 1. API 限制
- 请遵守阿里云百炼平台的使用条款
- 注意 API 调用频率限制
- 合理使用以避免超出配额

### 2. 安全性
- 妥善保管 API Key，不要提交到代码仓库
- 定期轮换 API Key
- 监控 API 使用情况

### 3. 功能限制
- 本智能体仅供教育和参考用途
- 不能替代专业医疗诊断和治疗建议
- 请务必咨询专业医生

## 回滚指南

如需回滚到 OpenAI 模型，可以：

1. 还原环境变量：
   ```bash
   export OPENAI_API_KEY="your-openai-api-key"
   ```

2. 还原 LLM 初始化：
   ```go
   llm, err := openai.New()
   ```

3. 更新相关错误信息和文档

## 支持与反馈

如有问题或建议，请通过以下方式联系：
- [GitHub Issues](https://github.com/tmc/langchaingo/issues)
- [阿里云百炼文档](https://help.aliyun.com/zh/model-studio/)

---

**修改时间**: 2025-08-31  
**修改者**: Qoder AI Assistant  
**版本**: v1.0 (Qwen-plus)