package database

import (
	"database/sql"
	"go_assignment/logger"
	"go_assignment/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var DB *sql.DB

func ConnectToMysql() (err error) {
	DB, err = sql.Open("mysql", viper.GetString(utils.DSN))
	if err != nil {
		logger.E(err)
		return
	}
	if err = DB.Ping(); err != nil {
		logger.E(err)
		return
	}
	logger.I("MySql Instance Started")
	return
}
