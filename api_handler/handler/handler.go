package handler

import (
	"context"
	"encoding/json"
	"go_assignment/logger"
	"go_assignment/models"
	"go_assignment/service/driver"
	"go_assignment/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadFile(svc driver.ExcelProcessing) gin.HandlerFunc {

	return func(c *gin.Context) {
		file, err := c.FormFile(utils.File)
		if err != nil {
			logger.E(err)
			c.JSON(http.StatusInternalServerError, "request Error")
		}
		filename := file.Filename
		err = c.SaveUploadedFile(file, utils.SaveFilePath+filename)
		if err != nil {
			logger.E(err)
			c.JSON(http.StatusInternalServerError, "Couldn't Save The File")
		}
		ctx := context.Context(c)
		response := svc.UploadFile(ctx)
		c.JSON(http.StatusOK, response)
	}
}
func ViewData(svc driver.ExcelProcessing) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Context(c)
		response := svc.ViewData(ctx)
		c.JSON(http.StatusOK, response)
	}
}
func EditFile(svc driver.ExcelProcessing) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Context(c)
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "request Error")
		}
		logger.I("REQUEST RECEIVED", string(body))

		editRequest := models.EditRequest{}
		err = json.Unmarshal(body, &editRequest)
		if err != nil {
			logger.E(err)
			c.JSON(http.StatusInternalServerError, "Invalid Request")
			return
		}
		response := svc.EditData(ctx, editRequest)
		c.JSON(http.StatusOK, response)
	}
}
