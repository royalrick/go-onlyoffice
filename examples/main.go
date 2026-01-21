package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== OnlyOffice Go SDK 示例 ===")
	fmt.Println()
	fmt.Println("本目录包含多个独立的示例程序，每个展示不同的功能：")
	fmt.Println()
	fmt.Println("1. 文档格式转换 (convert/)")
	fmt.Println("   cd convert && go run main.go")
	fmt.Println("   展示如何将文档从一种格式转换为另一种格式")
	fmt.Println()
	fmt.Println("2. Web 编辑器集成 (editor/)")
	fmt.Println("   cd editor && go run main.go")
	fmt.Println("   展示如何在网页中嵌入 OnlyOffice 编辑器")
	fmt.Println()
	fmt.Println("3. 回调处理与文档保存 (callback/)")
	fmt.Println("   cd callback && go run main.go")
	fmt.Println("   展示完整的编辑-保存流程和回调处理")
	fmt.Println()
	fmt.Println("4. 版本历史管理 (history/)")
	fmt.Println("   cd history && go run main.go")
	fmt.Println("   展示如何记录和查询文档版本历史")
	fmt.Println()
	fmt.Println("详细说明请查看: README.md")
	fmt.Println()
	fmt.Println("前置要求:")
	fmt.Println("  - OnlyOffice Document Server 必须正在运行")
	fmt.Println("  - 使用 Docker 快速启动:")
	fmt.Println("    docker run -i -t -d -p 80:80 onlyoffice/documentserver")
}
