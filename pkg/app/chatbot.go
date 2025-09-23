package app

import (
	"context"
	"fmt"
	"os"

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

// ChatBot AI聊天机器人
type ChatBot struct {
	executor *agents.Executor
	llm      llms.Model
	ctx      context.Context
	config   *Config
}

// NewChatBot 创建新的聊天机器人实例
func NewChatBot(ctx context.Context, config *Config) (*ChatBot, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 初始化OpenAI LLM
	llm, err := openai.New()
	if err != nil {
		return nil, fmt.Errorf("初始化LLM失败: %w", err)
	}

	// 初始化对话窗口缓冲内存 (保持最近N轮对话)
	conversationMemory := memory.NewConversationWindowBuffer(config.ConversationBufferSize)

	// 创建工具列表
	toolsList := createToolsList(config)

	// 创建智能终端助手的专用系统提示
	systemPromptPrefix := prompt.CreateSystemPrompt()

	// 创建使用内存和自定义系统提示的对话代理
	agent := agents.NewConversationalAgent(llm, toolsList,
		agents.WithMemory(conversationMemory),
		agents.WithPromptPrefix(systemPromptPrefix),
	)

	// 创建执行器选项
	executorOptions := createExecutorOptions(config, conversationMemory)
	executor := agents.NewExecutor(agent, executorOptions...)

	return &ChatBot{
		executor: executor,
		llm:      llm,
		ctx:      ctx,
		config:   config,
	}, nil
}

// ProcessInput 处理用户输入
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

// GetConfig 获取配置
func (cb *ChatBot) GetConfig() *Config {
	return cb.config
}

// Close 关闭聊天机器人，清理资源
func (cb *ChatBot) Close() error {
	// 这里可以添加清理逻辑，比如保存对话历史等
	return nil
}

// createToolsList 创建工具列表
func createToolsList(config *Config) []tools.Tool {
	toolsList := []tools.Tool{
		tools.Calculator{},
		localtools.NewSystemCommand(),
		localtools.NewFileReader(),
		localtools.NewFileWriter(),
	}

	// 如果设置了SERPAPI_API_KEY，添加搜索工具
	if config.HasSearchAPI {
		searchTool, err := serpapi.New()
		if err == nil {
			toolsList = append(toolsList, searchTool)
		}
	}

	return toolsList
}

// createExecutorOptions 创建执行器选项
func createExecutorOptions(config *Config, conversationMemory *memory.ConversationWindowBuffer) []agents.Option {
	var executorOptions []agents.Option
	
	executorOptions = append(executorOptions, agents.WithMaxIterations(config.MaxExecutorIterations))
	executorOptions = append(executorOptions, agents.WithMemory(conversationMemory))

	if config.DebugMode {
		fmt.Println("🔍 调试模式已启用 - 将显示详细的执行日志")
		debugHandler := callbacks.LogHandler{}
		executorOptions = append(executorOptions, agents.WithCallbacksHandler(debugHandler))
	}

	return executorOptions
}

// ValidateRequirements 验证运行环境要求
func ValidateRequirements() error {
	if !HasOpenAIAPI() {
		return fmt.Errorf("未设置OPENAI_API_KEY环境变量，请设置: export OPENAI_API_KEY=your_api_key")
	}
	return nil
}

// init 初始化包
func init() {
	// 设置默认的环境变量获取函数
	getOSEnv = os.Getenv
}
