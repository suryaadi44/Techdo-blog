package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gorilla/mux"
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

func InitializeRouter() *mux.Router {
	router := mux.NewRouter()
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
		os.Getenv("PORT"),
		router,
	)

	server.Run()
}
