package main

import (
	"easyinvesting/config"
	"easyinvesting/pkg/models"
	"easyinvesting/pkg/server"
	"fmt"
	"os"
	"strconv"
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

	port := 8000
	if portEnv := os.Getenv("EASYINVESTING_PORT"); portEnv != "" {
		var err error
		port, err = strconv.Atoi(portEnv)
		if err != nil || port <= 0 {
			fmt.Printf("EASYINVESTING_PORT is a invalid number: %s\n", portEnv)
			os.Exit(1)
		}
	}
	s := server.NewServer(port)
	s.Start()
}
