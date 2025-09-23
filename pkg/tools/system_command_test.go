package tools

import (
	"context"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestSystemCommand_Name(t *testing.T) {
	cmd := NewSystemCommand()
	if cmd.Name() != "system_command" {
		t.Errorf("Expected name 'system_command', got '%s'", cmd.Name())
	}
}

func TestSystemCommand_Description(t *testing.T) {
	cmd := NewSystemCommand()
	desc := cmd.Description()
	if len(desc) == 0 {
		t.Error("Description should not be empty")
	}
	if !strings.Contains(desc, "系统命令") {
		t.Error("Description should contain '系统命令'")
	}
}

func TestSystemCommand_DangerousCommands(t *testing.T) {
	cmd := NewSystemCommand()
	
	// 测试默认危险命令
	dangerousCommands := cmd.GetDangerousCommands()
	if len(dangerousCommands) == 0 {
		t.Error("Should have default dangerous commands")
	}

	// 检查是否包含rm命令
	found := false
	for _, dangerous := range dangerousCommands {
		if dangerous == "rm" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should include 'rm' in dangerous commands")
	}
}

func TestSystemCommand_AddDangerousCommand(t *testing.T) {
	cmd := NewSystemCommand()
	initialCount := len(cmd.GetDangerousCommands())
	
	cmd.AddDangerousCommand("test_dangerous_command")
	newCount := len(cmd.GetDangerousCommands())
	
	if newCount != initialCount+1 {
		t.Error("Adding dangerous command should increase count by 1")
	}
}

func TestSystemCommand_Call_EmptyCommand(t *testing.T) {
	cmd := NewSystemCommand()
	ctx := context.Background()

	result, err := cmd.Call(ctx, "")
	if err != nil {
		t.Errorf("Call should not return error for empty command, got: %v", err)
	}

	if !strings.Contains(result, "错误：命令不能为空") {
		t.Error("Should return error message for empty command")
	}
}

func TestSystemCommand_Call_SafeCommand(t *testing.T) {
	cmd := NewSystemCommand()
	ctx := context.Background()
	
	// 测试安全命令（不在危险列表中）应该直接执行
	result, err := cmd.Call(ctx, "echo Hello World")
	if err != nil {
		t.Errorf("Call should not return error for safe command, got: %v", err)
	}
	
	if !strings.Contains(result, "Hello World") || !strings.Contains(result, "命令执行成功") {
		t.Errorf("Should execute safe command successfully, got: %s", result)
	}
}

func TestSystemCommand_Call_AllowedCommand(t *testing.T) {
	cmd := NewSystemCommand()
	// 移除了SetAllowedCommands，现在不需要设置允许列表
	ctx := context.Background()

	// 测试不同操作系统的echo命令
	var testCommand string
	if runtime.GOOS == "windows" {
		testCommand = "echo Hello World"
	} else {
		testCommand = "echo Hello World"
	}

	result, err := cmd.Call(ctx, testCommand)
	if err != nil {
		t.Errorf("Call should not return error for allowed command, got: %v", err)
	}

	if !strings.Contains(result, "Hello World") || !strings.Contains(result, "命令执行成功") {
		t.Errorf("Should execute echo command successfully, got: %s", result)
	}
}

func TestSystemCommand_Call_WithTimeout(t *testing.T) {
	cmd := NewSystemCommand()
	cmd.Timeout = 100 * time.Millisecond // 设置很短的超时时间
	// 移除了SetAllowedCommands，现在使用危险命令机制

	ctx := context.Background()

	// 测试超时场景
	var testCommand string
	if runtime.GOOS == "windows" {
		testCommand = "timeout 2" // Windows下等待2秒
	} else {
		testCommand = "sleep 2" // Unix下等待2秒
	}

	result, err := cmd.Call(ctx, testCommand)
	if err != nil {
		t.Errorf("Call should not return error even with timeout, got: %v", err)
	}

	// 应该包含超时或执行失败的信息
	if !strings.Contains(result, "命令执行失败") {
		t.Errorf("Should indicate command execution failure due to timeout, got: %s", result)
	}
}

func TestSystemCommand_isDangerousCommand(t *testing.T) {
	cmd := NewSystemCommand()
	
	tests := []struct {
		command  string
		expected bool
	}{
		{"rm", true},      // 危险命令
		{"del", true},     // 危险命令
		{"shutdown", true}, // 危险命令
		{"echo", false},   // 安全命令
		{"ls", false},     // 安全命令
		{"pwd", false},    // 安全命令
		{"RM", true},      // 大小写不敏感
		{"", false},       // 空命令
	}
	
	for _, test := range tests {
		result := cmd.isDangerousCommand(test.command)
		if result != test.expected {
			t.Errorf("isDangerousCommand('%s') = %v, expected %v", test.command, result, test.expected)
		}
	}
}

func TestSystemCommand_SetDangerousCommands(t *testing.T) {
	cmd := NewSystemCommand()
	newCommands := []string{"test1", "test2", "test3"}
	
	cmd.SetDangerousCommands(newCommands)
	dangerousCommands := cmd.GetDangerousCommands()
	
	if len(dangerousCommands) != 3 {
		t.Errorf("Expected 3 dangerous commands, got %d", len(dangerousCommands))
	}
	
	for i, expected := range newCommands {
		if dangerousCommands[i] != expected {
			t.Errorf("Expected command %s at index %d, got %s", expected, i, dangerousCommands[i])
		}
	}
}
