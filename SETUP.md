# GitHub 仓库设置指南

## 1. 创建 GitHub 仓库

1. 访问 [GitHub](https://github.com) 并登录
2. 点击 "New repository"
3. 填写仓库信息：
   - Repository name: `NodeImageBackup`
   - Description: `NodeImage Backup Tool - 高性能图片同步工具`
   - 选择 Public 或 Private
   - 不要初始化 README（已有文件）

## 2. 推送代码到 GitHub

```bash
# 初始化 Git 仓库（如果还没有）
git init

# 添加远程仓库
git remote add origin https://github.com/your-username/NodeImageBackup.git

# 添加所有文件
git add .

# 提交代码
git commit -m "Initial commit: NodeImage Backup Tool"

# 推送到 GitHub
git push -u origin main
```

## 3. 配置 GitHub Actions

### 3.1 启用 Actions

1. 在仓库页面点击 "Actions" 标签
2. 点击 "Enable Actions"
3. 选择 "Go" 模板（可选，我们已有自定义工作流）

### 3.2 设置权限

1. 进入仓库 Settings → Actions → General
2. 在 "Workflow permissions" 部分：
   - 选择 "Read and write permissions"
   - 勾选 "Allow GitHub Actions to create and approve pull requests"
3. 点击 "Save"

## 4. 测试构建

### 4.1 测试 PR 构建

1. 创建新分支：
   ```bash
   git checkout -b test-build
   git push origin test-build
   ```

2. 在 GitHub 创建 Pull Request
3. 检查 Actions 是否自动运行

### 4.2 测试 Release 构建

```bash
# 创建版本标签
git tag v1.0.0

# 推送标签
git push origin v1.0.0
```

检查 Actions 是否：
1. 在所有平台构建成功
2. 自动创建 Release
3. 上传构建产物

## 5. 验证构建产物

在 [Releases](https://github.com/your-username/NodeImageBackup/releases) 页面应该看到：

- `nib-linux-amd64.tar.gz`
- `nib-linux-arm64.tar.gz`
- `nib-windows-amd64.zip`
- `nib-windows-arm64.zip`
- `nib-darwin-amd64.tar.gz`
- `nib-darwin-arm64.tar.gz`
- `checksums.txt`

## 6. 更新文档链接

记得更新以下文件中的链接：

- `README.md`: 更新仓库 URL
- `DEPLOYMENT.md`: 更新下载链接示例

## 7. 故障排除

### Actions 不运行
- 检查 `.github/workflows/` 目录是否存在
- 确认工作流文件语法正确
- 查看 Actions 权限设置

### 构建失败
- 检查 Go 版本兼容性
- 验证依赖项完整性
- 查看详细错误日志

### Release 未创建
- 确认标签格式为 `v*`
- 检查工作流条件设置
- 验证权限配置

## 8. 后续维护

### 发布新版本
```bash
git tag v1.1.0
git push origin v1.1.0
```

### 更新代码
```bash
git add .
git commit -m "feat: 新功能描述"
git push origin main
```

### 查看构建状态
- 访问 Actions 页面查看构建历史
- 检查 Release 页面查看发布历史 