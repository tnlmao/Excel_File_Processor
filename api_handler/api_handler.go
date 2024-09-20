package apihandler

import (
	"go_assignment/api_handler/handler"
	"go_assignment/middleware"
	"go_assignment/service/logic"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.Use(middleware.DbMiddleware)
	router.Use(middleware.RedisMiddleware)

	svc := logic.NewExcelService()

	uploadFile := handler.UploadFile(svc)
	viewFile := handler.ViewData(svc)
	editFile := handler.EditFile(svc)

	router.POST("/upload", uploadFile)
	router.GET("/view", viewFile)
	router.PATCH("/edit", editFile)

}
