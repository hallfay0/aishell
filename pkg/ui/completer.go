package ui

import (
	"github.com/chzyer/readline"
)

// CompleterConfig 自动补全配置
type CompleterConfig struct {
	EnableSystemCommands bool
	EnableFileOperations bool
	EnableCalculations   bool
	EnableSearching      bool
	EnableDiagnostics    bool
}

// DefaultCompleterConfig 返回默认的自动补全配置
func DefaultCompleterConfig() *CompleterConfig {
	return &CompleterConfig{
		EnableSystemCommands: true,
		EnableFileOperations: true,
		EnableCalculations:   true,
		EnableSearching:      true,
		EnableDiagnostics:    true,
	}
}

// CreateCompleter 创建自动补全器
func CreateCompleter(config *CompleterConfig) *readline.PrefixCompleter {
	if config == nil {
		config = DefaultCompleterConfig()
	}

	var items []readline.PrefixCompleterInterface

	// 系统命令补全
	if config.EnableSystemCommands {
		items = append(items, getSystemCommands()...)
	}

	// 文件操作补全
	if config.EnableFileOperations {
		items = append(items, getFileOperations()...)
	}

	// 计算分析补全
	if config.EnableCalculations {
		items = append(items, getCalculationCommands()...)
	}

	// 搜索功能补全
	if config.EnableSearching {
		items = append(items, getSearchCommands()...)
	}

	// 诊断功能补全
	if config.EnableDiagnostics {
		items = append(items, getDiagnosticCommands()...)
	}

	// 控制命令补全
	items = append(items, getControlCommands()...)

	return readline.NewPrefixCompleter(items...)
}

// getSystemCommands 获取系统管理命令补全
func getSystemCommands() []readline.PrefixCompleterInterface {
	return []readline.PrefixCompleterInterface{
		readline.PcItem("帮我安装Python"),
		readline.PcItem("帮我安装nodejs"),
		readline.PcItem("安装Docker"),
		readline.PcItem("查看系统配置"),
		readline.PcItem("检查磁盘空间"),
		readline.PcItem("查看内存使用"),
		readline.PcItem("创建项目目录"),
		readline.PcItem("查看当前文件"),
		readline.PcItem("查看运行的服务"),
		readline.PcItem("检查端口占用"),
		readline.PcItem("查看进程列表"),
		readline.PcItem("检查网络连接"),
	}
}

// getFileOperations 获取文件操作命令补全
func getFileOperations() []readline.PrefixCompleterInterface {
	return []readline.PrefixCompleterInterface{
		// 文件读取
		readline.PcItem("帮我读取main.go"),
		readline.PcItem("读取main.go的前10行"),
		readline.PcItem("查看config文件"),
		readline.PcItem("读取第20-30行"),
		readline.PcItem("查看package.json"),
		readline.PcItem("读取README.md"),
		
		// 文件写入
		readline.PcItem("创建一个config.txt文件"),
		readline.PcItem("写入Hello World到test.txt"),
		readline.PcItem("更新main.go文件"),
		readline.PcItem("创建新的代码文件"),
		readline.PcItem("创建Dockerfile"),
		readline.PcItem("生成README文件"),
	}
}

// getCalculationCommands 获取计算分析命令补全
func getCalculationCommands() []readline.PrefixCompleterInterface {
	return []readline.PrefixCompleterInterface{
		readline.PcItem("计算 (15 + 25) * 2"),
		readline.PcItem("计算"),
		readline.PcItem("分析数据"),
		readline.PcItem("转换单位"),
		readline.PcItem("1GB等于多少MB"),
		readline.PcItem("求解方程"),
		readline.PcItem("统计分析"),
	}
}

// getSearchCommands 获取搜索命令补全
func getSearchCommands() []readline.PrefixCompleterInterface {
	return []readline.PrefixCompleterInterface{
		readline.PcItem("搜索Go语言最佳实践"),
		readline.PcItem("查找解决方案"),
		readline.PcItem("最新技术动态"),
		readline.PcItem("搜索Docker教程"),
		readline.PcItem("查找Python库"),
		readline.PcItem("搜索前端框架"),
	}
}

// getDiagnosticCommands 获取诊断命令补全
func getDiagnosticCommands() []readline.PrefixCompleterInterface {
	return []readline.PrefixCompleterInterface{
		readline.PcItem("分析系统性能"),
		readline.PcItem("优化建议"),
		readline.PcItem("故障排查"),
		readline.PcItem("性能瓶颈分析"),
		readline.PcItem("内存泄漏检查"),
		readline.PcItem("服务启动失败"),
	}
}

// getControlCommands 获取控制命令补全
func getControlCommands() []readline.PrefixCompleterInterface {
	return []readline.PrefixCompleterInterface{
		readline.PcItem("help"),
		readline.PcItem("帮助"),
		readline.PcItem("history"),
		readline.PcItem("命令历史"),
		readline.PcItem("exit"),
		readline.PcItem("quit"),
		readline.PcItem("clear"),
		readline.PcItem("cls"),
	}
}
