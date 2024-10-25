package main

// @title 餐厅助手
// @version 1.0
// @description 一个在校大学生学习java语言时书上的案例，用go的web方法重构一遍，按照课程进度更新

import (
	"RestaurantAssistant/Database"
	"RestaurantAssistant/routes"
	"RestaurantAssistant/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadConfig()
	router := gin.Default()
	_, err := Database.ConnectDB()
	if err != nil {
		return
	}
	routes.InitRoutes(router)

	router.Run(utils.AppConfig.Services.Port)
}
