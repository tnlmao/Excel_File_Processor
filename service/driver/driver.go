package driver

import (
	"context"
	"go_assignment/models"
)

type ExcelProcessing interface {
	UploadFile(c context.Context) models.Response
	ViewData(c context.Context) models.Response
	EditData(c context.Context, editRequest models.EditRequest) models.Response
}
