package ui

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

// PrintCommandHistory 显示命令历史
func PrintCommandHistory(rl *readline.Instance) {
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

// PrintUsageTips 打印使用提示
func PrintUsageTips() {
	fmt.Println("💡 提示: 使用 ↑↓ 浏览历史，Tab 键自动补全，Ctrl+C 中断")
	fmt.Println()
}
