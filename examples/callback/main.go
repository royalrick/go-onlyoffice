package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/royalrick/go-onlyoffice"
	"github.com/royalrick/go-onlyoffice/models"
)

var (
	savedFiles = make(map[string]string) // key -> saved path
	mu         sync.Mutex
)

func main() {
	fmt.Println("=== OnlyOffice 回调处理与文档保存示例 ===")

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

	// 创建存储目录
	storageDir := "./storage"
	savedDir := filepath.Join(storageDir, "saved")
	if err := os.MkdirAll(savedDir, 0755); err != nil {
		log.Fatalf("创建存储目录失败: %v", err)
	}

	// 创建示例文档
	sampleFile := filepath.Join(storageDir, "document.docx")
	if _, err := os.Stat(sampleFile); os.IsNotExist(err) {
		content := []byte("请编辑这个文档，然后保存。\n\n编辑后的内容将通过回调自动保存到服务器。")
		if err := os.WriteFile(sampleFile, content, 0644); err != nil {
			log.Fatalf("创建示例文件失败: %v", err)
		}
		fmt.Printf("✓ 已创建示例文档: %s\n", sampleFile)
	}

	// 设置路由
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(storageDir))))

	// 编辑器页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveEditorPage(w, r, client)
	})

	// 编辑器配置 API
	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		serveEditorConfig(w, r, client)
	})

	// 回调处理器 - 使用新的 CallbackHandler API
	handler := client.CallbackHandler(onlyoffice.CallbackHandlers{
		OnEditing: func(cb *models.Callback) error {
			log.Printf("📝 文档正在编辑 - Key: %s, Users: %v", cb.Key, cb.Users)
			return nil
		},
		OnSave: func(cb *models.Callback) error {
			log.Printf("💾 文档准备保存 - Key: %s, URL: %s", cb.Key, cb.Url)
			return saveDocument(cb, savedDir)
		},
		OnSaveError: func(cb *models.Callback) error {
			log.Printf("❌ 文档保存出错 - Key: %s", cb.Key)
			return nil
		},
		OnClose: func(cb *models.Callback) error {
			log.Printf("🚪 文档已关闭(无修改) - Key: %s", cb.Key)
			return nil
		},
		OnForceSave: func(cb *models.Callback) error {
			log.Printf("⚡ 文档强制保存 - Key: %s, URL: %s", cb.Key, cb.Url)
			return saveDocument(cb, savedDir)
		},
		OnCorrupt: func(cb *models.Callback) error {
			log.Printf("⚠️  文档已损坏 - Key: %s", cb.Key)
			return nil
		},
	})

	http.Handle("/callback", handler)

	// 查看已保存的文档
	http.HandleFunc("/saved", func(w http.ResponseWriter, r *http.Request) {
		serveSavedList(w, r)
	})

	// 启动服务器
	port := getEnv("PORT", "8083")
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("\n✓ 服务器启动成功\n")
	fmt.Printf("  编辑器页面: http://localhost:%s\n", port)
	fmt.Printf("  回调地址: http://localhost:%s/callback\n", port)
	fmt.Printf("  已保存文档: http://localhost:%s/saved\n\n", port)
	fmt.Println("💡 提示: 编辑文档并保存，服务器将自动接收回调并下载保存")

	log.Fatal(http.ListenAndServe(addr, nil))
}

