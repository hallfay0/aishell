package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/tmc/langchaingo/callbacks"
)

// SystemCommand 是一个可以执行系统命令的工具
type SystemCommand struct {
	CallbacksHandler callbacks.Handler
	// Timeout 命令执行超时时间，默认30秒
	Timeout time.Duration
	// DangerousCommands 危险命令列表，需要用户确认才能执行
	DangerousCommands []string
}

// NewSystemCommand 创建一个新的系统命令工具
func NewSystemCommand() *SystemCommand {
	return &SystemCommand{
		Timeout: 30 * time.Second,
		// 危险命令列表，需要用户确认才能执行
		DangerousCommands: []string{
			// 文件删除命令
			"rm", "del", "erase", "rmdir", "rd",
			// 系统关机重启
			"shutdown", "reboot", "halt", "poweroff", "init",
			// 磁盘操作
			"dd", "fdisk", "mkfs", "format", "parted", "gdisk",
			// 权限修改
			"chmod", "chown", "chgrp", "icacls", "takeown",
			// 网络配置
			"iptables", "netsh", "route", "ifconfig", "ip",
			// 服务管理
			"systemctl", "service", "sc", "net", "kill", "killall", "taskkill",
			// 内核模块
			"modprobe", "rmmod", "insmod",
			// 压缩解压（可能覆盖文件）
			"tar", "unzip", "7z", "rar",
			// 系统配置修改
			"crontab", "at", "schtasks",
			// 用户管理
			"useradd", "userdel", "usermod", "passwd", "su", "sudo",
			// 软件安装/卸载（可能影响系统）
			"rpm", "dpkg", "msiexec",
		},
	}
}

// Name 返回工具名称
func (s *SystemCommand) Name() string {
	return "system_command"
}

// Description 返回工具描述
func (s *SystemCommand) Description() string {
	return `执行系统命令的工具。可以执行跨平台的系统命令，如包管理器安装软件、文件操作、系统信息查询等。
输入格式：要执行的完整命令，例如：
- Linux/macOS: "apt install python3", "brew install node", "ls -la"
- Windows: "choco install nodejs", "dir", "systeminfo"
安全机制：大部分命令可直接执行，危险命令(如rm删除、shutdown关机等)需要用户确认。`
}

// Call 执行系统命令
func (s *SystemCommand) Call(ctx context.Context, input string) (string, error) {
	if s.CallbacksHandler != nil {
		s.CallbacksHandler.HandleToolStart(ctx, input)
	}

	// 清理输入
	command := strings.TrimSpace(input)
	if command == "" {
		return "错误：命令不能为空", nil
	}

	// 解析命令
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "错误：无效的命令格式", nil
	}

	baseCommand := parts[0]

	// 安全检查：检查是否是危险命令
	if s.isDangerousCommand(baseCommand) {
		shouldExecute := s.askUserPermission(baseCommand)
		if !shouldExecute {
			return fmt.Sprintf("危险命令 '%s' 执行已被取消", baseCommand), nil
		}
		// 用户选择执行，显示警告信息
		fmt.Printf("\n⚠️  警告：正在执行危险命令: %s\n", baseCommand)
	}

	// 设置超时上下文
	if s.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.Timeout)
		defer cancel()
	}

	// 根据操作系统执行命令
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// Windows使用cmd /c执行命令
		cmd = exec.CommandContext(ctx, "cmd", "/c", command)
	default:
		// Linux/macOS使用sh -c执行命令
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()

	result := ""
	if err != nil {
		// 如果命令执行失败，返回错误信息和输出
		result = fmt.Sprintf("命令执行失败: %v\n输出: %s", err, string(output))
	} else {
		// 命令执行成功
		result = fmt.Sprintf("命令执行成功:\n%s", string(output))
	}

	if s.CallbacksHandler != nil {
		s.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}

// isDangerousCommand 检查命令是否在危险命令列表中
func (s *SystemCommand) isDangerousCommand(command string) bool {
	command = strings.ToLower(command)
	for _, dangerous := range s.DangerousCommands {
		if strings.ToLower(dangerous) == command {
			return true
		}
	}
	return false
}

// askUserPermission 询问用户是否允许执行危险命令
func (s *SystemCommand) askUserPermission(command string) bool {
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	green := color.New(color.FgGreen)

	fmt.Println()
	red.Printf("🚨 危险命令警告: '%s' 是潜在危险命令!\n", command)
	yellow.Println("执行此命令可能对系统造成不可逆损害。")
	fmt.Println()

	// 显示具体风险提示
	s.showCommandRisks(command)
	fmt.Println()

	for {
		green.Print("确定要执行这个危险命令吗? [yes/no]: ")

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return false
		}

		response := strings.ToLower(strings.TrimSpace(scanner.Text()))

		switch response {
		case "yes", "y", "是", "确定":
			yellow.Println("⚠️  用户确认执行危险命令")
			return true
		case "no", "n", "否", "取消":
			fmt.Println("✅ 危险命令已取消，系统安全得到保护")
			return false
		default:
			red.Println("❌ 请输入 'yes' 或 'no' (或 'y'/'n')")
			continue
		}
	}
}

// showCommandRisks 显示特定命令的风险提示
func (s *SystemCommand) showCommandRisks(command string) {
	command = strings.ToLower(command)

	fmt.Println("⚠️  具体风险:")
	switch command {
	case "rm", "del", "erase":
		fmt.Println("  • 可能永久删除重要文件和数据")
		fmt.Println("  • 删除操作通常无法撤销")
		fmt.Println("  • 建议先备份重要数据")
	case "shutdown", "reboot", "halt":
		fmt.Println("  • 将关闭或重启系统")
		fmt.Println("  • 可能导致正在运行的程序丢失数据")
		fmt.Println("  • 建议保存所有工作后再执行")
	case "chmod", "chown":
		fmt.Println("  • 将修改文件或目录权限")
		fmt.Println("  • 错误的权限设置可能导致系统无法正常运行")
		fmt.Println("  • 可能影响系统安全性")
	case "dd", "fdisk", "mkfs":
		fmt.Println("  • 可能覆盖或破坏磁盘数据")
		fmt.Println("  • 错误使用可能导致整个系统无法启动")
		fmt.Println("  • 强烈建议备份重要数据")
	case "kill", "killall", "taskkill":
		fmt.Println("  • 将强制终止进程")
		fmt.Println("  • 可能导致数据丢失或系统不稳定")
		fmt.Println("  • 建议先尝试优雅关闭进程")
	default:
		fmt.Println("  • 此命令可能对系统造成意外影响")
		fmt.Println("  • 请确保您了解此命令的具体作用")
		fmt.Println("  • 建议在非生产环境中先行测试")
	}
}

// AddDangerousCommand 添加危险命令
func (s *SystemCommand) AddDangerousCommand(command string) {
	s.DangerousCommands = append(s.DangerousCommands, command)
}

// SetDangerousCommands 设置危险命令列表
func (s *SystemCommand) SetDangerousCommands(commands []string) {
	s.DangerousCommands = commands
}

// GetDangerousCommands 获取危险命令列表
func (s *SystemCommand) GetDangerousCommands() []string {
	return s.DangerousCommands
}
