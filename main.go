package main

import (
	"fmt"

	"net/http"

	"example.com/incident-api/config"
	"example.com/incident-api/database"
	"example.com/incident-api/routes"
	"go.uber.org/zap"
)

func main() {
	if err := config.InitConfig(); err != nil {
		fmt.Println("Error in initialising config ", err.Error())
		return
	}
	if err := database.InitDatabase(); err != nil {
		config.Logger.Error("Unable to init database",
			zap.String("error", err.Error()),
			zap.String("type", config.ConfigObj.DBConfig.DatabaseType))
		return
	}
	if err := database.DB.Connect(); err != nil {
		config.Logger.Error("Unable to connect to databse",
			zap.String("error", err.Error()),
			zap.String("type", config.ConfigObj.DBConfig.DatabaseType),
			zap.String("connection", config.ConfigObj.DBConfig.ConnectionString))
		return
	}

	routes := routes.SetUpRoutes()
	config.Logger.Info("Server Started", zap.Int("PORT", config.ConfigObj.Port))
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.ConfigObj.Port), routes))
}
