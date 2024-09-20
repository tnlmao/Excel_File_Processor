package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"go_assignment/config/database"
	r "go_assignment/config/redis"
	"go_assignment/logger"
	"go_assignment/models"
	"go_assignment/service"
	"go_assignment/service/driver"
	"go_assignment/utils"
	"path/filepath"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-redis/redis/v8"
)

type ExcelService struct{}

func NewExcelService() driver.ExcelProcessing {
	return &ExcelService{}
}

func (e *ExcelService) UploadFile(c context.Context) models.Response {
	filePath := utils.FilePath
	ext := filepath.Ext(filePath)

	if ext == ".xlsx" || ext == ".xls" {
		logger.I("File  has a valid Excel extension: ", filePath, ext)
	} else {
		return service.ExcelServiceResponse(500, "Unsupported file format, Please Upload a valid File")
	}
	file, err := excelize.OpenFile("")
	if err != nil {
		logger.E(err)
	}

	rows := file.GetRows(file.GetSheetName(1))
	if err != nil {
		logger.E(err)
		return service.ExcelServiceResponse(500, "Couldnt fetch Rows from the File")
	}
	if len(rows) == 0 {
		return service.ExcelServiceResponse(200, "File Empty")
	}

	ColumnNames := rows[0]
	tableName := CreateTableName(file.GetSheetName(1))
	query := CreateTableQuery(ColumnNames, tableName)

	DB := database.DB
	defer DB.Close()
	_, err = DB.Exec(query)
	if err != nil {
		return service.ExcelServiceResponse(500, err.Error())
	}
	RedisClient := r.Client
	count := RedisClient.DBSize(c).Val()
	if count == 0 {
		AddRecordsToRedis(c, RedisClient, len(rows), rows)
	}
	Batching(c, DB, rows[1:], tableName)
	return models.Response{
		Code: 200,
		Msg:  "Success",
	}
}

func (e *ExcelService) ViewData(c context.Context) models.Response {
	RedisClient := r.Client

	var results []models.OrderedRecord

	size, _ := RedisClient.DBSize(c).Result()
	logger.I("Size --------", size)
	if size > 0 {
		logger.I("Inside Redis")
		for i := 1; i <= int(size); i++ {
			key := fmt.Sprintf("%d", i)
			val, err := RedisClient.Get(c, key).Bytes()
			if err != nil {
				if err == redis.Nil {
					logger.E("Key does not exist", key)
					continue
				}
				logger.E("Error fetching key ", key, err)
			}

			var data models.Record
			if err := json.Unmarshal(val, &data); err != nil {
				logger.E("Error unmarshalling data for key ", key, err)
				continue
			}
			results = append(results, models.OrderedRecord{Key: key, Data: data})
		}
		return models.Response{
			Code:     200,
			Msg:      "success",
			Response: results,
		}
	} else {
		logger.I("Inside DB")
		DB := database.DB
		defer DB.Close()
		var results []models.Record
		query := "SELECT address, city, company_name, county, email, first_name, last_name, phone, postal, web FROM uk_500"
		rows, err := DB.QueryContext(c, query)
		if err != nil {
			logger.E("Error querying database:", err)
			return models.Response{
				Code:     500,
				Msg:      "Internal Server Error",
				Response: nil,
			}
		}
		defer rows.Close()

		for rows.Next() {
			var data models.Record
			if err := rows.Scan(&data.Address, &data.City, &data.CompanyName, &data.County, &data.Email, &data.FirstName, &data.LastName, &data.Phone, &data.Postal, &data.Web); err != nil {
				logger.E("Error scanning row:", err)
				continue
			}
			results = append(results, data)
		}

		if err := rows.Err(); err != nil {
			logger.E("Error during row iteration:", err)
			return models.Response{
				Code:     500,
				Msg:      "Internal Server Error",
				Response: nil,
			}
		}

		return models.Response{
			Code:     200,
			Msg:      "success",
			Response: results,
		}
	}
}
func (e *ExcelService) EditData(c context.Context, editRequest models.EditRequest) models.Response {
	DB := database.DB
	id := editRequest.Id
	currentData, err := GetCurrentData(DB, c, *id)
	if err != nil {
		return models.Response{
			Code: 500,
			Msg:  "Couldnt get current data",
		}
	}
	if editRequest.Address != nil {
		currentData.Address = editRequest.Address
	}
	if editRequest.City != nil {
		currentData.City = editRequest.City
	}
	if editRequest.CompanyName != nil {
		currentData.CompanyName = editRequest.CompanyName
	}
	if editRequest.County != nil {
		currentData.County = editRequest.County
	}
	if editRequest.Email != nil {
		currentData.Email = editRequest.Email
	}
	if editRequest.FirstName != nil {
		currentData.FirstName = editRequest.FirstName
	}
	if editRequest.LastName != nil {
		currentData.LastName = editRequest.LastName
	}
	if editRequest.Phone != nil {
		currentData.Phone = editRequest.Phone
	}
	if editRequest.Postal != nil {
		currentData.Postal = editRequest.Postal
	}
	if editRequest.Web != nil {
		currentData.Web = editRequest.Web
	}

	err = SaveUpdatedData(DB, c, currentData, *id)
	if err != nil {
		return models.Response{
			Code: 500,
			Msg:  "Couldnt save the data",
		}
	}
	redisUpdate := models.Record{
		FirstName:   *currentData.FirstName,
		Address:     *currentData.Address,
		City:        *currentData.City,
		CompanyName: *currentData.CompanyName,
		County:      *currentData.County,
		Email:       *currentData.Email,
		LastName:    *currentData.LastName,
		Phone:       *currentData.Phone,
		Postal:      *currentData.Postal,
		Web:         *currentData.Web,
	}
	mapUpdate, err := json.Marshal(redisUpdate)
	if err != nil {
		logger.E("Error while marshalling ")
	}
	size := r.Client.DBSize(c).Val()
	if size > 0 {
		key := fmt.Sprintf("%d", *id)
		err := r.Client.Set(c, key, mapUpdate, time.Minute*5).Err()
		if err != nil {
			logger.E("Error while updating in redis")
		}
	}
	return models.Response{
		Code: 200,
		Msg:  "success",
	}
}
