package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"

	"github.com/dean2027/aishell/pkg/app"
	"github.com/dean2027/aishell/pkg/ui"
)

// Runner CLI运行器
type Runner struct {
	chatBot         *app.ChatBot
	rl             *readline.Instance
	inputProcessor *InputProcessor
	config         *app.Config
	ctx            context.Context
}

// RunnerConfig 运行器配置
type RunnerConfig struct {
	HistoryFile string
	Prompt      string
}

// NewRunner 创建新的CLI运行器
func NewRunner(ctx context.Context, config *app.Config) (*Runner, error) {
	// 验证环境要求
	if err := app.ValidateRequirements(); err != nil {
		return nil, err
	}

	// 创建聊天机器人
	chatBot, err := app.NewChatBot(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("初始化聊天机器人失败: %w", err)
	}

	// 配置 readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          config.Prompt,
		HistoryFile:     config.HistoryFile,
		AutoComplete:    ui.CreateCompleter(ui.DefaultCompleterConfig()),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: FilterInput,
	})
	if err != nil {
		return nil, fmt.Errorf("初始化readline失败: %w", err)
	}

	inputProcessor := NewInputProcessor(rl)

	return &Runner{
		chatBot:         chatBot,
		rl:             rl,
		inputProcessor: inputProcessor,
		config:         config,
		ctx:            ctx,
	}, nil
}

// Run 运行CLI应用
func (r *Runner) Run() error {
	defer r.Close()

	// 打印欢迎信息
	ui.PrintWelcome()
	ui.PrintUsageTips()

	// 主循环
	return r.mainLoop()
}

// mainLoop 主循环逻辑
func (r *Runner) mainLoop() error {
	validation := DefaultInputValidation()

	for {
		// 读取输入
		input, err := r.inputProcessor.ReadInput()
		if err == readline.ErrInterrupt {
			if len(input) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		// 验证输入
		if err := validation.ValidateInput(input); err != nil {
			if input == "" {
				continue // 空输入时继续循环
			}
			ui.PrintError("输入验证失败", err)
			continue
		}

		// 处理特殊命令
		if r.handleSpecialCommands(input) {
			continue
		}

		// 处理用户输入
		if err := r.processUserInput(input); err != nil {
			ui.PrintError("处理输入失败", err)
			continue
		}
	}

	ui.PrintGoodbye()
	return nil
}

// handleSpecialCommands 处理特殊命令
func (r *Runner) handleSpecialCommands(input string) bool {
	switch {
	case IsExitCommand(input):
		ui.PrintGoodbye()
		os.Exit(0)
		return true
	case IsHelpCommand(input):
		ui.PrintHelp()
		return true
	case IsHistoryCommand(input):
		ui.PrintCommandHistory(r.rl)
		return true
	case IsClearCommand(input):
		r.clearScreen()
		return true
	}
	return false
}

// processUserInput 处理用户输入
func (r *Runner) processUserInput(input string) error {
	// 显示思考状态
	ui.PrintThinking()
	
	// 处理输入
	response, err := r.chatBot.ProcessInput(input)
	
	// 清除思考状态
	ui.ClearThinking()
	
	if err != nil {
		return err
	}

	// 显示回复
	ui.PrintResponse(response)
	return nil
}

// clearScreen 清屏
func (r *Runner) clearScreen() {
	// 使用 ANSI 转义序列清屏
	fmt.Print("\033[2J\033[H")
	
	// 重新打印欢迎信息（可选）
	ui.PrintWelcome()
}

// GetConfig 获取配置
func (r *Runner) GetConfig() *app.Config {
	return r.config
}

// Close 关闭运行器，清理资源
func (r *Runner) Close() error {
	if r.rl != nil {
		r.rl.Close()
	}
	if r.chatBot != nil {
		r.chatBot.Close()
	}
	return nil
}

// SetInterruptHandler 设置中断处理器
func (r *Runner) SetInterruptHandler(handler func()) {
	// 这里可以设置自定义的中断处理逻辑
	// 比如优雅地保存状态、清理资源等
}

// GetChatBot 获取聊天机器人实例（用于测试或扩展）
func (r *Runner) GetChatBot() *app.ChatBot {
	return r.chatBot
}
