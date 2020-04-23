package main

import (
	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/home"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("**/*.tmpl.html")
	router.Static("/static", "static")

	router.Any("/", home.Home)
	router.Any("/password", func(c *gin.Context) {
		c.HTML(http.StatusOK, "password.tmpl.html", nil)
	})
	router.GET("/practice", func(c *gin.Context) {
		c.HTML(http.StatusOK, "practice.tmpl.html", nil)
	})
	router.GET("/files/:name", func(c *gin.Context) {
		log.Printf("Downloading file %s\n", c.Param("name"))
		c.File("static/practice/" + c.Param("name"))
	})

	_ = router.Run(":" + port)
}
