package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/serpapi"

	"github.com/dean2027/aishell/pkg/prompt"
	localtools "github.com/dean2027/aishell/pkg/tools"
)

// 配置常量
const (
	// ConversationBufferSize 对话窗口缓冲大小，控制保持的对话轮数
	ConversationBufferSize = 100

	// MaxExecutorIterations 执行器最大迭代次数，控制单次交互的推理深度
	MaxExecutorIterations = 30
)

type ChatBot struct {
	executor *agents.Executor
	memory   *memory.ConversationWindowBuffer
	llm      llms.Model
	ctx      context.Context
}

// 初始化聊天机器人
func NewChatBot(ctx context.Context) (*ChatBot, error) {
	// 初始化OpenAI LLM
	llm, err := openai.New()
	if err != nil {
		return nil, fmt.Errorf("初始化LLM失败: %w", err)
	}

	// 初始化对话窗口缓冲内存 (保持最近N轮对话)
	conversationMemory := memory.NewConversationWindowBuffer(ConversationBufferSize)

	// 创建工具列表
	toolsList := []tools.Tool{
		tools.Calculator{},
		localtools.NewSystemCommand(),
		localtools.NewFileReader(),
		localtools.NewFileWriter(),
	}

	// 如果设置了SERPAPI_API_KEY，添加搜索工具
	if os.Getenv("SERPAPI_API_KEY") != "" {
		searchTool, err := serpapi.New()
		if err == nil {
			toolsList = append(toolsList, searchTool)
		}
	}

	// 创建智能终端助手的专用系统提示
	systemPromptPrefix := prompt.CreateSystemPrompt()

	// 创建使用内存和自定义系统提示的对话代理
	agent := agents.NewConversationalAgent(llm, toolsList,
		agents.WithMemory(conversationMemory),
		agents.WithPromptPrefix(systemPromptPrefix),
	)

	// 根据环境变量决定是否启用调试模式
	var executorOptions []agents.Option
	executorOptions = append(executorOptions, agents.WithMaxIterations(MaxExecutorIterations))
	// 关键修复：Executor 也需要同样的 memory 实例
	executorOptions = append(executorOptions, agents.WithMemory(conversationMemory))

	if os.Getenv("AISHELL_DEBUG") == "true" {
		fmt.Println("🔍 调试模式已启用 - 将显示详细的执行日志")
		debugHandler := callbacks.LogHandler{}
		executorOptions = append(executorOptions, agents.WithCallbacksHandler(debugHandler))
	}

	executor := agents.NewExecutor(agent, executorOptions...)

	return &ChatBot{
		executor: executor,
		memory:   conversationMemory,
		llm:      llm,
		ctx:      ctx,
	}, nil
}

// 处理用户输入
func (cb *ChatBot) ProcessInput(input string) (string, error) {
	// 调用执行器处理输入
	result, err := chains.Run(cb.ctx, cb.executor, input)
	if err != nil {
		// ConversationalAgent 现在应该足够稳定，直接返回错误
		// 如果频繁出现解析错误，可以考虑重新启用 fallback 机制
		return "", fmt.Errorf("处理输入失败: %w", err)
	}

	return result, nil
}

