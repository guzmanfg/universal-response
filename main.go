package main

import (
	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/home"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.Any("/", home.Home)

	router.Run(":" + port)
}
