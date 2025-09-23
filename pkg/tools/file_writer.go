package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
)

// FileWriteParams 写文件的参数结构
type FileWriteParams struct {
	FilePath   string `json:"file_path"`
	Content    string `json:"content"`
	CreateDirs bool   `json:"create_dirs,omitempty"`
}

// FileWriter 文件写入工具
type FileWriter struct {
	CallbacksHandler callbacks.Handler
}

// NewFileWriter 创建新的文件写入工具
func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

// Name 返回工具名称
func (f *FileWriter) Name() string {
	return "file_writer"
}

// Description 返回工具描述
func (f *FileWriter) Description() string {
	return `写入文件内容的工具。支持创建新文件或覆盖现有文件，可选择是否自动创建目录。
输入格式：JSON字符串
{
  "file_path": "文件路径（必需）",
  "content": "要写入的内容（必需）",
  "create_dirs": true/false（可选，默认false）
}

参数说明：
- file_path (必需): 要写入的文件路径，支持相对路径和绝对路径，仅支持文本文件
- content (必需): 要写入的内容，支持任意文本内容，自动处理特殊字符和编码
- create_dirs (可选): 是否自动创建不存在的目录，默认为false

示例：
{"file_path": "config.txt", "content": "debug=true\nport=8080"}
{"file_path": "/tmp/test.log", "content": "Application started", "create_dirs": true}
{"file_path": "src/main.go", "content": "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}", "create_dirs": true}`
}

// Call 执行文件写入
func (f *FileWriter) Call(ctx context.Context, input string) (string, error) {
	if f.CallbacksHandler != nil {
		f.CallbacksHandler.HandleToolStart(ctx, input)
	}

	// 解析输入参数
	params, err := f.parseInput(input)
	if err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	// 验证参数
	if err := f.validateParams(params); err != nil {
		return "", fmt.Errorf("参数验证失败: %w", err)
	}

	// 写入文件
	bytesWritten, err := f.writeFile(params)
	if err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	result := fmt.Sprintf("成功写入文件: %s\n写入内容: %d 字节\n路径: %s", 
		params.FilePath, bytesWritten, f.getAbsolutePath(params.FilePath))

	if f.CallbacksHandler != nil {
		f.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}

// parseInput 解析输入参数
func (f *FileWriter) parseInput(input string) (*FileWriteParams, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("输入不能为空")
	}

	var params FileWriteParams

	// 尝试解析JSON格式
	if strings.HasPrefix(input, "{") && strings.HasSuffix(input, "}") {
		if err := json.Unmarshal([]byte(input), &params); err != nil {
			return nil, fmt.Errorf("JSON解析失败: %w", err)
		}
		return &params, nil
	}

	// fallback: 尝试解析简单格式 "file_path|||content|||create_dirs"
	parts := strings.Split(input, "|||")
	if len(parts) < 2 {
		return nil, fmt.Errorf("参数格式错误，请使用JSON格式或 'file_path|||content|||create_dirs' 格式")
	}

	params.FilePath = strings.TrimSpace(parts[0])
	params.Content = parts[1] // 保留原始内容，包括空白字符

	// 解析可选的create_dirs参数
	if len(parts) > 2 {
		createDirsStr := strings.TrimSpace(strings.ToLower(parts[2]))
		params.CreateDirs = createDirsStr == "true" || createDirsStr == "1"
	}

	return &params, nil
}

// validateParams 验证参数
func (f *FileWriter) validateParams(params *FileWriteParams) error {
	if params.FilePath == "" {
		return fmt.Errorf("文件路径不能为空")
	}

	// 安全检查：防止路径遍历攻击
	cleanPath := filepath.Clean(params.FilePath)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("不允许使用相对路径符号 '..'")
	}

	// 检查文件扩展名，确保是文本文件
	ext := strings.ToLower(filepath.Ext(params.FilePath))
	textExtensions := []string{
		".txt", ".log", ".md", ".json", ".xml", ".yaml", ".yml", 
		".go", ".py", ".js", ".html", ".css", ".sql", ".sh", ".bat",
		".c", ".cpp", ".h", ".java", ".php", ".rb", ".rs", ".swift",
		".conf", ".config", ".ini", ".env", ".properties",
	}
	
	if ext != "" {
		isTextFile := false
		for _, validExt := range textExtensions {
			if ext == validExt {
				isTextFile = true
				break
			}
		}
		if !isTextFile {
			return fmt.Errorf("不支持的文件类型: %s，仅支持文本文件", ext)
		}
	}

	return nil
}

// writeFile 写入文件
func (f *FileWriter) writeFile(params *FileWriteParams) (int, error) {
	// 获取绝对路径
	absPath := f.getAbsolutePath(params.FilePath)
	
	// 获取目录路径
	dirPath := filepath.Dir(absPath)

	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if params.CreateDirs {
			// 创建目录
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return 0, fmt.Errorf("创建目录失败: %w", err)
			}
		} else {
			return 0, fmt.Errorf("目录不存在: %s，请设置 create_dirs=true 来自动创建", dirPath)
		}
	}

	// 写入文件
	err := os.WriteFile(absPath, []byte(params.Content), 0644)
	if err != nil {
		return 0, fmt.Errorf("写入文件失败: %w", err)
	}

	return len(params.Content), nil
}

// getAbsolutePath 获取绝对路径
func (f *FileWriter) getAbsolutePath(filePath string) string {
	// 如果已经是绝对路径，直接返回
	if filepath.IsAbs(filePath) {
		return filepath.Clean(filePath)
	}

	// 获取当前工作目录
	pwd, err := os.Getwd()
	if err != nil {
		return filePath // fallback
	}

	// 构建绝对路径
	absPath := filepath.Join(pwd, filePath)
	return filepath.Clean(absPath)
}

// checkFileExists 检查文件是否存在
func (f *FileWriter) checkFileExists(filePath string) bool {
	absPath := f.getAbsolutePath(filePath)
	_, err := os.Stat(absPath)
	return !os.IsNotExist(err)
}

// getFileInfo 获取文件信息
func (f *FileWriter) getFileInfo(filePath string) (bool, int64, error) {
	absPath := f.getAbsolutePath(filePath)
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	return true, info.Size(), nil
}
