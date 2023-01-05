package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

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
	authorized := router.Group("/authorized", Authorization())
	authorized.GET()
	authorized.POST()
	authorized.PATCH()
	authorized.DELETE()
}

//handlers should utilize service that makes call to redis cahce
