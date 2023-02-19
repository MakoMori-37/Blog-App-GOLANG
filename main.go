package main

import (
	"goBlogApp/config"
	"goBlogApp/migrations"
	"goBlogApp/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()
	defer config.CloseDB()
	migrations.Migrate()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")

	r := gin.Default()
	r.Use(cors.New(corsConfig))
	r.Static("/uploads", "./uploads")

	uploadDirs := [...]string{"atticles", "users"}
	for _, dir := range uploadDirs {
		os.MkdirAll("uploads/" + dir, 0755)
		
	}

	routes.Serve(r)
	r.Run()
}