// 打印欢迎信息
func printWelcome() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Println("🤖 AI Shell - 智能终端助手")
	cyan.Println("============================")

	// 身份介绍
	fmt.Println("👨‍💻 我是您的智能终端助手，专门帮助您解决各种系统和技术问题")
	fmt.Println()

	yellow.Println("🎯 我能为您做什么:")
	fmt.Println("  🔧 系统管理 - 软件安装、配置查看、文件操作")
	fmt.Println("  🧮 数据计算 - 数学运算、数据分析、公式求解")
	fmt.Println("  🔍 信息搜索 - 技术文档、解决方案、最新资讯")
	fmt.Println("  💡 问题诊断 - 系统问题分析、性能优化建议")
	fmt.Println("  📝 代码协助 - 代码分析、开发环境配置")
	fmt.Println("  📄 文件读取 - 按行号范围读取文件内容")
	fmt.Println("  📝 文件写入 - 创建和编辑文本文件")
	fmt.Println()

	yellow.Println("💬 交互方式:")
	fmt.Println("  • 用自然语言描述您的需求，我会智能选择最合适的工具")
	fmt.Println("  • 支持 ↑↓ 浏览历史，Tab 自动补全，Ctrl+R 搜索历史")
	fmt.Println("  • 输入 'exit' 退出 | 'help' 查看功能")
	fmt.Println("")

	if os.Getenv("OPENAI_API_KEY") == "" {
		color.Red("⚠️  警告: 未设置OPENAI_API_KEY环境变量")
		fmt.Println("   请设置: export OPENAI_API_KEY=your_api_key")
		fmt.Println("")
	}

	if os.Getenv("SERPAPI_API_KEY") == "" {
		color.Yellow("💡 提示: 设置SERPAPI_API_KEY可启用网络搜索功能")
		fmt.Println("")
	}

	if os.Getenv("AISHELL_DEBUG") == "true" {
		color.Green("🔍 调试模式已启用 - 将显示详细的执行日志")
		fmt.Println("")
	} else {
		color.Yellow("💡 提示: 设置AISHELL_DEBUG=true可启用详细调试输出")
		fmt.Println("")
	}
}

// 打印帮助信息
func printHelp() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	green := color.New(color.FgGreen)

	cyan.Println("🤖 智能终端助手 - 功能说明")
	cyan.Println("===============================")

	yellow.Println("🔧 系统管理功能:")
	green.Println("  • 软件安装: '帮我安装Python', '安装nodejs'")
	green.Println("  • 系统信息: '查看系统配置', '检查磁盘空间'")
	green.Println("  • 文件操作: '创建项目目录', '查看当前文件'")
	green.Println("  • 进程管理: '查看运行的服务', '检查端口占用'")
	fmt.Println()

	yellow.Println("📄 文件读取功能:")
	green.Println("  • 读取完整文件: '帮我读取main.go', '查看config.json文件'")
	green.Println("  • 按行号范围: '读取main.go的前10行', '查看第20-30行'")
	green.Println("  • 支持相对和绝对路径: '/path/to/file', './src/main.go'")
	fmt.Println()

	yellow.Println("📝 文件写入功能:")
	green.Println("  • 创建新文件: '创建一个config.txt文件', '写入Hello World到test.txt'")
	green.Println("  • 编辑现有文件: '更新main.go中的代码', '修改配置文件'")
	green.Println("  • 自动创建目录: '在新目录中创建文件', '创建完整的目录结构'")
	green.Println("  • 支持多种文本格式: .txt, .go, .py, .js, .json, .md等")
	fmt.Println()

	yellow.Println("🧮 计算分析功能:")
	green.Println("  • 数学计算: '计算 (15 + 25) * 2', '求解方程'")
	green.Println("  • 数据处理: '分析这组数据的统计特征'")
	green.Println("  • 单位转换: '1GB等于多少MB'")
	fmt.Println()

	if os.Getenv("SERPAPI_API_KEY") != "" {
		yellow.Println("🔍 信息搜索功能:")
		green.Println("  • 技术搜索: '搜索Go语言最佳实践'")
		green.Println("  • 问题解决: '查找Redis连接错误的解决方案'")
		green.Println("  • 资讯获取: '最新的Docker更新内容'")
		fmt.Println()
	}

	yellow.Println("💡 智能诊断功能:")
	green.Println("  • 问题分析: '分析系统性能瓶颈'")
	green.Println("  • 优化建议: '如何提升服务器性能'")
	green.Println("  • 故障排查: '为什么我的应用启动失败'")
	fmt.Println()

	yellow.Println("⌨️  快捷键:")
	green.Println("  • ↑↓ 方向键 - 浏览历史命令")
	green.Println("  • Tab 键 - 自动补全命令")
	green.Println("  • Ctrl+R - 搜索历史命令")
	green.Println("  • Ctrl+C - 中断当前输入")
	green.Println("  • Ctrl+D 或 'exit' - 退出程序")
	fmt.Println()

	yellow.Println("💡 使用技巧:")
	fmt.Println("  • 用自然语言描述您的需求，无需记忆复杂命令")
	fmt.Println("  • 我会根据您的操作系统自动适配命令")
	fmt.Println("  • 告诉我您的工作背景，我能提供更精准的帮助")
	fmt.Println()
}

