# 部署说明

## GitHub Actions 自动构建

本项目使用 GitHub Actions 实现自动跨平台构建，支持以下平台：

- **Linux**: AMD64, ARM64
- **Windows**: AMD64, ARM64  
- **macOS**: Intel (AMD64), Apple Silicon (ARM64)

### 触发构建

#### 1. 发布版本 (自动创建 Release)

推送版本标签即可触发构建并创建 Release：

```bash
# 创建并推送版本标签
git tag v1.0.0
git push origin v1.0.0
```

这将：
- 在所有平台构建可执行文件
- 创建压缩包 (.tar.gz, .zip)
- 生成校验和文件
- 自动创建 GitHub Release

#### 2. PR 测试

每次 Pull Request 都会触发测试：
- 构建验证
- 运行测试脚本
- 检查代码质量

### 构建产物

Release 包含以下文件：

```
nib-linux-amd64.tar.gz      # Linux Intel/AMD
nib-linux-arm64.tar.gz      # Linux ARM64
nib-windows-amd64.zip       # Windows Intel/AMD
nib-windows-arm64.zip       # Windows ARM64
nib-darwin-amd64.tar.gz     # macOS Intel
nib-darwin-arm64.tar.gz     # macOS Apple Silicon
checksums.txt               # 文件校验和
```

### 本地构建 (可选)

如果需要本地构建：

```bash
# 构建所有平台
./build.sh

# 或使用 Makefile
make cross-build
```

### 工作流文件

- `.github/workflows/build.yml` - 跨平台构建和发布
- `.github/workflows/test.yml` - PR 测试验证

### 配置要求

确保仓库设置：
1. **Actions 权限**: Settings → Actions → General → Workflow permissions → "Read and write permissions"
2. **Release 权限**: 自动使用 `GITHUB_TOKEN`

### 使用示例

用户下载和使用：

```bash
# 1. 下载对应平台的版本
# 例如 macOS Apple Silicon:
curl -L -o nib.tar.gz https://github.com/your-repo/NodeImageBackup/releases/latest/download/nib-darwin-arm64.tar.gz

# 2. 解压
tar -xzf nib.tar.gz

# 3. 设置权限
chmod +x nib

# 4. 创建配置文件
cp nib.yaml.example nib.yaml
nano nib.yaml

# 5. 使用
./nib
```

### 故障排除

#### 构建失败
- 检查 Go 版本兼容性
- 验证依赖项完整性
- 查看 Actions 日志

#### Release 未创建
- 确认标签格式为 `v*` (如 v1.0.0)
- 检查仓库权限设置
- 验证工作流文件语法

#### 文件缺失
- 检查构建矩阵配置
- 验证 artifact 上传
- 查看 release 步骤日志

### 自定义构建

如需修改构建配置：

1. **添加平台**: 在 `matrix.include` 中添加新条目
2. **修改参数**: 调整 `go build` 参数
3. **添加步骤**: 在 steps 中添加自定义操作

### 性能优化

- **缓存**: 使用 Go modules 缓存加速构建
- **并行**: 矩阵构建并行执行
- **精简**: 使用 `-ldflags "-s -w"` 减小文件大小 