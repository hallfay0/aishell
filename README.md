# 🤖 AI Shell - 智能终端助手

一个基于Go语言和LangChain的智能终端助手，集成OpenAI GPT，具备记忆功能和多种实用工具，让命令行操作变得更加智能和高效。

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ 功能特性

### 🧠 智能对话
- **上下文记忆**: 基于滑动窗口的对话记忆，记住最近100轮对话
- **环境感知**: 自动识别操作系统、架构、当前目录等环境信息
- **自然语言交互**: 用自然语言描述需求，AI智能选择合适的工具

### 🛠️ 强大工具集
- **🔧 系统命令**: 跨平台系统命令执行，安全检查，危险命令确认
- **📄 文件读取**: 按行号范围读取文件内容，支持大文件处理
- **📝 文件写入**: 创建和编辑文本文件，自动创建目录结构
- **🧮 数学计算**: 复杂数学运算和数据分析
- **🔍 网络搜索**: 集成SerpAPI的实时信息搜索（可选）

### 🎛️ 用户体验
- **智能补全**: Tab键自动补全，支持历史命令
- **历史记录**: ↑↓键浏览历史，Ctrl+R搜索历史
- **调试模式**: 详细的执行日志，便于开发调试
- **彩色输出**: 美观的界面和清晰的信息层级

## 📦 安装

### 系统要求
- Go 1.24+ 
- OpenAI API Key

### 快速安装

```bash
# 克隆项目
git clone https://github.com/dean2027/aishell.git
cd aishell

# 安装依赖
go mod tidy

# 编译项目
go build

# 运行
./aishell
```

## ⚙️ 配置

### 环境变量配置

创建环境配置文件：

```bash
cp env.example .env
```

编辑 `.env` 文件：

```bash
# OpenAI API 密钥 (必需)
export OPENAI_API_KEY="your_openai_api_key_here"

# SerpAPI 密钥 (可选，用于网络搜索)
export SERPAPI_API_KEY="your_serpapi_key_here"

# 调试模式 (可选)
export AISHELL_DEBUG="true"
```

加载配置：

```bash
source .env
```

### 配置参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `ConversationBufferSize` | 100 | 对话记忆窗口大小 |
| `MaxExecutorIterations` | 30 | 最大推理迭代次数 |
| `AISHELL_DEBUG` | false | 调试模式开关 |

## 🚀 使用方法

### 基本使用

```bash
# 启动AI Shell
./aishell

# 启用调试模式
AISHELL_DEBUG=true ./aishell
```

### 使用示例

#### 系统管理
```bash
💻 智能终端> 帮我查看系统信息
💻 智能终端> 检查磁盘使用情况
💻 智能终端> 安装Python最新版本
```

#### 文件操作
```bash
💻 智能终端> 读取main.go的前50行
💻 智能终端> 创建一个config.json配置文件
💻 智能终端> 在src目录中创建一个新的Go文件
```

#### 计算分析
```bash
💻 智能终端> 计算 (15 + 25) * 3 / 2
💻 智能终端> 1TB等于多少GB
```

#### 信息搜索
```bash
💻 智能终端> 搜索Go语言最佳实践
💻 智能终端> 查找Docker容器优化方案
```

## 🏗️ 项目结构

```
aishell/
├── main.go                 # 主程序入口
├── go.mod                  # Go模块定义
├── go.sum                  # 依赖锁定文件
├── .gitignore              # Git忽略文件
├── env.example             # 环境配置模板
├── README.md               # 项目文档
└── pkg/                    # 核心包
    ├── prompt/             # 系统提示模块
    │   └── system_prompt.go
    ├── tools/              # 工具模块
    │   ├── file_reader.go      # 文件读取工具
    │   ├── file_writer.go      # 文件写入工具
    │   ├── system_command.go   # 系统命令工具
    │   └── *_test.go           # 单元测试
    └── utils/              # 工具函数
        └── environment.go      # 环境信息
```

## 🔧 开发

### 本地开发

```bash
# 运行测试
go test ./...

# 运行特定测试
go test ./pkg/tools/... -v

# 代码格式化
go fmt ./...

# 静态检查
go vet ./...

# 构建
go build -o aishell
```

### 添加新工具

1. 在 `pkg/tools/` 目录创建新工具文件
2. 实现 `tools.Tool` 接口：
   ```go
   type Tool interface {
       Name() string
       Description() string
       Call(ctx context.Context, input string) (string, error)
   }
   ```
3. 在 `main.go` 中注册新工具
4. 添加相应的单元测试

### 测试覆盖率

```bash
# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 添加适当的注释和文档
- 确保所有测试通过
- 保持测试覆盖率在90%以上

## 🐛 问题排查

### 常见问题

**Q: 提示 "初始化LLM失败"**
A: 请检查 `OPENAI_API_KEY` 环境变量是否正确设置

**Q: 搜索功能不可用**
A: 需要设置 `SERPAPI_API_KEY` 环境变量启用搜索功能

**Q: 命令执行权限问题**
A: 危险命令会提示确认，输入 `y` 确认执行

**Q: 文件写入失败**
A: 检查目录权限，或设置 `create_dirs: true` 自动创建目录

### 调试模式

启用调试模式查看详细执行日志：

```bash
AISHELL_DEBUG=true ./aishell
```

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🙏 致谢

- [LangChain Go](https://github.com/tmc/langchaingo) - 强大的LLM框架
- [OpenAI](https://openai.com/) - GPT模型支持
- [Readline](https://github.com/chzyer/readline) - 交互式命令行
- [Color](https://github.com/fatih/color) - 彩色终端输出

## 📞 联系

- 项目地址: https://github.com/dean2027/aishell
- 问题反馈: [GitHub Issues](https://github.com/dean2027/aishell/issues)

---

⭐ 如果这个项目对您有帮助，请给个星星支持！