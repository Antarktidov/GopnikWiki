package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/api/articles", getArticles)

	router.Run("localhost:8080")
}

func getArticles(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, nil)
}
