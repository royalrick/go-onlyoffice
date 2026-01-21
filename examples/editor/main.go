package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/royalrick/go-onlyoffice"
	"github.com/royalrick/go-onlyoffice/models"
)

func main() {
	fmt.Println("=== OnlyOffice Web 编辑器示例 ===")

	// 配置 OnlyOffice 客户端
	config := &onlyoffice.Config{
		DocumentServerURL: getEnv("ONLYOFFICE_URL", "http://localhost"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
		JWTEnabled:        getEnv("JWT_ENABLED", "false") == "true",
	}

	client, err := onlyoffice.NewClient(config)
	if err != nil {
		log.Fatalf("初始化客户端失败: %v", err)
	}

	// 创建存储目录和示例文档
	storageDir := "./storage"
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("创建存储目录失败: %v", err)
	}

	sampleFile := filepath.Join(storageDir, "document.docx")
	if _, err := os.Stat(sampleFile); os.IsNotExist(err) {
		// 创建一个简单的文本文件作为示例
		content := []byte("欢迎使用 OnlyOffice 编辑器！\n\n这是一个示例文档。")
		if err := os.WriteFile(sampleFile, content, 0644); err != nil {
			log.Fatalf("创建示例文件失败: %v", err)
		}
		fmt.Printf("✓ 已创建示例文档: %s\n", sampleFile)
	}

	// 设置路由
	// 1. 文件服务
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(storageDir))))

	// 2. 编辑器页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveEditorPage(w, r, client)
	})

	// 3. 编辑器配置 API
	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		serveEditorConfig(w, r, client)
	})

	// 启动服务器
	port := getEnv("PORT", "8082")
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("\n✓ Web 服务器启动成功\n")
	fmt.Printf("  访问编辑器: http://localhost:%s\n", port)
	fmt.Printf("  文档文件: http://localhost:%s/files/document.docx\n\n", port)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func serveEditorPage(w http.ResponseWriter, r *http.Request, client *onlyoffice.Client) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>OnlyOffice 编辑器示例</title>
    <style>
        body {
            margin: 0;
            padding: 20px;
            font-family: Arial, sans-serif;
            background: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            margin-bottom: 20px;
        }
        #editor {
            width: 100%;
            height: 600px;
            border: 1px solid #ddd;
        }
        .info {
            background: #e3f2fd;
            padding: 15px;
            border-radius: 4px;
            margin-bottom: 20px;
        }
        .info p {
            margin: 5px 0;
            color: #1976d2;
        }
    </style>
    <script src="http://localhost/web-apps/apps/api/documents/api.js"></script>
</head>
<body>
    <div class="container">
        <h1>🚀 OnlyOffice Web 编辑器示例</h1>
        <div class="info">
            <p><strong>提示:</strong> 此示例展示如何在网页中嵌入 OnlyOffice 编辑器</p>
            <p><strong>功能:</strong> 在线编辑文档、协作、评论等</p>
            <p><strong>配置:</strong> 点击页面即可加载编辑器配置</p>
        </div>
        <div id="editor"></div>
    </div>

    <script>
        // 从 API 获取编辑器配置
        fetch('/api/config')
            .then(response => response.json())
            .then(config => {
                console.log('编辑器配置:', config);

                // 初始化 OnlyOffice 编辑器
                new DocsAPI.DocEditor("editor", config);
            })
            .catch(error => {
                console.error('加载配置失败:', error);
                document.getElementById('editor').innerHTML =
                    '<p style="color: red; padding: 20px;">错误: 无法加载编辑器配置。请确保 OnlyOffice Document Server 正在运行。</p>';
            });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func serveEditorConfig(w http.ResponseWriter, r *http.Request, client *onlyoffice.Client) {
	// 构建编辑器配置
	params := models.EditorParams{
		Filename:    "document.docx",
		Mode:        "edit",
		Language:    "zh-CN",
		UserId:      "user-" + randString(6),
		UserName:    "测试用户",
		UserEmail:   "user@example.com",
		CallbackUrl: "", // 简单示例不需要回调
		CanEdit:     true,
		CanDownload: true,
	}

	// 获取当前主机地址
	host := r.Host
	fileURL := fmt.Sprintf("http://%s/files/document.docx", host)

	cfg, err := client.BuildEditorConfig(params, fileURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("生成配置失败: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
