package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/suryaadi44/Techdo-blog/pkg/controller"
	"github.com/suryaadi44/Techdo-blog/pkg/database"
	"github.com/suryaadi44/Techdo-blog/pkg/server"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[ENV] Error loading .env file")
	}

	appName, present := os.LookupEnv("APP_NAME")

	if !present {
		log.Fatal("[ENV] Env variable not configure correctly")
	}

	log.Printf("[APP] app %s started\n", appName)
}

func InitializeRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.SetTrustedProxies(nil)
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency - param.Latency%time.Second
		}

		return fmt.Sprintf("%v [GIN] |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	return router
}

func InitializeDatabase() *sql.DB {
	return database.ConnectDB(
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
}

func main() {
	db := InitializeDatabase()
	router := InitializeRouter()

	controller.InitializeController(router, db)

	server := server.NewServer(
		os.Getenv("ADDRESS"),
		router,
	)

	server.Run()
}
