#!/bin/bash

# 获取当前脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 定义构建参数
BINARY_NAME="minichat"
PLATFORMS="windows/amd64 windows/arm64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"

# 清理之前的构建
rm -rf ./build
mkdir -p ./build

# 构建每个平台
for PLATFORM in $PLATFORMS; do
    OS=$(echo $PLATFORM | cut -d'/' -f1)
    ARCH=$(echo $PLATFORM | cut -d'/' -f2)

    echo "Building for $OS/$ARCH..."

    # 设置输出路径
    if [ "$OS" = "windows" ]; then
        EXT=".exe"
    else
        EXT=""
    fi

    # 执行构建
    GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w" -o "./build/$BINARY_NAME-$OS-$ARCH$EXT"

    # 检查构建是否成功
    if [ $? -eq 0 ]; then
        echo "$OS/$ARCH build completed successfully"
    else
        echo "$OS/$ARCH build failed"
        exit 1
    fi
done

echo "All builds completed!"
echo "Binary files are located in the ./build directory"

# 回到原始目录
cd - > /dev/null