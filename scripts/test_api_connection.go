package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	fmt.Println("🔍 API连接测试开始...")

	// 获取环境变量
	apiKey := os.Getenv("OPENAI_API_KEY")
	baseURL := os.Getenv("OPENAI_BASE_URL")

	fmt.Printf("🔍 API Key: %s... (前10个字符)\n", apiKey[:min(10, len(apiKey))])
	fmt.Printf("🔍 Base URL: %s\n", baseURL)

	// 测试HTTP连接
	testURL := "https://api.openai.com/v1/models"
	if baseURL != "" {
		testURL = baseURL + "/v1/models"
	}

	fmt.Printf("🔍 测试URL: %s\n", testURL)

	// 创建HTTP客户端测试连接
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		fmt.Printf("❌ 创建请求失败: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ HTTP请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("🔍 HTTP状态码: %d\n", resp.StatusCode)

	if resp.StatusCode == 200 {
		fmt.Println("✅ API连接测试成功")
	} else {
		fmt.Printf("❌ API连接测试失败，状态码: %d\n", resp.StatusCode)
	}

	// 测试OpenAI客户端初始化
	fmt.Println("\n🔍 测试OpenAI客户端初始化...")

	var llm *openai.LLM
	var initErr error

	if baseURL != "" {
		fmt.Println("🔍 使用自定义BaseURL初始化...")
		llm, initErr = openai.New(openai.WithBaseURL(baseURL))
	} else {
		fmt.Println("🔍 使用默认端点初始化...")
		llm, initErr = openai.New()
	}

	if initErr != nil {
		fmt.Printf("❌ OpenAI客户端初始化失败: %v\n", initErr)
		return
	}

	fmt.Println("✅ OpenAI客户端初始化成功")

	// 测试简单的API调用
	fmt.Println("\n🔍 测试简单API调用...")
	ctx := context.Background()

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, "Hello"),
	}

	_, err = llm.GenerateContent(ctx, messages)
	if err != nil {
		fmt.Printf("❌ API调用失败: %v\n", err)
		return
	}

	fmt.Println("✅ API调用测试成功")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