// 创建自动补全器
func createCompleter() *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		// 系统管理命令
		readline.PcItem("帮我安装Python"),
		readline.PcItem("帮我安装nodejs"),
		readline.PcItem("查看系统配置"),
		readline.PcItem("检查磁盘空间"),
		readline.PcItem("创建项目目录"),
		readline.PcItem("查看当前文件"),
		readline.PcItem("查看运行的服务"),
		readline.PcItem("检查端口占用"),

		// 计算分析命令
		readline.PcItem("计算"),
		readline.PcItem("分析数据"),
		readline.PcItem("转换单位"),

		// 文件读取命令
		readline.PcItem("帮我读取main.go"),
		readline.PcItem("读取main.go的前10行"),
		readline.PcItem("查看config文件"),
		readline.PcItem("读取第20-30行"),

		// 文件写入命令
		readline.PcItem("创建一个config.txt文件"),
		readline.PcItem("写入Hello World到test.txt"),
		readline.PcItem("更新main.go文件"),
		readline.PcItem("创建新的代码文件"),

		// 信息搜索命令
		readline.PcItem("搜索Go语言最佳实践"),
		readline.PcItem("查找解决方案"),
		readline.PcItem("最新技术动态"),

		// 诊断优化命令
		readline.PcItem("分析系统性能"),
		readline.PcItem("优化建议"),
		readline.PcItem("故障排查"),

		// 系统命令
		readline.PcItem("help"),
		readline.PcItem("帮助"),
		readline.PcItem("history"),
		readline.PcItem("命令历史"),
		readline.PcItem("exit"),
		readline.PcItem("quit"),
	)
}

func main() {
	ctx := context.Background()

	printWelcome()

	// 初始化聊天机器人
	chatBot, err := NewChatBot(ctx)
	if err != nil {
		log.Fatal("初始化聊天机器人失败:", err)
	}

	// 配置 readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "💻 智能终端> ",
		HistoryFile:     "/tmp/aishell_history",
		AutoComplete:    createCompleter(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		log.Fatal("初始化readline失败:", err)
	}
	defer rl.Close()

	// 设置颜色
	blue := color.New(color.FgBlue)
	red := color.New(color.FgRed)

	// 欢迎提示
	fmt.Println("💡 提示: 使用 ↑↓ 浏览历史，Tab 键自动补全，Ctrl+C 中断")
	fmt.Println()

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		// 处理特殊命令
		switch strings.ToLower(input) {
		case "exit", "quit":
			blue.Println("👋 再见！感谢使用智能终端助手，祝您工作顺利！")
			return
		case "help", "帮助":
			printHelp()
			continue
		case "history", "命令历史":
			printCommandHistory(rl)
			continue
		}

		// 处理用户输入
		fmt.Print("\n🤔 思考中...")
		response, err := chatBot.ProcessInput(input)
		fmt.Print("\r                    \r") // 清除"思考中"提示

		if err != nil {
			red.Printf("❌ 错误: %v\n\n", err)
			continue
		}

		// 显示回复
		blue.Println("🤖 终端助手:")
		fmt.Println(response)
		fmt.Println("")
	}
}

// filterInput 过滤输入字符
func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

// printCommandHistory 显示命令历史
func printCommandHistory(rl *readline.Instance) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Println("📜 命令历史")
	cyan.Println("==========")

	// 获取历史记录 (readline 库的历史记录功能)
	yellow.Println("💡 使用 ↑↓ 方向键浏览历史命令")
	yellow.Println("💡 使用 Ctrl+R 进行历史搜索")
	yellow.Println("💡 历史记录已保存到 /tmp/aishell_history")

	fmt.Println()
}