func saveDocument(cb *models.Callback, savedDir string) error {
	if cb.Url == "" {
		return fmt.Errorf("empty download URL")
	}

	// 下载文档
	log.Printf("  开始下载文档: %s", cb.Url)
	resp, err := http.Get(cb.Url)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 保存文档
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("document-%s.docx", timestamp)
	savePath := filepath.Join(savedDir, filename)

	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	written, err := io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}

	log.Printf("  ✓ 文档已保存: %s (%d 字节)", savePath, written)

	// 记录保存的文件
	mu.Lock()
	savedFiles[cb.Key] = savePath
	mu.Unlock()

	return nil
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
    <title>OnlyOffice 回调处理示例</title>
    <style>
        body {
            margin: 0;
            padding: 20px;
            font-family: Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
        }
        h1 {
            color: #333;
            margin-bottom: 10px;
        }
        .subtitle {
            color: #666;
            font-size: 14px;
            margin-bottom: 25px;
        }
        #editor {
            width: 100%;
            height: 600px;
            border: 2px solid #667eea;
            border-radius: 8px;
            overflow: hidden;
        }
        .info {
            background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 25px;
            border-left: 4px solid #667eea;
        }
        .info p {
            margin: 8px 0;
            color: #333;
            line-height: 1.6;
        }
        .info strong {
            color: #667eea;
        }
        .saved-link {
            display: inline-block;
            margin-top: 15px;
            padding: 10px 20px;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 6px;
            transition: background 0.3s;
        }
        .saved-link:hover {
            background: #764ba2;
        }
    </style>
    <script src="http://localhost/web-apps/apps/api/documents/api.js"></script>
</head>
<body>
    <div class="container">
        <h1>💾 OnlyOffice 回调处理示例</h1>
        <div class="subtitle">编辑文档并保存，服务器将自动接收回调并下载保存</div>

        <div class="info">
            <p><strong>📋 功能说明:</strong></p>
            <p>• 在下方编辑器中修改文档内容</p>
            <p>• 点击保存按钮（Ctrl+S 或工具栏保存）</p>
            <p>• OnlyOffice 将发送回调通知到服务器</p>
            <p>• 服务器自动下载并保存修改后的文档</p>
            <p>• 查看服务器控制台日志了解回调处理过程</p>
            <a href="/saved" class="saved-link" target="_blank">📂 查看已保存的文档</a>
        </div>

        <div id="editor"></div>
    </div>

    <script>
        fetch('/api/config')
            .then(response => response.json())
            .then(config => {
                console.log('编辑器配置:', config);
                new DocsAPI.DocEditor("editor", config);
            })
            .catch(error => {
                console.error('加载配置失败:', error);
                document.getElementById('editor').innerHTML =
                    '<p style="color: red; padding: 20px;">错误: 无法加载编辑器配置</p>';
            });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func serveEditorConfig(w http.ResponseWriter, r *http.Request, client *onlyoffice.Client) {
	host := r.Host
	fileURL := fmt.Sprintf("http://%s/files/document.docx", host)
	callbackURL := fmt.Sprintf("http://%s/callback", host)

	params := models.EditorParams{
		Filename:    "document.docx",
		Mode:        "edit",
		Language:    "zh-CN",
		UserId:      "user123",
		UserName:    "测试用户",
		UserEmail:   "user@example.com",
		CallbackUrl: callbackURL,
		CanEdit:     true,
		CanDownload: true,
	}

	cfg, err := client.BuildEditorConfig(params, fileURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("生成配置失败: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func serveSavedList(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>已保存的文档</title>
    <style>
        body {
            margin: 0;
            padding: 20px;
            font-family: Arial, sans-serif;
            background: #f5f5f5;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            margin-bottom: 20px;
        }
        .file-list {
            list-style: none;
            padding: 0;
        }
        .file-item {
            padding: 15px;
            margin: 10px 0;
            background: #f9f9f9;
            border-radius: 6px;
            border-left: 4px solid #667eea;
        }
        .file-key {
            color: #666;
            font-size: 12px;
            margin-top: 5px;
        }
        .empty {
            text-align: center;
            color: #999;
            padding: 40px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📂 已保存的文档</h1>`

	if len(savedFiles) == 0 {
		html += `<div class="empty">暂无保存的文档。请编辑文档并保存。</div>`
	} else {
		html += `<ul class="file-list">`
		for key, path := range savedFiles {
			html += fmt.Sprintf(`
            <li class="file-item">
                <div><strong>%s</strong></div>
                <div class="file-key">Key: %s</div>
            </li>`, filepath.Base(path), key)
		}
		html += `</ul>`
	}

	html += `
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
