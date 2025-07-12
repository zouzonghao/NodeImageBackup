package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// ImageInfo 图片信息结构
type ImageInfo struct {
	ID       string `json:"image_id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Links    struct {
		Direct string `json:"direct"`
	} `json:"links"`
}

// APIResponse API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Images  []ImageInfo `json:"images,omitempty"`
}

// Config 配置结构
type Config struct {
	Token    string `yaml:"token"`
	LocalDir string `yaml:"dir"`
	APIBase  string `yaml:"api_base"`
	Workers  int    `yaml:"workers"`
}

// 全局配置
var config Config

// 主命令
var rootCmd = &cobra.Command{
	Use:   "nib",
	Short: "NodeImage Backup Tool - 高性能图片同步工具",
	Long: `NodeImage Backup Tool (nib) 是一个高性能的图片同步工具，
支持从NodeImage API同步图片到本地目录，实现单向同步功能。

默认执行同步操作，使用子命令可执行其他功能。`,
	RunE: runSync, // 默认执行同步
}

// 同步命令
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "同步远程图片到本地",
	RunE:  runSync,
}

// 列表命令
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出远程图片",
	RunE:  runList,
}

func init() {
	// 添加配置文件参数
	rootCmd.Flags().StringP("config", "c", "", "配置文件路径 (默认: nib.yaml 或 nib.yml)")
	syncCmd.Flags().StringP("config", "c", "", "配置文件路径 (默认: nib.yaml 或 nib.yml)")
	listCmd.Flags().StringP("config", "c", "", "配置文件路径 (默认: nib.yaml 或 nib.yml)")

	// 添加token参数
	rootCmd.Flags().StringP("token", "t", "", "API Token (可通过配置文件指定)")
	syncCmd.Flags().StringP("token", "t", "", "API Token (可通过配置文件指定)")
	listCmd.Flags().StringP("token", "t", "", "API Token (可通过配置文件指定)")

	// 添加目录参数
	rootCmd.Flags().StringP("dir", "d", "", "本地同步目录 (可通过配置文件指定，默认: 程序目录/images)")
	syncCmd.Flags().StringP("dir", "d", "", "本地同步目录 (可通过配置文件指定，默认: 程序目录/images)")
	listCmd.Flags().StringP("dir", "d", "", "本地同步目录 (可通过配置文件指定，默认: 程序目录/images)")

	// 添加并发数参数
	rootCmd.Flags().IntP("workers", "w", 0, "并发下载数量 (可通过配置文件指定，默认: 10)")
	syncCmd.Flags().IntP("workers", "w", 0, "并发下载数量 (可通过配置文件指定，默认: 10)")

	// 添加调试参数
	rootCmd.Flags().Bool("debug", false, "显示调试信息")
	syncCmd.Flags().Bool("debug", false, "显示调试信息")
	listCmd.Flags().Bool("debug", false, "显示调试信息")

	rootCmd.AddCommand(syncCmd, listCmd)
}

// 生成默认配置文件
func generateDefaultConfig() error {
	// 获取程序所在目录
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取程序路径: %v", err)
	}
	exeDir := filepath.Dir(exe)
	configPath := filepath.Join(exeDir, "nib.yaml")

	// 检查配置文件是否已存在
	if _, err := os.Stat(configPath); err == nil {
		return nil // 配置文件已存在，不需要生成
	}

	// 创建默认配置内容
	defaultConfig := `# NodeImage Backup Tool 配置文件
# 请将 YOUR_API_TOKEN_HERE 替换为您的真实 API Token

# API Token (必需)
token: "YOUR_API_TOKEN_HERE"

# 本地同步目录 (可选，默认: 程序目录/images)
dir: ""

# API 基础地址 (可选，默认: https://api.nodeimage.com)
api_base: "https://api.nodeimage.com"

# 并发下载数量 (可选，默认: 10)
workers: 10
`

	// 写入配置文件
	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("创建配置文件失败: %v", err)
	}

	fmt.Printf("📝 已自动生成配置文件: %s\n", configPath)
	fmt.Printf("🔑 请编辑配置文件，将 YOUR_API_TOKEN_HERE 替换为您的真实 API Token\n")
	fmt.Printf("💡 配置文件说明:\n")
	fmt.Printf("   - token: 您的 NodeImage API Token\n")
	fmt.Printf("   - dir: 本地同步目录 (留空使用默认目录)\n")
	fmt.Printf("   - workers: 并发下载数量\n")
	fmt.Printf("\n")

	return nil
}

// 根据token生成配置文件
func generateConfigWithToken(token string) error {
	// 获取程序所在目录
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取程序路径: %v", err)
	}
	exeDir := filepath.Dir(exe)
	configPath := filepath.Join(exeDir, "nib.yaml")

	// 检查配置文件是否已存在
	if _, err := os.Stat(configPath); err == nil {
		return nil // 配置文件已存在，不需要生成
	}

	// 创建包含token的配置内容
	configContent := fmt.Sprintf(`# NodeImage Backup Tool 配置文件
# 此配置文件已根据您提供的token自动生成

# API Token (必需)
token: "%s"

# 本地同步目录 (可选，默认: 程序目录/images)
dir: ""

# API 基础地址 (可选，默认: https://api.nodeimage.com)
api_base: "https://api.nodeimage.com"

# 并发下载数量 (可选，默认: 10)
workers: 10
`, token)

	// 写入配置文件
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("创建配置文件失败: %v", err)
	}

	fmt.Printf("📝 已根据您提供的token自动生成配置文件: %s\n", configPath)
	fmt.Printf("✅ 下次运行程序时无需再提供token参数\n")
	fmt.Printf("💡 如需修改配置，请编辑此配置文件\n\n")

	return nil
}

// 读取配置文件
func loadConfig(configPath string) (*Config, error) {
	paths := []string{}
	if configPath != "" {
		paths = append(paths, configPath)
	} else {
		paths = append(paths, "nib.yaml", "nib.yml")
	}

	configFound := false
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			configFound = true
			f, err := os.Open(path)
			if err != nil {
				return nil, fmt.Errorf("无法打开配置文件: %v", err)
			}
			defer f.Close()
			var cfg Config
			d := yaml.NewDecoder(f)
			if err := d.Decode(&cfg); err != nil {
				return nil, fmt.Errorf("解析配置文件失败: %v", err)
			}
			return &cfg, nil
		}
	}

	// 如果没有找到配置文件，自动生成一个
	if !configFound {
		if err := generateDefaultConfig(); err != nil {
			return nil, fmt.Errorf("生成配置文件失败: %v", err)
		}
	}

	return &Config{}, nil // 返回空配置，让用户先编辑配置文件
}

// 合并命令行参数和配置文件，命令行优先
type cliParams struct {
	token   string
	dir     string
	workers int
	config  string
	debug   bool
}

func mergeConfig(cfg *Config, cli cliParams) (Config, error) {
	final := *cfg

	// 如果命令行提供了token，优先使用
	if cli.token != "" {
		final.Token = cli.token

		// 如果配置文件中没有token（说明是第一次运行），自动生成配置文件
		if cfg.Token == "" {
			if err := generateConfigWithToken(cli.token); err != nil {
				return final, fmt.Errorf("自动生成配置文件失败: %v", err)
			}
		}
	}

	if cli.dir != "" {
		final.LocalDir = cli.dir
	}
	if cli.workers > 0 {
		final.Workers = cli.workers
	}
	if final.LocalDir == "" {
		// 获取程序所在目录
		exe, err := os.Executable()
		if err != nil {
			final.LocalDir = "./images"
		} else {
			exeDir := filepath.Dir(exe)
			final.LocalDir = filepath.Join(exeDir, "images")
		}
	}
	if final.Workers == 0 {
		final.Workers = 10
	}
	if final.APIBase == "" {
		final.APIBase = "https://api.nodeimage.com"
	}
	return final, nil
}

// 获取远程图片列表
func getRemoteImages(token string, debug bool) ([]ImageInfo, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", "https://api.nodeimage.com/api/v1/list", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("X-API-Key", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 调试：打印API原始响应
	if debug {
		fmt.Printf("[DEBUG] API响应原文: %s\n", string(body))
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API返回错误: %s", apiResp.Message)
	}

	return apiResp.Images, nil
}

// 获取本地图片列表
func getLocalImages(localDir string) (map[string]string, error) {
	localImages := make(map[string]string)

	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查是否为图片文件
		ext := strings.ToLower(filepath.Ext(path))
		if isImageFile(ext) {
			// 使用相对路径作为key
			relPath, err := filepath.Rel(localDir, path)
			if err != nil {
				return err
			}
			localImages[relPath] = path
		}

		return nil
	})

	return localImages, err
}

// 检查是否为图片文件
func isImageFile(ext string) bool {
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".avif", ".svg"}
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// 下载图片
func downloadImage(url, localPath string) error {
	client := &http.Client{Timeout: 60 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 创建目录
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建临时文件
	tmpFile := localPath + ".tmp"
	file, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()
	defer os.Remove(tmpFile) // 清理临时文件

	// 写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	// 原子性重命名
	if err := os.Rename(tmpFile, localPath); err != nil {
		return fmt.Errorf("重命名文件失败: %v", err)
	}

	return nil
}

// 计算文件MD5
func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// 用户确认函数
func askForConfirmation(prompt string) bool {
	fmt.Print(prompt + " (y/N 回车默认N): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取输入失败: %v，默认取消\n", err)
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// 运行同步
func runSync(cmd *cobra.Command, args []string) error {
	// 读取命令行参数
	cli := cliParams{}
	cli.token, _ = cmd.Flags().GetString("token")
	cli.dir, _ = cmd.Flags().GetString("dir")
	cli.workers, _ = cmd.Flags().GetInt("workers")
	cli.config, _ = cmd.Flags().GetString("config")
	cli.debug, _ = cmd.Flags().GetBool("debug")

	// 读取配置文件
	cfgFile, err := loadConfig(cli.config)
	if err != nil {
		return err
	}
	config, err = mergeConfig(cfgFile, cli)
	if err != nil {
		return err
	}

	if config.Token == "" || config.Token == "YOUR_API_TOKEN_HERE" {
		fmt.Printf("❌ 错误: 请通过配置文件或 -t 参数提供API Token\n")
		fmt.Printf("💡 如果配置文件已生成，请编辑配置文件并设置正确的 token\n")
		fmt.Printf("💡 或者使用命令行参数: ./nib -t YOUR_TOKEN\n")
		return fmt.Errorf("API Token 未配置")
	}

	// 新增：自动创建本地目录
	if err := os.MkdirAll(config.LocalDir, 0755); err != nil {
		return fmt.Errorf("创建本地目录失败: %v", err)
	}

	fmt.Printf("🔍 正在获取远程图片列表...\n")
	remoteImages, err := getRemoteImages(config.Token, cli.debug)
	if err != nil {
		return fmt.Errorf("获取远程图片失败: %v", err)
	}

	fmt.Printf("📁 正在扫描本地图片...\n")
	localImages, err := getLocalImages(config.LocalDir)
	if err != nil {
		return fmt.Errorf("扫描本地图片失败: %v", err)
	}

	fmt.Printf("📊 统计信息:\n")
	fmt.Printf("   远程图片: %d 张\n", len(remoteImages))
	fmt.Printf("   本地图片: %d 张\n", len(localImages))

	// 创建远程图片映射
	remoteMap := make(map[string]ImageInfo)
	for _, img := range remoteImages {
		remoteMap[img.Filename] = img
	}

	// 需要下载的图片
	var toDownload []ImageInfo
	for _, img := range remoteImages {
		if _, exists := localImages[img.Filename]; !exists {
			toDownload = append(toDownload, img)
		}
	}

	// 需要删除的本地图片
	var toDelete []string
	for filename, localPath := range localImages {
		if _, exists := remoteMap[filename]; !exists {
			toDelete = append(toDelete, localPath)
		}
	}

	fmt.Printf("\n🔄 同步计划:\n")
	fmt.Printf("   需要下载: %d 张\n", len(toDownload))
	fmt.Printf("   需要删除: %d 张\n", len(toDelete))

	if len(toDownload) == 0 && len(toDelete) == 0 {
		fmt.Printf("✅ 本地与远程已同步，无需操作\n")
		return nil
	}

	// 记录是否有实际执行的操作
	hasExecuted := false

	// 删除本地多余文件
	if len(toDelete) > 0 {
		fmt.Printf("\n🗑️  正在删除本地多余文件...\n")
		if !askForConfirmation(fmt.Sprintf("确认删除 %d 个本地文件?", len(toDelete))) {
			fmt.Printf("用户取消删除操作\n")
		} else {
			hasExecuted = true
			for _, filePath := range toDelete {
				if err := os.Remove(filePath); err != nil {
					fmt.Printf("   删除失败 %s: %v\n", filepath.Base(filePath), err)
				} else {
					fmt.Printf("   ✅ 已删除: %s\n", filepath.Base(filePath))
				}
			}
		}
	}

	// 并发下载
	if len(toDownload) > 0 {
		fmt.Printf("\n⬇️  正在下载图片...\n")
		if !askForConfirmation(fmt.Sprintf("确认下载 %d 个远程文件?", len(toDownload))) {
			fmt.Printf("用户取消下载操作\n")
		} else {
			hasExecuted = true
			semaphore := make(chan struct{}, config.Workers)
			var wg sync.WaitGroup

			for _, img := range toDownload {
				wg.Add(1)
				go func(img ImageInfo) {
					defer wg.Done()
					semaphore <- struct{}{}
					defer func() { <-semaphore }()

					localPath := filepath.Join(config.LocalDir, img.Filename)
					if err := downloadImage(img.Links.Direct, localPath); err != nil {
						fmt.Printf("   ❌ 下载失败 %s: %v\n", img.Filename, err)
					} else {
						fmt.Printf("   ✅ 已下载: %s (%d bytes)\n", img.Filename, img.Size)
					}
				}(img)
			}

			wg.Wait()
		}
	}

	// 只有实际执行了操作才显示同步完成
	if hasExecuted {
		fmt.Printf("\n🎉 同步完成!\n")
	}
	return nil
}

// 运行列表
func runList(cmd *cobra.Command, args []string) error {
	cli := cliParams{}
	cli.token, _ = cmd.Flags().GetString("token")
	cli.dir, _ = cmd.Flags().GetString("dir")
	cli.config, _ = cmd.Flags().GetString("config")
	cli.debug, _ = cmd.Flags().GetBool("debug")
	cfgFile, err := loadConfig(cli.config)
	if err != nil {
		return err
	}
	config, err = mergeConfig(cfgFile, cli)
	if err != nil {
		return err
	}
	if config.Token == "" || config.Token == "YOUR_API_TOKEN_HERE" {
		fmt.Printf("❌ 错误: 请通过配置文件或 -t 参数提供API Token\n")
		fmt.Printf("💡 如果配置文件已生成，请编辑配置文件并设置正确的 token\n")
		fmt.Printf("💡 或者使用命令行参数: ./nib list -t YOUR_TOKEN\n")
		return fmt.Errorf("API Token 未配置")
	}
	fmt.Printf("🔍 正在获取远程图片列表...\n")
	remoteImages, err := getRemoteImages(config.Token, cli.debug)
	if err != nil {
		return fmt.Errorf("获取远程图片失败: %v", err)
	}

	fmt.Printf("\n📋 远程图片列表 (%d 张):\n", len(remoteImages))
	for i, img := range remoteImages {
		fmt.Printf("%3d. %s (%d bytes)\n", i+1, img.Filename, img.Size)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
