package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
)

// FileReader 文件读取工具
type FileReader struct {
	CallbacksHandler callbacks.Handler
}

// NewFileReader 创建新的文件读取工具
func NewFileReader() *FileReader {
	return &FileReader{}
}

// Name 返回工具名称
func (f *FileReader) Name() string {
	return "file_reader"
}

// Description 返回工具描述
func (f *FileReader) Description() string {
	return `读取文件内容的工具。可以按行号范围读取文件，支持相对路径和绝对路径。
输入格式：file_path[,start_line,end_line]
参数说明：
- file_path (必需): 要读取的文件路径，支持相对路径和绝对路径
- start_line (可选): 起始行号，从1开始计数，默认为1
- end_line (可选): 结束行号，必须大于等于start_line，默认为100

示例：
- "main.go" - 读取main.go文件的前100行
- "main.go,1,50" - 读取main.go文件的第1-50行
- "/path/to/file.txt,10,20" - 读取文件的第10-20行`
}

// Call 执行文件读取
func (f *FileReader) Call(ctx context.Context, input string) (string, error) {
	if f.CallbacksHandler != nil {
		f.CallbacksHandler.HandleToolStart(ctx, input)
	}

	// 解析输入参数
	filePath, startLine, endLine, err := f.parseInput(input)
	if err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	// 读取文件内容
	content, err := f.readFileLines(filePath, startLine, endLine)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	result := fmt.Sprintf("文件: %s (第%d-%d行)\n%s", filePath, startLine, endLine, content)

	if f.CallbacksHandler != nil {
		f.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}

// parseInput 解析输入参数
func (f *FileReader) parseInput(input string) (filePath string, startLine, endLine int, err error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", 0, 0, fmt.Errorf("文件路径不能为空")
	}

	parts := strings.Split(input, ",")

	// 解析文件路径
	filePath = strings.TrimSpace(parts[0])
	if filePath == "" {
		return "", 0, 0, fmt.Errorf("文件路径不能为空")
	}

	// 默认值
	startLine = 1
	endLine = 100

	// 解析起始行号
	if len(parts) > 1 {
		startStr := strings.TrimSpace(parts[1])
		if startStr != "" {
			startLine, err = strconv.Atoi(startStr)
			if err != nil {
				return "", 0, 0, fmt.Errorf("起始行号格式错误: %s", startStr)
			}
			if startLine < 1 {
				return "", 0, 0, fmt.Errorf("起始行号必须大于0")
			}
		}
	}

	// 解析结束行号
	if len(parts) > 2 {
		endStr := strings.TrimSpace(parts[2])
		if endStr != "" {
			endLine, err = strconv.Atoi(endStr)
			if err != nil {
				return "", 0, 0, fmt.Errorf("结束行号格式错误: %s", endStr)
			}
		}
	}

	// 验证行号范围
	if endLine < startLine {
		return "", 0, 0, fmt.Errorf("结束行号(%d)不能小于起始行号(%d)", endLine, startLine)
	}

	return filePath, startLine, endLine, nil
}

// readFileLines 读取文件指定行号范围的内容
func (f *FileReader) readFileLines(filePath string, startLine, endLine int) (string, error) {
	// 处理相对路径
	absPath, err := f.getAbsolutePath(filePath)
	if err != nil {
		return "", err
	}

	// 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("文件不存在: %s", absPath)
	}

	// 打开文件
	file, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// 逐行读取
	scanner := bufio.NewScanner(file)
	var lines []string
	currentLine := 1

	for scanner.Scan() {
		if currentLine >= startLine && currentLine <= endLine {
			// 格式化输出：行号|内容
			lines = append(lines, fmt.Sprintf("%6d|%s", currentLine, scanner.Text()))
		}

		// 如果超过结束行，提前退出
		if currentLine > endLine {
			break
		}

		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取文件时出错: %w", err)
	}

	// 检查是否读取到内容
	if len(lines) == 0 {
		if currentLine-1 < startLine {
			return "", fmt.Errorf("文件只有%d行，起始行号%d超出范围", currentLine-1, startLine)
		}
		return "指定范围内没有内容", nil
	}

	return strings.Join(lines, "\n"), nil
}

// getAbsolutePath 获取绝对路径
func (f *FileReader) getAbsolutePath(filePath string) (string, error) {
	// 如果已经是绝对路径，直接返回
	if filepath.IsAbs(filePath) {
		return filePath, nil
	}

	// 获取当前工作目录
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("无法获取当前工作目录: %w", err)
	}

	// 构建绝对路径
	absPath := filepath.Join(pwd, filePath)
	return absPath, nil
}
