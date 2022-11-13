package main

import (
  "github.com/gin-gonic/gin"
)

import (
	"net/http"
)

func main() {
  router := gin.Default()

  api := router.Group("/api")
  {
    api.GET("/test", func(ctx *gin.Context) {
      ctx.JSON(200, gin.H{"msg": "Hellow World"})
    })
  }

  router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

  router.Run(":8080")
}
