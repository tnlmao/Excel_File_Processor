package logic

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"go_assignment/logger"
	"go_assignment/models"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Record struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	County      string `json:"county"`
	Postal      string `json:"postal"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Web         string `json:"web"`
}

var wg sync.WaitGroup

func CreateTableName(fileName string) string {
	cleanFileName := strings.ReplaceAll(fileName, " ", "_")
	cleanFileName = strings.ReplaceAll(cleanFileName, ".", "_")
	cleanFileName = strings.ReplaceAll(cleanFileName, "-", "_")

	return cleanFileName
}

func CreateTableQuery(ColumnNames []string, tableName string) string {
	var columnsWithTypes []string

	columnsWithTypes = append(columnsWithTypes, "id INT AUTO_INCREMENT PRIMARY KEY")

	for _, col := range ColumnNames {
		cleanCol := strings.ReplaceAll(col, " ", "_")
		cleanCol = strings.ReplaceAll(cleanCol, ".", "_")
		columnsWithTypes = append(columnsWithTypes, fmt.Sprintf("%s VARCHAR(255)", cleanCol))
	}

	columnsDefinition := strings.Join(columnsWithTypes, ", ")

	createTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, columnsDefinition)
	return createTableQuery
}
func CreateViewDataQuery(ColumnNames []string, tableName string) string {
	var columnsWithTypes []string

	columnsWithTypes = append(columnsWithTypes, "id INT AUTO_INCREMENT PRIMARY KEY")

	for _, col := range ColumnNames {
		cleanCol := strings.ReplaceAll(col, " ", "_")
		cleanCol = strings.ReplaceAll(cleanCol, ".", "_")
		columnsWithTypes = append(columnsWithTypes, fmt.Sprintf("%s VARCHAR(255)", cleanCol))
	}

	ViewData := strings.Join(columnsWithTypes, ", ")

	createViewDataQuery := fmt.Sprintf("SELECT %s FROM %s ;", ViewData, tableName)
	return createViewDataQuery
}
func Batching(ctx context.Context, DB *sql.DB, rows [][]string, tableName string) error {
	batchSize := 100
	for i := 0; i < len(rows); i += batchSize {
		batch := rows[i:min(i+batchSize, len(rows))]
		wg.Add(1)
		go func(batch [][]string) {
			defer wg.Done()
			err := StoreDetails(ctx, DB, batch, tableName)
			if err != nil {
				logger.E("Error storing details: ", err)
			}
		}(batch)

	}
	wg.Wait()
	return nil
}
func StoreDetails(ctx context.Context, DB *sql.DB, batch [][]string, tableName string) error {
	var err error
	valueStrings := []string{}
	valueArgs := []interface{}{}

	for _, row := range batch {
		if len(row) != 10 {
			return errors.New("row has incorrect number of columns")
		}
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		for _, col := range row {
			valueArgs = append(valueArgs, col)
		}
	}
	query := fmt.Sprintf("INSERT INTO %s (first_name, last_name, company_name, address, city, county, postal, phone, email, web) VALUES %s", tableName, strings.Join(valueStrings, ","))

	_, err = DB.Exec(query, valueArgs...)
	if err != nil {
		return fmt.Errorf("error executing batch insert: %v", err)
	}
	return err

}
func ConvertExcelToCSV(rows [][]string) error {

	file, err := os.Create("/var/lib/mysql-files/csvFile.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
func LoadCSVIntoDB(db *sql.DB, tableName string) error {
	query := fmt.Sprintf("LOAD DATA  INFILE '/var/lib/mysql-files/csvFile.csv' INTO TABLE %s FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '\"'  LINES TERMINATED BY '\n' IGNORE 1 LINES (first_name, last_name, company_name, address, city, county, postal, phone, email, web);", tableName)
	logger.I(query)
	_, err := db.Exec(query)
	return err
}
func AddRecordsToRedis(ctx context.Context, client *redis.Client, numRecords int, rows [][]string) {
	for i := 1; i <= 2000; i++ {
		key := fmt.Sprintf("%d", i)
		// redisData := make(map[string]interface{})
		redisData := map[string]interface{}{
			"first_name":   rows[i][0],
			"last_name":    rows[i][1],
			"company_name": rows[i][2],
			"address":      rows[i][3],
			"city":         rows[i][4],
			"county":       rows[i][5],
			"postal":       rows[i][6],
			"phone":        rows[i][7],
			"email":        rows[i][8],
			"web":          rows[i][9],
		}
		recordJSON, err := json.Marshal(redisData)
		if err != nil {
			logger.E(err)
		}

		err = client.Set(ctx, key, recordJSON, 5*time.Minute).Err()
		if err != nil {
			logger.E("could not set record ", err)
		}
	}

	fmt.Printf("Added %d records to Redis\n", numRecords)
}
func GetCurrentData(DB *sql.DB, ctx context.Context, id int) (*models.EditRequest, error) {
	var data models.EditRequest
	query := "SELECT id, address, city, company_name, county, email, first_name, last_name, phone, postal, web FROM uk_500 WHERE id = ?"
	err := DB.QueryRowContext(ctx, query, id).Scan(&data.Id, &data.Address, &data.City, &data.CompanyName, &data.County, &data.Email, &data.FirstName, &data.LastName, &data.Phone, &data.Postal, &data.Web)
	if err != nil {
		logger.E(err)
		return &models.EditRequest{}, err
	}
	return &data, nil
}

func SaveUpdatedData(DB *sql.DB, ctx context.Context, data *models.EditRequest, id int) error {
	query := "UPDATE uk_500 SET address = ?, city = ?, company_name = ?, county = ?, email = ?, first_name = ?, last_name = ?, phone = ?, postal = ?, web = ? WHERE id = ?"
	_, err := DB.ExecContext(ctx, query, data.Address, data.City, data.CompanyName, data.County, data.Email, data.FirstName, data.LastName, data.Phone, data.Postal, data.Web, id)
	return err
}
