# NodeImage Backup Tool (nib)

一个高性能的 NodeImage 图片同步工具，支持从远程 API 单向同步图片到本地目录。

## ✨ 特性

- 🚀 **高性能**: 支持并发下载，默认10个并发线程
- 📦 **小体积**: 使用 Go 语言编译，生成的可执行文件体积小
- 🔄 **单向同步**: 从远程到本地的单向同步
- 🖥️ **跨平台**: 支持 macOS、Windows、Linux 多平台
- 🛡️ **安全**: 使用临时文件确保下载完整性，用户确认机制
- 📊 **详细日志**: 提供详细的同步进度和统计信息
- ⚙️ **配置灵活**: 支持配置文件，命令行参数优先

## 📋 功能

1. **远程图片列表**: 获取 NodeImage API 中的所有图片
2. **本地图片扫描**: 扫描本地目录中的图片文件
3. **智能同步**: 
   - 远程有，本地没有 → 下载到本地
   - 远程有，本地也有 → 保持不变
   - 远程没有，本地有 → 删除本地文件
4. **用户确认**: 执行删除和下载前需要用户确认
5. **自动创建目录**: 本地目录不存在时自动创建

## 🚀 快速开始

### 构建

#### 方式一：GitHub Actions (推荐)

项目已配置 GitHub Actions 自动构建，支持所有平台：

```bash
# 推送版本标签触发构建
git tag v1.0.0
git push origin v1.0.0
```

构建完成后，在 [Releases](https://github.com/your-repo/NodeImageBackup/releases) 页面下载对应平台的可执行文件。

#### 方式二：本地构建

```bash
# 克隆项目
git clone <repository-url>
cd NodeImageBackup

# 构建所有平台版本
./build.sh
```

详细说明请参考 [DEPLOYMENT.md](DEPLOYMENT.md)。

### 使用

```bash
# 基本同步命令（推荐）
./nib

# 指定本地目录
./nib -d /path/to/local/directory

# 调整并发数量
./nib -w 20

# 查看远程图片列表
./nib list

# 显示调试信息
./nib --debug
```

## 📖 命令说明

### 默认同步命令

直接运行 `./nib` 执行同步操作。

**参数:**
- `-t, --token`: API Token (可通过配置文件指定)
- `-d, --dir`: 本地同步目录 (可通过配置文件指定，默认: 程序目录/images)
- `-w, --workers`: 并发下载数量 (可通过配置文件指定，默认: 10)
- `-c, --config`: 配置文件路径 (默认: nib.yaml 或 nib.yml)
- `--debug`: 显示调试信息

**示例:**
```bash
# 使用配置文件（推荐）
./nib

# 指定token
./nib -t YOUR_API_TOKEN

# 指定目录
./nib -d /Users/username/Pictures/NodeImage

# 使用20个并发线程
./nib -w 20

# 显示调试信息
./nib --debug
```

### `sync` 命令

显式执行同步操作（与默认命令效果相同）。

**示例:**
```bash
./nib sync -t YOUR_API_TOKEN
```

### `list` 命令

列出远程图片信息。

**参数:**
- `-t, --token`: API Token (可通过配置文件指定)
- `-c, --config`: 配置文件路径
- `--debug`: 显示调试信息

**示例:**
```bash
./nib list
./nib list -t YOUR_API_TOKEN
```

## ⚙️ 配置文件

创建 `nib.yaml` 文件来配置默认参数：

```yaml
token: YOUR_API_TOKEN
dir: ./images
workers: 10
```

**配置文件优先级:**
1. 命令行参数（最高）
2. 配置文件
3. 默认值（最低）

### 快速配置

```bash
# 复制配置文件模板
cp nib.yaml.example nib.yaml

# 编辑配置文件，填入你的API Token
nano nib.yaml
```

**注意:** 真实的 `nib.yaml` 文件包含敏感信息，不会被提交到GitHub。请使用 `nib.yaml.example` 作为模板。

## 🔧 支持的图片格式

- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- BMP (.bmp)
- WebP (.webp)
- AVIF (.avif)
- SVG (.svg)

## 📦 构建产物

构建脚本会生成以下文件：

### macOS
- `nib-darwin-amd64` - Intel Mac
- `nib-darwin-arm64` - Apple Silicon Mac

### Linux
- `nib-linux-amd64` - Intel/AMD 64位
- `nib-linux-arm64` - ARM 64位

### Windows
- `nib-windows-amd64.exe` - Intel/AMD 64位
- `nib-windows-arm64.exe` - ARM 64位

## 🔍 工作原理

1. **获取远程列表**: 调用 NodeImage API 获取所有图片信息
2. **扫描本地文件**: 递归扫描本地目录中的图片文件
3. **对比分析**: 比较远程和本地文件列表
4. **用户确认**: 执行操作前询问用户确认
5. **执行同步**:
   - 删除本地多余文件
   - 并发下载缺失文件
   - 使用临时文件确保下载完整性

## ⚡ 性能优化

- **并发下载**: 默认10个并发线程，可自定义
- **临时文件**: 使用临时文件确保下载完整性
- **原子操作**: 下载完成后原子性重命名
- **内存优化**: 流式处理，避免大文件内存占用

## 🛠️ 开发

### 环境要求

- Go 1.21+
- Git

### 本地开发

```bash
# 安装依赖
go mod tidy

# 本地运行
go run main.go

# 本地构建
go build -o nib main.go
```

### 测试

```bash
# 运行测试
make test

# 构建测试
make build
```

## 📝 日志输出示例

```
🔍 正在获取远程图片列表...
📁 正在扫描本地图片...
📊 统计信息:
   远程图片: 15 张
   本地图片: 8 张

🔄 同步计划:
   需要下载: 7 张
   需要删除: 0 张

⬇️  正在下载图片...
确认下载 7 个远程文件? (Y/n): y
   ✅ 已下载: image1.jpg (25376 bytes)
   ✅ 已下载: image2.png (15678 bytes)
   ✅ 已下载: image3.webp (8923 bytes)
   ...

🎉 同步完成!
```

## 🔒 安全说明

- API Token 仅在内存中使用，不会保存到文件
- 使用 HTTPS 进行所有网络通信
- 下载文件使用临时文件，确保完整性
- 支持超时控制，避免长时间等待
- 用户确认机制，防止误操作
- 不会产生垃圾文件，临时文件自动清理

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 支持

如有问题，请提交 Issue 或联系开发者。 