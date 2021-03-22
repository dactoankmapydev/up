package main

import (
	"backup/conn"
	"backup/handler"
	"backup/repository/repo_impl"
	"backup/router"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Not environment variable")
	}
}

func main() {
	// redis details
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	// connect redis
	client := &conn.RedisDB{
		Host: redisHost,
		Port: redisPort,
	}
	client.NewRedisDB()

	e := echo.New()
	uploadHandler := handler.UploadHandler{
		UploadRepo: repo_impl.NewBackupRepo(client),
	}

	api := router.API{
		Echo: e,
		UploadHandler: uploadHandler,
	}
	api.SetupRouter()
	e.Logger.Fatal(e.Start(":3000"))
}