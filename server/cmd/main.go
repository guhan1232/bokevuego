package main

import (
	"database/sql"
	"log"

	"bokeui/internal/model"
	"bokeui/internal/router"
)

func main() {
	db, err := model.InitDB("data.db")
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	defer db.Close()

	model.AutoMigrate(db)

	r := router.SetupRouter(db)
	log.Println("服务启动在 http://localhost:9088")
	log.Println("管理后台: http://localhost:9088/admin")
	log.Println("前台首页: http://localhost:9088")
	if err := r.Run(":9088"); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}

var _ *sql.DB
