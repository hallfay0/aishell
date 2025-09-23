package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// WelcomeInfo 欢迎信息配置
type WelcomeInfo struct {
	Title       string
	Description string
	ShowTips    bool
}

// PrintWelcome 打印欢迎信息
func PrintWelcome() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Println("🤖 AI Shell - 智能终端助手")
	cyan.Println("============================")

	// 身份介绍
	fmt.Println("👨‍💻 我是您的智能终端助手，专门帮助您解决各种系统和技术问题")
	fmt.Println()

	yellow.Println("💬 交互方式:")
	fmt.Println("  • 用自然语言描述您的需求，我会智能选择最合适的工具")
	fmt.Println("  • 支持 ↑↓ 浏览历史，Tab 自动补全，Ctrl+R 搜索历史")
	fmt.Println("  • 输入 'exit' 退出 | 'help' 查看功能")
	fmt.Println("")

	printEnvironmentStatus()
}

// printEnvironmentStatus 打印环境状态信息
func printEnvironmentStatus() {
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

// PrintHelp 打印帮助信息
func PrintHelp() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	green := color.New(color.FgGreen)

	cyan.Println("🤖 智能终端助手 - 功能说明")
	cyan.Println("===============================")

	printSystemFeatures(yellow, green)
	printFileFeatures(yellow, green)
	printCalculationFeatures(yellow, green)
	printSearchFeatures(yellow, green)
	printDiagnosticFeatures(yellow, green)
	printShortcuts(yellow, green)
	printTips()
}

// printSystemFeatures 打印系统管理功能
func printSystemFeatures(yellow, green *color.Color) {
	yellow.Println("🔧 系统管理功能:")
	green.Println("  • 软件安装: '帮我安装Python', '安装nodejs'")
	green.Println("  • 系统信息: '查看系统配置', '检查磁盘空间'")
	green.Println("  • 文件操作: '创建项目目录', '查看当前文件'")
	green.Println("  • 进程管理: '查看运行的服务', '检查端口占用'")
	fmt.Println()
}

// printFileFeatures 打印文件操作功能
func printFileFeatures(yellow, green *color.Color) {
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
}

// printCalculationFeatures 打印计算分析功能
func printCalculationFeatures(yellow, green *color.Color) {
	yellow.Println("🧮 计算分析功能:")
	green.Println("  • 数学计算: '计算 (15 + 25) * 2', '求解方程'")
	green.Println("  • 数据处理: '分析这组数据的统计特征'")
	green.Println("  • 单位转换: '1GB等于多少MB'")
	fmt.Println()
}

// printSearchFeatures 打印搜索功能
func printSearchFeatures(yellow, green *color.Color) {
	if os.Getenv("SERPAPI_API_KEY") != "" {
		yellow.Println("🔍 信息搜索功能:")
		green.Println("  • 技术搜索: '搜索Go语言最佳实践'")
		green.Println("  • 问题解决: '查找Redis连接错误的解决方案'")
		green.Println("  • 资讯获取: '最新的Docker更新内容'")
		fmt.Println()
	}
}

// printDiagnosticFeatures 打印诊断功能
func printDiagnosticFeatures(yellow, green *color.Color) {
	yellow.Println("💡 智能诊断功能:")
	green.Println("  • 问题分析: '分析系统性能瓶颈'")
	green.Println("  • 优化建议: '如何提升服务器性能'")
	green.Println("  • 故障排查: '为什么我的应用启动失败'")
	fmt.Println()
}

// printShortcuts 打印快捷键
func printShortcuts(yellow, green *color.Color) {
	yellow.Println("⌨️  快捷键:")
	green.Println("  • ↑↓ 方向键 - 浏览历史命令")
	green.Println("  • Tab 键 - 自动补全命令")
	green.Println("  • Ctrl+R - 搜索历史命令")
	green.Println("  • Ctrl+C - 中断当前输入")
	green.Println("  • Ctrl+D 或 'exit' - 退出程序")
	fmt.Println()
}

// printTips 打印使用技巧
func printTips() {
	yellow := color.New(color.FgYellow, color.Bold)
	
	yellow.Println("💡 使用技巧:")
	fmt.Println("  • 用自然语言描述您的需求，无需记忆复杂命令")
	fmt.Println("  • 我会根据您的操作系统自动适配命令")
	fmt.Println("  • 告诉我您的工作背景，我能提供更精准的帮助")
	fmt.Println()
}

// PrintGoodbye 打印告别信息
func PrintGoodbye() {
	blue := color.New(color.FgBlue)
	blue.Println("👋 再见！感谢使用智能终端助手，祝您工作顺利！")
}

// PrintError 打印错误信息
func PrintError(msg string, err error) {
	red := color.New(color.FgRed)
	red.Printf("❌ %s: %v\n\n", msg, err)
}

// PrintThinking 打印思考状态
func PrintThinking() {
	fmt.Print("\n🤔 思考中...")
}

// ClearThinking 清除思考状态
func ClearThinking() {
	fmt.Print("\r                    \r") // 清除"思考中"提示
}

// PrintResponse 打印AI响应
func PrintResponse(response string) {
	blue := color.New(color.FgBlue)
	blue.Println("🤖 终端助手:")
	fmt.Println(response)
	fmt.Println("")
}
