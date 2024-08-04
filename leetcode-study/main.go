package main

import (
	_ "leetcode/dao"
	"leetcode/routes"
	"log"
)

func main() {

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	log.Println("Server started at http://localhost:8080")
	err := r.Run(":8080") // 使用 Gin 的 Run 方法
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
