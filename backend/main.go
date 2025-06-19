package main

import (
	"easyinvesting/config"
	"easyinvesting/pkg/models"
	"easyinvesting/pkg/server"
	"fmt"
)

func main() {
	models.Migrate()
	defer func() {
		db, _ := config.DB.DB()
		if err := db.Close(); err != nil {
			fmt.Printf("Failed to close database connection: %v\n", err)
		}
	}()
	fmt.Println("Hello")

	s := server.NewServer(8000)
	s.Start()
}
