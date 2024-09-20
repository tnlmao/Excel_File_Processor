package app

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
}

func New() *App {
	return &App{
		Router: gin.Default(),
	}
}
