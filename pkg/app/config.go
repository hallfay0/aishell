package app

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
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		ConversationBufferSize: 100,
		MaxExecutorIterations:  30,
		HistoryFile:           "/tmp/aishell_history",
		Prompt:               "💻 智能终端> ",
		DebugMode:            false,
		HasSearchAPI:         false,
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
