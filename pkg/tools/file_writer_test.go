package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestFileWriter_parseInput(t *testing.T) {
	fw := NewFileWriter()

	tests := []struct {
		name    string
		input   string
		want    *FileWriteParams
		wantErr bool
	}{
		{
			name:  "JSON格式 - 完整参数",
			input: `{"file_path": "test.txt", "content": "Hello World", "create_dirs": true}`,
			want: &FileWriteParams{
				FilePath:   "test.txt",
				Content:    "Hello World",
				CreateDirs: true,
			},
			wantErr: false,
		},
		{
			name:  "JSON格式 - 必需参数",
			input: `{"file_path": "test.txt", "content": "Hello World"}`,
			want: &FileWriteParams{
				FilePath:   "test.txt",
				Content:    "Hello World",
				CreateDirs: false,
			},
			wantErr: false,
		},
		{
			name:  "简单格式 - 完整参数",
			input: "test.txt|||Hello World|||true",
			want: &FileWriteParams{
				FilePath:   "test.txt",
				Content:    "Hello World",
				CreateDirs: true,
			},
			wantErr: false,
		},
		{
			name:  "简单格式 - 必需参数",
			input: "test.txt|||Hello World",
			want: &FileWriteParams{
				FilePath:   "test.txt",
				Content:    "Hello World",
				CreateDirs: false,
			},
			wantErr: false,
		},
		{
			name:  "JSON格式 - 多行内容",
			input: `{"file_path": "multi.txt", "content": "line1\nline2\nline3"}`,
			want: &FileWriteParams{
				FilePath: "multi.txt",
				Content:  "line1\nline2\nline3",
			},
			wantErr: false,
		},
		{
			name:    "空输入",
			input:   "",
			wantErr: true,
		},
		{
			name:    "无效JSON",
			input:   `{"file_path": "test.txt", "content":`,
			wantErr: true,
		},
		{
			name:    "简单格式缺少参数",
			input:   "test.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fw.parseInput(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseInput() 期望出现错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("parseInput() 出现意外错误 = %v", err)
				return
			}

			if got.FilePath != tt.want.FilePath {
				t.Errorf("parseInput() 文件路径 = %v, 期望 %v", got.FilePath, tt.want.FilePath)
			}
			if got.Content != tt.want.Content {
				t.Errorf("parseInput() 内容 = %v, 期望 %v", got.Content, tt.want.Content)
			}
			if got.CreateDirs != tt.want.CreateDirs {
				t.Errorf("parseInput() 创建目录 = %v, 期望 %v", got.CreateDirs, tt.want.CreateDirs)
			}
		})
	}
}

func TestFileWriter_validateParams(t *testing.T) {
	fw := NewFileWriter()

	tests := []struct {
		name    string
		params  *FileWriteParams
		wantErr bool
	}{
		{
			name: "有效参数 - txt文件",
			params: &FileWriteParams{
				FilePath: "test.txt",
				Content:  "Hello",
			},
			wantErr: false,
		},
		{
			name: "有效参数 - go文件",
			params: &FileWriteParams{
				FilePath: "main.go",
				Content:  "package main",
			},
			wantErr: false,
		},
		{
			name: "有效参数 - 无扩展名",
			params: &FileWriteParams{
				FilePath: "config",
				Content:  "debug=true",
			},
			wantErr: false,
		},
		{
			name: "空文件路径",
			params: &FileWriteParams{
				FilePath: "",
				Content:  "Hello",
			},
			wantErr: true,
		},
		{
			name: "路径遍历攻击",
			params: &FileWriteParams{
				FilePath: "../../../etc/passwd",
				Content:  "malicious",
			},
			wantErr: true,
		},
		{
			name: "不支持的文件类型",
			params: &FileWriteParams{
				FilePath: "image.jpg",
				Content:  "not text",
			},
			wantErr: true,
		},
		{
			name: "可执行文件",
			params: &FileWriteParams{
				FilePath: "app.exe",
				Content:  "binary",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fw.validateParams(tt.params)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateParams() 期望出现错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("validateParams() 出现意外错误 = %v", err)
			}
		})
	}
}

