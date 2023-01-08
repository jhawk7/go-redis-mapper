package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jhawk7/go-redis-mapper/pkg/redis_client"
)

var redisClient *redis_client.RedisClient

// custom middleware
func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := os.Getenv("SECRET_HEADER")
		authSecret := os.Getenv("AUTH_SECRET")

		if reqSecret := c.Request.Header.Get(authHeader); reqSecret == authSecret {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, nil)
		}
	}
}

func main() {
	router := gin.Default()
	redisClient = redis_client.InitClient()
	authorized := router.Group("/authorized", Authorization())
	authorized.GET()
	authorized.POST()
	authorized.PATCH()
	authorized.DELETE()
}

//handlers should utilize service that makes call to redis cahce
