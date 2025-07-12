#!/bin/bash

# NodeImage Backup Tool 测试脚本

set -e

echo "🧪 开始测试 NodeImage Backup Tool..."

# 检查程序是否存在
if [ ! -f "./nib" ]; then
    echo "❌ 错误: 程序文件 ./nib 不存在，请先构建程序"
    exit 1
fi

echo "✅ 程序文件存在"

# 测试帮助命令
echo ""
echo "📋 测试帮助命令..."
./nib --help > /dev/null && echo "✅ 主帮助命令正常"
./nib sync --help > /dev/null && echo "✅ sync帮助命令正常"
./nib list --help > /dev/null && echo "✅ list帮助命令正常"

# 测试参数验证
echo ""
echo "🔍 测试参数验证..."
# 临时移除配置文件进行测试
if [ -f "nib.yaml" ]; then
    mv nib.yaml nib.yaml.test
fi
if ./nib 2>&1 | grep -q "请通过配置文件或 -t 参数提供API Token"; then
    echo "✅ token参数验证正常"
else
    echo "❌ token参数验证失败"
fi
# 恢复配置文件
if [ -f "nib.yaml.test" ]; then
    mv nib.yaml.test nib.yaml
fi

# 测试无效token
echo ""
echo "🔑 测试无效token..."
if ./nib -t "invalid_token" 2>&1 | grep -q "API"; then
    echo "✅ 无效token处理正常"
else
    echo "❌ 无效token处理异常"
fi

echo ""
echo "🎉 基础测试完成!"
echo ""
echo "💡 提示: 使用真实token进行完整测试:"
echo "   ./nib list -t YOUR_REAL_TOKEN"
echo "   ./nib -t YOUR_REAL_TOKEN" 