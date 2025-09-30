package app

import "fmt"

// Config 应用配置
type Config struct {
	// ConversationBufferSize 对话窗口缓冲大小，控制保持的对话轮数
	ConversationBufferSize int

	// MaxExecutorIterations 执行器最大迭代次数，控制单次交互的推理深度
	MaxExecutorIterations int

	// HistoryFile 历史文件路径
	HistoryFile string

	// Prompt 命令行提示符
	Prompt string

	// DebugMode 调试模式
	DebugMode bool

	// HasSearchAPI 是否有搜索API
	HasSearchAPI bool

	// OpenAIBaseURL OpenAI API基础URL，用于自定义端点
	OpenAIBaseURL string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		ConversationBufferSize: 100,
		MaxExecutorIterations:  30,
		HistoryFile:            "/tmp/aishell_history",
		Prompt:                 "💻 智能终端> ",
		DebugMode:              false,
		HasSearchAPI:           false,
		OpenAIBaseURL:          "", // 默认为空，使用OpenAI官方端点
	}
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	config := DefaultConfig()

	// 从环境变量读取配置
	if isDebugEnabled() {
		config.DebugMode = true
	}

	if hasSearchAPI() {
		config.HasSearchAPI = true
	}

	// 读取OpenAI BaseURL配置
	if baseURL := getEnv("OPENAI_BASE_URL"); baseURL != "" {
		config.OpenAIBaseURL = baseURL
	}

	// 添加调试日志
	if config.DebugMode {
		apiKey := getEnv("OPENAI_API_KEY")
		if apiKey != "" {
			fmt.Printf("🔍 [DEBUG] OpenAI API Key: %s... (前10个字符)\n", apiKey[:min(10, len(apiKey))])
		} else {
			fmt.Printf("🔍 [DEBUG] OpenAI API Key: (未设置)\n")
		}
		fmt.Printf("🔍 [DEBUG] OpenAI Base URL: %s\n", config.OpenAIBaseURL)
		fmt.Printf("🔍 [DEBUG] Has Search API: %v\n", config.HasSearchAPI)
		fmt.Printf("🔍 [DEBUG] Debug Mode: %v\n", config.DebugMode)
	}

	return config
}

// isDebugEnabled 检查是否启用调试模式
func isDebugEnabled() bool {
	return getEnv("AISHELL_DEBUG") == "true"
}

// hasSearchAPI 检查是否有搜索API配置
func hasSearchAPI() bool {
	return getEnv("SERPAPI_API_KEY") != ""
}

// hasOpenAIAPI 检查是否有OpenAI API配置
func HasOpenAIAPI() bool {
	return getEnv("OPENAI_API_KEY") != ""
}

// getOpenAIBaseURL 获取OpenAI BaseURL
func getOpenAIBaseURL() string {
	return getEnv("OPENAI_BASE_URL")
}

// getEnv 安全获取环境变量
func getEnv(key string) string {
	// 这里可以添加更复杂的环境变量处理逻辑
	// 比如从配置文件读取、设置默认值等
	return getOSEnv(key)
}

// getOSEnv 获取操作系统环境变量的抽象接口
var getOSEnv = func(key string) string {
	// 默认实现，可以在测试中替换
	return ""
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
