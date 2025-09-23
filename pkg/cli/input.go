package cli

import (
	"strings"

	"github.com/chzyer/readline"
)

// InputProcessor 输入处理器
type InputProcessor struct {
	rl *readline.Instance
}

// NewInputProcessor 创建新的输入处理器
func NewInputProcessor(rl *readline.Instance) *InputProcessor {
	return &InputProcessor{
		rl: rl,
	}
}

// ReadInput 读取用户输入
func (ip *InputProcessor) ReadInput() (string, error) {
	line, err := ip.rl.Readline()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(line), nil
}

// IsExitCommand 检查是否为退出命令
func IsExitCommand(input string) bool {
	lower := strings.ToLower(input)
	return lower == "exit" || lower == "quit"
}

// IsHelpCommand 检查是否为帮助命令
func IsHelpCommand(input string) bool {
	lower := strings.ToLower(input)
	return lower == "help" || lower == "帮助"
}

// IsHistoryCommand 检查是否为历史命令
func IsHistoryCommand(input string) bool {
	lower := strings.ToLower(input)
	return lower == "history" || lower == "命令历史"
}

// IsClearCommand 检查是否为清屏命令
func IsClearCommand(input string) bool {
	lower := strings.ToLower(input)
	return lower == "clear" || lower == "cls"
}

// FilterInput 过滤输入字符
func FilterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

// InputValidation 输入验证
type InputValidation struct {
	MinLength int
	MaxLength int
	AllowEmpty bool
}

// DefaultInputValidation 默认输入验证配置
func DefaultInputValidation() *InputValidation {
	return &InputValidation{
		MinLength:  0,
		MaxLength:  1000,
		AllowEmpty: false,
	}
}

// ValidateInput 验证输入
func (iv *InputValidation) ValidateInput(input string) error {
	if !iv.AllowEmpty && len(input) == 0 {
		return NewInputError("输入不能为空")
	}
	
	if len(input) < iv.MinLength {
		return NewInputError("输入长度不足")
	}
	
	if len(input) > iv.MaxLength {
		return NewInputError("输入长度超出限制")
	}
	
	return nil
}

// InputError 输入错误
type InputError struct {
	message string
}

// NewInputError 创建输入错误
func NewInputError(message string) *InputError {
	return &InputError{message: message}
}

// Error 实现error接口
func (e *InputError) Error() string {
	return e.message
}
