package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileReader_parseInput(t *testing.T) {
	fr := NewFileReader()

	tests := []struct {
		name      string
		input     string
		wantPath  string
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{
			name:      "只有文件路径",
			input:     "main.go",
			wantPath:  "main.go",
			wantStart: 1,
			wantEnd:   100,
			wantErr:   false,
		},
		{
			name:      "文件路径和起始行",
			input:     "main.go,10",
			wantPath:  "main.go",
			wantStart: 10,
			wantEnd:   100,
			wantErr:   false,
		},
		{
			name:      "完整参数",
			input:     "main.go,5,15",
			wantPath:  "main.go",
			wantStart: 5,
			wantEnd:   15,
			wantErr:   false,
		},
		{
			name:      "绝对路径",
			input:     "/tmp/test.txt,1,20",
			wantPath:  "/tmp/test.txt",
			wantStart: 1,
			wantEnd:   20,
			wantErr:   false,
		},
		{
			name:    "空输入",
			input:   "",
			wantErr: true,
		},
		{
			name:    "结束行小于起始行",
			input:   "test.txt,10,5",
			wantErr: true,
		},
		{
			name:    "起始行号为0",
			input:   "test.txt,0,10",
			wantErr: true,
		},
		{
			name:    "无效的起始行号",
			input:   "test.txt,abc,10",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, start, end, err := fr.parseInput(tt.input)

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

			if path != tt.wantPath {
				t.Errorf("parseInput() 文件路径 = %v, 期望 %v", path, tt.wantPath)
			}
			if start != tt.wantStart {
				t.Errorf("parseInput() 起始行 = %v, 期望 %v", start, tt.wantStart)
			}
			if end != tt.wantEnd {
				t.Errorf("parseInput() 结束行 = %v, 期望 %v", end, tt.wantEnd)
			}
		})
	}
}

func TestFileReader_Call(t *testing.T) {
	// 创建临时测试文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// 写入测试内容
	testContent := `第一行内容
第二行内容
第三行内容
第四行内容
第五行内容`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("无法创建测试文件: %v", err)
	}

	fr := NewFileReader()
	ctx := context.Background()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "读取存在的文件",
			input:   testFile + ",1,3",
			wantErr: false,
		},
		{
			name:    "读取不存在的文件",
			input:   "/nonexistent/file.txt",
			wantErr: true,
		},
		{
			name:    "空输入",
			input:   "",
			wantErr: true,
		},
		{
			name:    "超出文件行数",
			input:   testFile + ",10,20",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fr.Call(ctx, tt.input)

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
		})
	}
}

func TestFileReader_ToolInterface(t *testing.T) {
	fr := NewFileReader()

	// 测试工具名称
	if fr.Name() != "file_reader" {
		t.Errorf("Name() = %v, 期望 file_reader", fr.Name())
	}

	// 测试工具描述
	desc := fr.Description()
	if desc == "" {
		t.Errorf("Description() 返回空字符串")
	}

	// 确保描述包含关键信息
	if !contains(desc, "file_path") {
		t.Errorf("Description() 不包含 'file_path' 参数说明")
	}
	if !contains(desc, "start_line") {
		t.Errorf("Description() 不包含 'start_line' 参数说明")
	}
	if !contains(desc, "end_line") {
		t.Errorf("Description() 不包含 'end_line' 参数说明")
	}
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