func TestFileWriter_Call(t *testing.T) {
	// 创建临时测试目录
	tmpDir := t.TempDir()
	
	fw := NewFileWriter()
	ctx := context.Background()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		setup   func() string // 返回测试文件路径
	}{
		{
			name: "写入新文件",
			setup: func() string {
				return filepath.Join(tmpDir, "new_file.txt")
			},
			wantErr: false,
		},
		{
			name: "覆盖现有文件",
			setup: func() string {
				testFile := filepath.Join(tmpDir, "existing_file.txt")
				os.WriteFile(testFile, []byte("old content"), 0644)
				return testFile
			},
			wantErr: false,
		},
		{
			name: "在不存在目录中创建文件 - 不创建目录",
			setup: func() string {
				return filepath.Join(tmpDir, "nonexistent", "file.txt")
			},
			wantErr: true,
		},
		{
			name: "在不存在目录中创建文件 - 创建目录",
			setup: func() string {
				return filepath.Join(tmpDir, "newdir", "file.txt")
			},
			wantErr: false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := tt.setup()
			
			var input string
			if i == len(tests)-1 { // 最后一个测试用例，启用创建目录
				input = fmt.Sprintf(`{"file_path": "%s", "content": "test content %d", "create_dirs": true}`, 
					testFile, i)
			} else {
				input = fmt.Sprintf(`{"file_path": "%s", "content": "test content %d"}`, 
					testFile, i)
			}
			
			result, err := fw.Call(ctx, input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Call() 期望出现错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("Call() 出现意外错误 = %v", err)
				return
			}

			if result == "" {
				t.Errorf("Call() 返回空结果")
			}

			// 验证文件是否被正确写入
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Errorf("无法读取测试文件: %v", err)
				return
			}

			expectedContent := fmt.Sprintf("test content %d", i)
			if string(content) != expectedContent {
				t.Errorf("文件内容 = %s, 期望 %s", string(content), expectedContent)
			}
		})
	}
}

func TestFileWriter_ToolInterface(t *testing.T) {
	fw := NewFileWriter()

	// 测试工具名称
	if fw.Name() != "file_writer" {
		t.Errorf("Name() = %v, 期望 file_writer", fw.Name())
	}

	// 测试工具描述
	desc := fw.Description()
	if desc == "" {
		t.Errorf("Description() 返回空字符串")
	}

	// 确保描述包含关键信息
	if !contains(desc, "file_path") {
		t.Errorf("Description() 不包含 'file_path' 参数说明")
	}
	if !contains(desc, "content") {
		t.Errorf("Description() 不包含 'content' 参数说明")
	}
	if !contains(desc, "create_dirs") {
		t.Errorf("Description() 不包含 'create_dirs' 参数说明")
	}
	if !contains(desc, "JSON") {
		t.Errorf("Description() 不包含 'JSON' 格式说明")
	}
}

func TestFileWriter_writeFile(t *testing.T) {
	tmpDir := t.TempDir()
	fw := NewFileWriter()

	tests := []struct {
		name    string
		params  *FileWriteParams
		setup   func() *FileWriteParams
		wantErr bool
	}{
		{
			name: "写入简单内容",
			setup: func() *FileWriteParams {
				return &FileWriteParams{
					FilePath: filepath.Join(tmpDir, "simple.txt"),
					Content:  "Hello, World!",
				}
			},
			wantErr: false,
		},
		{
			name: "写入多行内容",
			setup: func() *FileWriteParams {
				return &FileWriteParams{
					FilePath: filepath.Join(tmpDir, "multiline.txt"),
					Content:  "Line 1\nLine 2\nLine 3\n",
				}
			},
			wantErr: false,
		},
		{
			name: "写入特殊字符",
			setup: func() *FileWriteParams {
				return &FileWriteParams{
					FilePath: filepath.Join(tmpDir, "special.txt"),
					Content:  "特殊字符: @#$%^&*()_+ 中文测试 🚀",
				}
			},
			wantErr: false,
		},
		{
			name: "创建目录并写入文件",
			setup: func() *FileWriteParams {
				return &FileWriteParams{
					FilePath:   filepath.Join(tmpDir, "newdir", "subdir", "test.txt"),
					Content:    "Content in nested directory",
					CreateDirs: true,
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := tt.setup()
			
			bytesWritten, err := fw.writeFile(params)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("writeFile() 期望出现错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("writeFile() 出现意外错误 = %v", err)
				return
			}

			if bytesWritten != len(params.Content) {
				t.Errorf("writeFile() 写入字节数 = %v, 期望 %v", bytesWritten, len(params.Content))
			}

			// 验证文件内容
			content, err := os.ReadFile(params.FilePath)
			if err != nil {
				t.Errorf("无法读取测试文件: %v", err)
				return
			}

			if string(content) != params.Content {
				t.Errorf("文件内容不匹配，得到 %q, 期望 %q", string(content), params.Content)
			}
		})
	}
}
