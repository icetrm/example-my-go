package main

import (
	"my-go/controllers"
	"my-go/initializers"
	middleware "my-go/middleware"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	//gin.DisableConsoleColor()

	// r := gin.Default()
	gin.SetMode(os.Getenv("MODE"))

	r := gin.New()

	r.SetTrustedProxies([]string{"localhost/32"})

	s := &http.Server{
		Addr:           os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

	r.Use(middleware.ActionLog())
	r.Use(gin.Recovery())

	authennGroup := r.Group("/auth")

	authennGroup.POST("/signup", controllers.Signup)
	authennGroup.POST("/login", controllers.Login)
	// here RequireAuth is a middleware that we will be creating below. It protects the route
	authennGroup.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run()
}
