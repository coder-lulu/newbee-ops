package main

import (
	"context"
	"fmt"
	"log"

	"github.com/coder-lulu/newbee-ops-rpc/ent"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(192.168.26.130:3306)/newbee?parseTime=True&charset=utf8mb4"
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	// 自动迁移schema
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	fmt.Println("✅ Worker tables created successfully!")
}
