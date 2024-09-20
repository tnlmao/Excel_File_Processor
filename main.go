package main

import (
	apihandler "go_assignment/api_handler"
	"go_assignment/app"
	"go_assignment/logger"
	"go_assignment/utils"

	"github.com/spf13/viper"
)

func init() {
	logger.SetLogger()
	viper.SetConfigFile(utils.Env)
	viper.ReadInConfig()
	viper.AutomaticEnv()
}
func main() {
	application := app.New()
	apihandler.SetupRoutes(application.Router)
	application.Router.Run(":8080")
}
