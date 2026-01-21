# OnlyOffice Go SDK 示例

本目录包含 OnlyOffice Go SDK 的完整功能示例，展示了文档格式转换、在线编辑、回调处理和版本历史管理等核心功能。

## 前置要求

1. **OnlyOffice Document Server** 必须正在运行
   ```bash
   # 使用 Docker 快速启动
   docker run -i -t -d -p 80:80 onlyoffice/documentserver
   ```

2. **配置环境变量**（可选）
   ```bash
   export ONLYOFFICE_URL="http://localhost"
   export JWT_SECRET="your-secret-key"
   export JWT_ENABLED="false"  # 生产环境建议设置为 true
   ```

## 示例列表

### 1. 文档格式转换 (`convert/`)

展示如何使用 OnlyOffice 转换服务将文档从一种格式转换为另一种格式。

**功能:**
- 自动创建示例文档 (TXT)
- 启动本地文件服务器
- 调用转换 API (TXT → DOCX, TXT → PDF)
- 下载并保存转换后的文件

**运行:**
```bash
cd convert
go run main.go
```

**输出:**
- `storage/sample.txt` - 原始文档
- `storage/converted.docx` - 转换后的 DOCX
- `storage/converted.pdf` - 转换后的 PDF

---

### 2. Web 编辑器集成 (`editor/`)

展示如何在网页中嵌入 OnlyOffice 编辑器，实现在线文档编辑。

**功能:**
- 生成编辑器配置（JWT 签名）
- 提供 HTML 编辑器界面
- 文件服务（提供文档下载）

**运行:**
```bash
cd editor
go run main.go
```

**访问:**
打开浏览器访问 `http://localhost:8082`

**提示:**
- 可在浏览器中直接编辑文档
- 支持实时协作编辑
- 此示例不包含回调处理

---

### 3. 回调处理与文档保存 (`callback/`)

展示完整的编辑-保存流程，包括接收 OnlyOffice 回调并自动保存修改后的文档。

**功能:**
- 完整的在线编辑器界面
- 实现所有回调状态处理:
  - `OnEditing` (状态 1) - 文档正在编辑
  - `OnSave` (状态 2) - 文档准备保存
  - `OnSaveError` (状态 3) - 保存错误
  - `OnClose` (状态 4) - 文档关闭（无修改）
  - `OnForceSave` (状态 6) - 强制保存
  - `OnCorrupt` (状态 7) - 文档损坏
- 自动下载并保存修改后的文档
- 展示已保存文档列表

**运行:**
```bash
cd callback
go run main.go
```

**访问:**
- 编辑器: `http://localhost:8083`
- 已保存文档: `http://localhost:8083/saved`

**工作流程:**
1. 在编辑器中修改文档
2. 点击保存（或按 Ctrl+S）
3. OnlyOffice 发送回调到 `/callback`
4. 服务器自动下载并保存到 `storage/saved/`
5. 查看控制台日志了解回调处理过程

---

### 4. 版本历史管理 (`history/`)

展示如何记录和查询文档的版本历史。

**功能:**
- 创建文档历史记录
- 查询历史版本列表
- 统计版本数量
- 展示版本存储结构

**运行:**
```bash
cd history
go run main.go
```

**输出:**
```
=== OnlyOffice 文档版本历史管理示例 ===

--- 示例 1: 创建文档历史版本 ---
文档键: abc123...

1. 模拟第一次保存...
   ✓ 版本 1 已保存
2. 模拟第二次保存...
   ✓ 版本 2 已保存
3. 模拟第三次保存...
   ✓ 版本 3 已保存

--- 示例 2: 查询历史版本 ---

找到 3 个历史版本:

版本 1:
  文件键: abc123...
  创建时间: 2026-01-21 10:30:00
  用户: 张三 (ID: user1)
  修改数: 2
    - 修改时间: 2026-01-21 10:30:00, 用户: 张三
    - 修改时间: 2026-01-21 10:31:00, 用户: 张三
...
```

**存储结构:**
```
storage/
└── .history/
    └── {document-key}/
        └── changes.json
```

---

## 目录结构

```
examples/
├── README.md           # 本文件
├── convert/            # 格式转换示例
│   ├── main.go
│   └── storage/        # 文件存储目录
├── editor/             # Web 编辑器示例
│   ├── main.go
│   └── storage/
├── callback/           # 回调处理示例
│   ├── main.go
│   └── storage/
│       └── saved/      # 自动保存的文档
└── history/            # 版本历史示例
    ├── main.go
    └── storage/
        └── .history/   # 历史版本存储
```

## 常见问题

### Q: 运行示例时提示连接失败？
A: 确保 OnlyOffice Document Server 正在运行。使用 Docker 快速启动:
```bash
docker run -i -t -d -p 80:80 onlyoffice/documentserver
```

### Q: 编辑器无法加载？
A: 检查浏览器控制台，确保可以访问 `http://localhost/web-apps/apps/api/documents/api.js`

### Q: 回调一直不生效？
A: 确保:
1. OnlyOffice Document Server 可以访问你的回调 URL
2. 如果使用 localhost，Document Server 也必须在同一台机器上
3. 如果启用 JWT，确保 `JWT_SECRET` 配置一致

### Q: 如何在生产环境使用？
A:
1. 设置 `JWT_ENABLED=true` 启用 JWT 验证
2. 使用强随机 `JWT_SECRET`
3. 使用 HTTPS 协议
4. 确保 Document Server 可以通过外网访问回调 URL

## 下一步

- 阅读 [OnlyOffice API 文档](https://api.onlyoffice.com/)
- 查看项目根目录的 `CLAUDE.md` 了解架构设计
- 运行 `go test ./...` 执行单元测试

## 许可证

与主项目相同。
