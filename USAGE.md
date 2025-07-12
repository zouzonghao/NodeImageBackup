# NodeImage Backup Tool 使用示例

## 🚀 快速开始

### 1. 首次使用

```bash
# 首次运行，自动生成配置文件
./nib -t YOUR_API_TOKEN
```

### 2. 查看远程图片列表

```bash
# 使用配置文件（推荐）
./nib list

# 使用命令行参数
./nib list -t YOUR_API_TOKEN

# 使用交互式命令
make list
```

### 3. 同步图片到本地

```bash
# 使用配置文件（推荐）
./nib

# 强制同步，无需确认
./nib -y

# 同步到指定目录
./nib -d /path/to/images

# 使用20个并发线程
./nib -w 20

# 显示调试信息
./nib --debug

# 使用交互式命令
make sync
```

## 📋 实际使用示例

### 示例1: 基本同步

```bash
# 使用配置文件（推荐方式）
./nib

# 输出示例:
# 🔍 正在获取远程图片列表...
# 📁 正在扫描本地图片...
# 📊 统计信息:
#    远程图片: 5 张
#    本地图片: 2 张
# 
# 🔄 同步计划:
#    需要下载: 3 张
#    需要删除: 0 张
# 
# ⬇️  正在下载图片...
# 确认下载 3 个远程文件? (Y/n): y
#    ✅ 已下载: Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J.avif (0.02 MB)
#    ✅ 已下载: image2.png (0.02 MB)
#    ✅ 已下载: photo3.jpg (0.01 MB)
# 
# 🎉 同步完成!
```

### 示例2: 指定目录同步

```bash
# 同步到用户图片目录
./nib -d ~/Pictures/NodeImage

# 同步到桌面
./nib -d ~/Desktop/backup_images
```

### 示例3: 高性能同步

```bash
# 使用更多并发线程加速下载
./nib -w 20

# 适合网络较好的环境
./nib -w 50
```

### 示例4: 查看远程图片

```bash
# 查看所有远程图片
./nib list

# 输出示例:
# 🔍 正在获取远程图片列表...
# 
# 📋 远程图片列表 (5 张):
#   1. Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J.avif (0.02 MB)
#   2. screenshot.png (0.02 MB)
#   3. photo.jpg (0.01 MB)
#   4. logo.svg (0.00 MB)
#   5. banner.webp (0.00 MB)
```

### 示例5: 调试模式

```bash
# 显示API响应详情
./nib --debug

# 查看远程图片列表并显示调试信息
./nib list --debug
```

## 🔧 高级用法

### 1. 配置文件使用

#### 自动生成配置（推荐）

```bash
# 首次运行，自动生成配置文件
./nib -t YOUR_API_TOKEN
```

#### 手动创建配置

```bash
# 复制配置文件模板
cp nib.yaml.example nib.yaml

# 编辑配置文件
nano nib.yaml
```

配置文件内容示例：
```yaml
# NodeImage Backup Tool 配置文件
token: YOUR_API_TOKEN
dir: ./images
workers: 15
```

然后直接使用：
```bash
./nib
```

#### 更新配置

```bash
# 更新token，配置文件会自动更新
./nib -t NEW_API_TOKEN
```

**注意:** 真实的 `nib.yaml` 文件包含敏感信息，不会被提交到GitHub。请使用 `nib.yaml.example` 作为模板。

### 2. 定时同步

```bash
# 使用cron定时同步 (每天凌晨2点)
# 编辑crontab: crontab -e
# 添加以下行:
0 2 * * * /path/to/nib -y >> /var/log/nib.log 2>&1
```

### 3. 脚本化使用

```bash
#!/bin/bash
# sync_images.sh

BACKUP_DIR="/path/to/backup"

echo "开始同步图片..."
./nib -d $BACKUP_DIR -y

if [ $? -eq 0 ]; then
    echo "同步成功!"
    # 可以添加通知逻辑
else
    echo "同步失败!"
    exit 1
fi
```

### 4. 环境变量使用

```bash
# 设置环境变量
export NIB_DIR="/path/to/images"

# 在脚本中使用
./nib -d $NIB_DIR
```

## 🛠️ 故障排除

### 常见问题

1. **API Token 无效**
   ```
   错误: API返回错误: Invalid API key
   ```
   解决: 检查token是否正确，确保没有多余的空格

2. **网络连接问题**
   ```
   错误: 请求失败: dial tcp: lookup api.nodeimage.com: no such host
   ```
   解决: 检查网络连接，确认DNS设置

3. **权限问题**
   ```
   错误: 创建目录失败: permission denied
   ```
   解决: 检查目录权限，使用sudo或更改目录

4. **磁盘空间不足**
   ```
   错误: 写入文件失败: no space left on device
   ```
   解决: 清理磁盘空间或更改同步目录

### 调试模式

```bash
# 启用调试信息
./nib --debug

# 查看API响应详情
./nib list --debug
```

## 📊 性能优化建议

1. **并发数调整**
   - 网络较好: 20-50个并发
   - 网络一般: 10-20个并发
   - 网络较差: 5-10个并发

2. **目录结构**
   - 避免在根目录同步
   - 使用专门的备份目录
   - 定期清理不需要的文件

3. **定时同步**
   - 选择网络空闲时间
   - 避免与其他下载任务冲突
   - 监控磁盘空间使用

## 🔒 安全建议

1. **Token 管理**
   - 使用配置文件存储token
   - 设置适当的文件权限
   - 定期更换API token

2. **文件权限**
   - 设置适当的文件权限
   - 避免在公共目录同步
   - 定期检查文件完整性

3. **网络安全**
   - 使用HTTPS连接
   - 避免在公共网络使用
   - 监控网络流量

## 🎯 最佳实践

1. **首次使用**
   ```bash
   # 1. 复制配置文件模板
   cp nib.yaml.example nib.yaml
   
   # 2. 编辑配置文件，填入你的API Token
   nano nib.yaml
   
   # 3. 测试连接
   ./nib list
   
   # 4. 执行同步
   ./nib
   ```

2. **日常使用**
   ```bash
   # 直接同步
   ./nib
   
   # 查看状态
   ./nib list
   ```

3. **问题排查**
   ```bash
   # 启用调试
   ./nib --debug
   
   # 查看帮助
   ./nib --help
   ``` 