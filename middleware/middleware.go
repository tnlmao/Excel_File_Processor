package middleware

import (
	"go_assignment/config/database"
	"go_assignment/config/redis"
	"go_assignment/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DbMiddleware(c *gin.Context) {
	err := database.ConnectToMysql()
	if err != nil {
		logger.E(err)
		c.JSON(http.StatusInternalServerError, "Database Connection Error")
	}
	c.Next()
}
func RedisMiddleware(c *gin.Context) {
	err := redis.ConnectToRedis()
	if err != nil {
		logger.E(err)
		c.JSON(http.StatusInternalServerError, "Redis Connection Error")
	}
	c.Next()
}
