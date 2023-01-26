package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jhawk7/go-redis-mapper/pkg/redis_client"
	log "github.com/sirupsen/logrus"
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
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func main() {
	router := gin.Default()
	redisClient = redis_client.InitClient()
	authorized := router.Group("/authorized").Use(Authorization())
	authorized.GET("/mapper/:key", GetValueByKey)
	authorized.POST("/mapper", StoreKVPair)
	authorized.PATCH("/mapper", UpdateKVPair)
	authorized.DELETE("/mapper", DeleteKeys)
	router.Run() //running on port 8080
}

func GetValueByKey(c *gin.Context) {
	key := c.Param("key")
	value, err := redisClient.GetValue(c.Request.Context(), key)
	if err != nil {
		ErrorHandler(c, err, http.StatusBadRequest, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"value": value,
	})
}

func StoreKVPair(c *gin.Context) {
	var kv redis_client.KVPair
	if bindErr := c.BindJSON(&kv); bindErr != nil {
		err := fmt.Errorf("failed to bind request data to object; [error: %v]", bindErr.Error())
		ErrorHandler(c, err, 0, false)
		return
	}

	if storeErr := redisClient.Store(c.Request.Context(), kv); storeErr != nil {
		ErrorHandler(c, storeErr, http.StatusBadRequest, false)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func UpdateKVPair(c *gin.Context) {
	var kv redis_client.KVPair
	if bindErr := c.BindJSON(&kv); bindErr != nil {
		err := fmt.Errorf("failed to bind request data to object; [error: %v]", bindErr.Error())
		ErrorHandler(c, err, 0, false)
		return
	}

	if updateErr := redisClient.UpdateValue(c.Request.Context(), kv); updateErr != nil {
		ErrorHandler(c, updateErr, http.StatusBadRequest, false)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func DeleteKeys(c *gin.Context) {
	var keys redis_client.DeleteKeys
	if bindErr := c.BindJSON(&keys); bindErr != nil {
		err := fmt.Errorf("failed to bind request data to object; [error: %v]", bindErr.Error())
		ErrorHandler(c, err, 0, false)
		return
	}

	if deleteErr := redisClient.Delete(c.Request.Context(), keys); deleteErr != nil {
		ErrorHandler(c, deleteErr, http.StatusBadRequest, false)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func ErrorHandler(c *gin.Context, err error, status int, fatal bool) {
	if err != nil {
		log.Error(fmt.Errorf("error: %v", err.Error()))

		if fatal {
			panic(err)
		}

		if status != 0 {
			c.JSON(status, gin.H{
				"error": err.Error(),
			})
		}
	}
}
