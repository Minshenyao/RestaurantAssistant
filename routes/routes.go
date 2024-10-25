package routes

import (
	"RestaurantAssistant/Database"
	"RestaurantAssistant/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoutes(router *gin.Engine) {
	router.GET("/todaySpecials", func(context *gin.Context) {
		specials, err := services.GetTodaySpecials(Database.Db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		context.JSON(http.StatusOK, gin.H{"data": specials})
	})
	router.POST("/addSpecial", func(context *gin.Context) {
		var insertData services.SpecialDish
		if err := context.BindJSON(&insertData); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		status, err := services.AddSpecial(insertData, Database.Db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		context.JSON(http.StatusOK, gin.H{"status": status})
	})
	router.DELETE("/deleteSpecialOne", func(context *gin.Context) {
		var deleteName services.SpecialName
		if err := context.BindJSON(&deleteName); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		status, err := services.DeleteSpecialOne(deleteName, Database.Db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		context.JSON(http.StatusOK, gin.H{"status": status})
	})
	router.DELETE("/deleteSpecialAll", func(context *gin.Context) {
		var targetTime services.SpecialTime
		if err := context.BindQuery(&targetTime); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		status, err := services.DeleteSpecialAll(targetTime, Database.Db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		context.JSON(http.StatusOK, gin.H{"status": status})
	})
	router.GET("/BillAmount", func(context *gin.Context) {
		var targetTime services.SpecialTime
		if err := context.BindQuery(&targetTime); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		status, err := services.BillAmount(Database.Db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		context.JSON(http.StatusOK, gin.H{"status": status})
	})
	router.POST("/memberRecharge", func(context *gin.Context) {
		var userInfo services.User
		if err := context.BindJSON(&userInfo); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		err := services.MemberRecharge(userInfo.Phone, userInfo.Deposit, Database.Db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		context.JSON(http.StatusOK, gin.H{"status": userInfo.Phone + "充值成功!"})
	})
}
