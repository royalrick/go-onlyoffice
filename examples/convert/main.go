package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/royalrick/go-onlyoffice"
)

func main() {
	fmt.Println("=== OnlyOffice 文档格式转换示例 ===")

	// 配置 OnlyOffice 客户端
	config := &onlyoffice.Config{
		DocumentServerURL: getEnv("ONLYOFFICE_URL", "http://localhost:80"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
		JWTEnabled:        getEnv("JWT_ENABLED", "false") == "true",
	}

	client, err := onlyoffice.NewClient(config)
	if err != nil {
		log.Fatalf("初始化客户端失败: %v", err)
	}

	// 检查 OnlyOffice 服务器
	fmt.Println("\n--- 检查 OnlyOffice Document Server ---")
	if err := checkOnlyOfficeServer(config.DocumentServerURL); err != nil {
		log.Fatalf("❌ OnlyOffice 服务器不可用: %v\n\n提示: 使用 Docker 启动服务器:\n  docker run -i -t -d -p 80:80 onlyoffice/documentserver\n", err)
	}
	fmt.Println("✓ OnlyOffice 服务器运行正常")

	// 获取本机 IP 地址
	hostIP := getEnv("HOST_IP", "")
	if hostIP == "" {
		hostIP = getLocalIP()
		if hostIP == "" {
			log.Fatal("❌ 无法获取本机 IP 地址\n\n提示: 请设置环境变量 HOST_IP，例如:\n  export HOST_IP=192.168.1.100\n")
		}
	}
	fmt.Printf("✓ 使用主机 IP: %s\n", hostIP)

	// 启动简单的文件服务器来提供源文件
	storageDir := "./storage"
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("创建存储目录失败: %v", err)
	}

	// 创建示例文档
	sampleFile := filepath.Join(storageDir, "sample.txt")
	if _, err := os.Stat(sampleFile); os.IsNotExist(err) {
		content := []byte("这是一个测试文档。\n\n这个文件将被转换为不同的格式。\n\nOnlyOffice Document Server 支持多种格式转换。")
		if err := os.WriteFile(sampleFile, content, 0644); err != nil {
			log.Fatalf("创建示例文件失败: %v", err)
		}
		fmt.Printf("✓ 已创建示例文件: %s\n", sampleFile)
	}

	// 启动文件服务器（在后台）
	port := "8081"
	fileServer := http.FileServer(http.Dir(storageDir))
	go func() {
		http.Handle("/files/", http.StripPrefix("/files/", fileServer))
		log.Printf("文件服务器启动在 http://%s:%s/files/", hostIP, port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Printf("文件服务器错误: %v", err)
		}
	}()

	// 等待服务器启动
	fmt.Println("\n--- 启动文件服务器 ---")
	time.Sleep(time.Second)
	fileURL := fmt.Sprintf("http://%s:%s/files/sample.txt", hostIP, port)
	if err := checkFileServer(fileURL); err != nil {
		log.Fatalf("❌ 文件服务器启动失败: %v", err)
	}
	fmt.Printf("✓ 文件服务器就绪: %s\n", fileURL)

	// 执行文档转换
	fmt.Println("\n--- 开始转换文档 ---")

	// 示例 1: TXT -> DOCX
	fmt.Println("\n1. 转换 TXT 为 DOCX...")
	result1, err := convertDocument(client, fileURL, "txt", "docx")
	if err != nil {
		log.Printf("❌ 转换失败: %v", err)
	} else {
		if err := downloadAndSave(result1.FileURL, filepath.Join(storageDir, "converted.docx")); err != nil {
			log.Printf("❌ 下载失败: %v", err)
		}
	}

	// 示例 2: TXT -> PDF
	fmt.Println("\n2. 转换 TXT 为 PDF...")
	result2, err := convertDocument(client, fileURL, "txt", "pdf")
	if err != nil {
		log.Printf("❌ 转换失败: %v", err)
	} else {
		if err := downloadAndSave(result2.FileURL, filepath.Join(storageDir, "converted.pdf")); err != nil {
			log.Printf("❌ 下载失败: %v", err)
		}
	}

	fmt.Println("\n=== 转换完成 ===")
	fmt.Printf("转换后的文件保存在: %s\n", storageDir)
}

func convertDocument(client *onlyoffice.Client, fileURL, fromExt, toExt string) (*onlyoffice.ConvertResult, error) {
	key, err := client.GenerateFileHash(fileURL)
	if err != nil {
		return nil, fmt.Errorf("生成文件键失败: %w", err)
	}

	opts := onlyoffice.ConvertOptions{
		DocumentURL: fileURL,
		FromExt:     fromExt,
		ToExt:       toExt,
		DocumentKey: key,
		Title:       fmt.Sprintf("convert-%s-to-%s", fromExt, toExt),
		Async:       false,
	}

	fmt.Printf("  调用转换 API: %s -> %s\n", fromExt, toExt)
	result, err := client.ConvertDocument(opts)
	if err != nil {
		return nil, fmt.Errorf("转换请求失败: %w", err)
	}

	if result.Error != 0 {
		return nil, fmt.Errorf("转换错误，错误码: %d", result.Error)
	}

	fmt.Printf("  ✓ 转换成功，进度: %d%%\n", result.Percent)
	fmt.Printf("  ✓ 转换结果 URL: %s\n", result.FileURL)

	return result, nil
}

func downloadAndSave(url, savePath string) error {
	fmt.Printf("  下载文件到: %s\n", savePath)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	written, err := file.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	fmt.Printf("  ✓ 已保存 %d 字节\n", written)
	return nil
}

func checkOnlyOfficeServer(serverURL string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(serverURL + "/healthcheck")
	if err != nil {
		return fmt.Errorf("无法连接到服务器: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务器返回错误状态: %d", resp.StatusCode)
	}
	return nil
}

func checkFileServer(fileURL string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fileURL)
	if err != nil {
		return fmt.Errorf("无法访问文件: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("文件服务器返回错误状态: %d", resp.StatusCode)
	}
	return nil
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
