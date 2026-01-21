package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/royalrick/go-onlyoffice"
	"github.com/royalrick/go-onlyoffice/models"
)

func main() {
	fmt.Println("=== OnlyOffice 文档版本历史管理示例 ===")

	// 配置客户端
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
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("创建存储目录失败: %v", err)
	}

	fmt.Println("\n--- 示例 1: 创建文档历史版本 ---")

	// 模拟文档的多次编辑和保存
	documentKey, err := client.GenerateFileHash("test-document.docx")
	if err != nil {
		log.Fatalf("生成文档键失败: %v", err)
	}

	fmt.Printf("文档键: %s\n\n", documentKey)

	// 模拟第一次保存
	fmt.Println("1. 模拟第一次保存...")
	callback1 := createMockCallback(documentKey, 1, "user1", "张三")
	if err := client.CreateHistory(callback1, storageDir); err != nil {
		log.Printf("创建历史失败: %v", err)
	} else {
		fmt.Println("   ✓ 版本 1 已保存")
	}

	time.Sleep(time.Second)

	// 模拟第二次保存
	fmt.Println("2. 模拟第二次保存...")
	callback2 := createMockCallback(documentKey, 2, "user2", "李四")
	if err := client.CreateHistory(callback2, storageDir); err != nil {
		log.Printf("创建历史失败: %v", err)
	} else {
		fmt.Println("   ✓ 版本 2 已保存")
	}

	time.Sleep(time.Second)

	// 模拟第三次保存
	fmt.Println("3. 模拟第三次保存...")
	callback3 := createMockCallback(documentKey, 3, "user1", "张三")
	if err := client.CreateHistory(callback3, storageDir); err != nil {
		log.Printf("创建历史失败: %v", err)
	} else {
		fmt.Println("   ✓ 版本 3 已保存")
	}

	fmt.Println("\n--- 示例 2: 查询历史版本 ---")

	// 查询版本历史
	versions, err := client.GetHistory("test-document.docx", storageDir)
	if err != nil {
		log.Printf("查询历史失败: %v", err)
	} else {
		fmt.Printf("\n找到 %d 个历史版本:\n\n", len(versions))
		for i, version := range versions {
			fmt.Printf("版本 %d:\n", i+1)
			fmt.Printf("  文件键: %s\n", version.Key)
			fmt.Printf("  创建时间: %s\n", version.Created.Format("2006-01-02 15:04:05"))
			if version.User != nil {
				fmt.Printf("  用户: %s (ID: %s)\n", version.User.Name, version.User.Id)
			}
			fmt.Printf("  修改数: %d\n", len(version.ChangesData))
			if len(version.ChangesData) > 0 {
				for _, change := range version.ChangesData {
					fmt.Printf("    - 修改时间: %s, 用户: %s\n", change.Created, change.User.Name)
				}
			}
			fmt.Println()
		}
	}

	fmt.Println("\n--- 示例 3: 统计版本数量 ---")

	count := client.CountVersion(storageDir)
	fmt.Printf("共有 %d 个版本记录\n", count)

	fmt.Println("\n--- 示例 4: 查看存储结构 ---")

	historyDir := filepath.Join(storageDir, ".history")
	if _, err := os.Stat(historyDir); err == nil {
		fmt.Printf("\n历史记录存储在: %s\n", historyDir)
		fmt.Println("目录结构:")

		entries, err := os.ReadDir(historyDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					fmt.Printf("  📁 %s/\n", entry.Name())
					subDir := filepath.Join(historyDir, entry.Name())
					subEntries, _ := os.ReadDir(subDir)
					for _, subEntry := range subEntries {
						fmt.Printf("      📄 %s\n", subEntry.Name())
					}
				}
			}
		}
	}

	fmt.Println("\n=== 示例完成 ===")
	fmt.Printf("历史文件保存在: %s\n", historyDir)
}

func createMockCallback(key string, versionNum int, userId, userName string) models.Callback {
	now := time.Now()

	return models.Callback{
		Key:    key,
		Status: 2, // 保存状态
		Url:    fmt.Sprintf("http://localhost/file-%s-v%d.docx", key, versionNum),
		History: models.History{
			ServerVersion: fmt.Sprintf("7.%d.0", versionNum),
			Key:           key,
			Created:       now.Format("2006-01-02 15:04:05"),
			User: &models.User{
				Id:   userId,
				Name: userName,
			},
			Changes: []models.Change{
				{
					Created: now.Format("2006-01-02 15:04:05"),
					User: models.User{
						Id:   userId,
						Name: userName,
					},
				},
				{
					Created: now.Add(time.Minute).Format("2006-01-02 15:04:05"),
					User: models.User{
						Id:   userId,
						Name: userName,
					},
				},
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// printJSON 打印结构化 JSON（用于调试）
func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
