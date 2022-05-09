module github.com/suryaadi44/Techdo-blog

// +heroku goVersion go1.18
go 1.18

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.4.0
	golang.org/x/crypto v0.0.0-20220427172511-eb4f295cb31f
)

require github.com/codedius/imagekit-go v1.1.1

replace github.com/codedius/imagekit-go => github.com/suryaadi44/imagekit-go v1.1.2
