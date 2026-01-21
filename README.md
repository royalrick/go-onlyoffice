# go-onlyoffice

OnlyOffice 文档服务器的 Go 语言 SDK，提供核心功能集成。

## 特性

- 生成文档编辑器配置
- 文档格式转换
- JWT 认证支持
- 历史版本管理
- 回调数据处理

## 安装

```bash
go get github.com/royalrick/go-onlyoffice
```

## 使用

### 初始化客户端

```go
import "github.com/royalrick/go-onlyoffice"

config := &onlyoffice.Config{
    DocumentServerURL: "https://doc-server.com",
    JWTSecret:        "your-secret",
    JWTEnabled:       true,
}

client, err := onlyoffice.NewClient(config)
```

### 生成编辑器配置

```go
import "github.com/royalrick/go-onlyoffice/models"

params := models.EditorParams{
    Filename:    "document.docx",
    Mode:        "edit",
    Language:    "zh-CN",
    UserId:      "user123",
    UserName:    "John Doe",
    UserEmail:   "john@example.com",
    CallbackUrl: "https://your-server.com/callback",
    CanEdit:     true,
    CanDownload: true,
}

fileURL := "https://your-server.com/storage/document.docx"
cfg, err := client.BuildEditorConfig(params, fileURL)
```

### 文档转换

```go
opts := onlyoffice.ConvertOptions{
    DocumentURL: "https://example.com/document.docx",
    ToExt:       "pdf",
    DocumentKey: "unique-key",
    Async:       false,
}

result, err := client.ConvertDocument(opts)
```

### 处理回调

```go
jsonData := []byte(`{"status": 2, "key": "key", "url": "url"}`)
tokenHeader := r.Header.Get("Authorization")

callback, err := client.ParseCallback(jsonData, tokenHeader)
if err := client.ValidateCallback(callback); err != nil {
    return err
}

downloadURL, err := client.GetDownloadURL(callback)
```

### 历史版本管理

```go
err := client.CreateHistory(callback, "./storage")
versions, err := client.GetHistory("document.docx", "./storage")
count := client.CountVersion("./storage")
```

## 许可证

Apache License 2.0
