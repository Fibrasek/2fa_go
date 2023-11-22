package main

import (
	"fmt"
	"log"
	"net/http"

	auth "github.com/fibrasek/2fa_go/controllers"
	"github.com/fibrasek/2fa_go/models"
	"github.com/fibrasek/2fa_go/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	server *gin.Engine

	AuthController auth.AuthController
	AuthRoute      routes.AuthRoute
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("XX Failed to load .env")
	}

	DB, err := gorm.Open(sqlite.Open("2fa_go.db"), &gorm.Config{})
	DB.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatal("XX Failed to connect to the database")
	}
	fmt.Println(">> Database OK!")

	AuthRoute = routes.NewAuthRoute(auth.NewAuthController(DB))

	server = gin.Default()
}

func main() {
	corsConfig := cors.DefaultConfig()
	// Change to the actual front-end URL
	corsConfig.AllowOrigins = []string{"http://localhost:8080"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/health", func(ctx *gin.Context) {
		message := "Hey, I'm healthy :)"
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": message})
	})

	AuthRoute.AuthRouter(router)
	log.Fatal(server.Run(":3000"))
}
