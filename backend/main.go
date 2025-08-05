package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  
  // Configure CORS
  config := cors.DefaultConfig()
  config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173", "http://127.0.0.1:3000"}
  config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
  config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
  
  r.Use(cors.New(config))
  
  r.GET("/", func(c *gin.Context) {
    c.String(200, "Hello, Gin!")
  })
  r.Run(":8080") // listens on 0.0.0.0:8080 by default
}
