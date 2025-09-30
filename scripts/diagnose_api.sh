#!/bin/bash

echo "🔍 AI Shell API 404错误诊断工具"
echo "================================="

# 检查环境变量
echo "📋 检查环境变量配置..."
if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ OPENAI_API_KEY 未设置"
    echo "请运行: export OPENAI_API_KEY=your_api_key"
    exit 1
else
    echo "✅ OPENAI_API_KEY 已设置"
fi

if [ -n "$OPENAI_BASE_URL" ]; then
    echo "✅ OPENAI_BASE_URL 已设置: $OPENAI_BASE_URL"
else
    echo "ℹ️  OPENAI_BASE_URL 未设置，使用默认端点"
fi

# 启用调试模式
export AISHELL_DEBUG="true"
echo "✅ 调试模式已启用"

# 编译测试程序
echo ""
echo "🔧 编译测试程序..."
cd /Volumes/sea-2/CODE/aishell
go build -o test_api scripts/test_api_connection.go

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译成功"

# 运行测试
echo ""
echo "🚀 运行API连接测试..."
./test_api

# 清理
rm -f test_api

echo ""
echo "📝 诊断完成！"
echo ""
echo "如果测试失败，请检查："
echo "1. OPENAI_API_KEY 是否正确"
echo "2. OPENAI_BASE_URL 是否正确（如果设置了的话）"
echo "3. 网络连接是否正常"
echo "4. 是否需要配置代理"
echo ""
echo "要启用详细日志，请运行："
echo "export AISHELL_DEBUG=true"
echo "go run cmd/aishell/main.go"