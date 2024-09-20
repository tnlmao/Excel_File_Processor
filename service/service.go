package service

import "go_assignment/models"

func ExcelServiceResponse(code int, Msg string) models.Response {
	return models.Response{
		Code: code,
		Msg:  Msg,
	}
}
