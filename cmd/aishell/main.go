package main

import (
	"context"
	"log"
	"os"

	"github.com/dean2027/aishell/pkg/app"
	"github.com/dean2027/aishell/pkg/cli"
)

// 版本信息（构建时注入）
var (
	Version   = "unknown"
	Commit    = "unknown"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	// 创建上下文
	ctx := context.Background()

	// 加载配置
	config := app.LoadConfig()

	// 创建CLI运行器
	runner, err := cli.NewRunner(ctx, config)
	if err != nil {
		log.Fatal("初始化应用失败:", err)
	}

	// 运行应用
	if err := runner.Run(); err != nil {
		log.Fatal("运行应用失败:", err)
	}
}

// init 初始化函数
func init() {
	// 可以在这里添加初始化逻辑
	// 比如日志配置、信号处理等

	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 处理版本信息查询
	handleVersionFlag()
}

// handleVersionFlag 处理版本信息标志
func handleVersionFlag() {
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			printVersion()
			os.Exit(0)
		}
		if arg == "--help" || arg == "-h" {
			printHelp()
			os.Exit(0)
		}
	}
}

// printVersion 打印版本信息
func printVersion() {
	println("🤖 AI Shell - 智能终端助手")
	println("版本:", Version)
	println("提交:", Commit)
	println("构建时间:", BuildTime)
	println("Go版本:", GoVersion)
}

// printHelp 打印命令行帮助
func printHelp() {
	println("🤖 AI Shell - 智能终端助手")
	println("")
	println("用法:")
	println("  aishell [选项]")
	println("")
	println("选项:")
	println("  -h, --help     显示此帮助信息")
	println("  -v, --version  显示版本信息")
	println("")
	println("环境变量:")
	println("  OPENAI_API_KEY     OpenAI API密钥 (必需)")
	println("  SERPAPI_API_KEY    SerpAPI密钥 (可选，用于搜索功能)")
	println("  AISHELL_DEBUG      启用调试模式 (true/false)")
	println("")
	println("示例:")
	println("  export OPENAI_API_KEY=your_key")
	println("  aishell")
	println("")
	println("  AISHELL_DEBUG=true aishell")
}