package main

// @title 餐厅助手
// @version 1.0
// @description 一个在校大学生学习java语言时书上的案例，用go的web方法重构一遍，按照课程进度更新

import (
	"RestaurantAssistant/Database"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type SpecialDish struct {
	Name   string  `json:"name"`
	Number int     `json:"number"`
	Price  float64 `json:"price"`
}

type SpecialTime struct {
	Time string `form:"time"`
}
type SpecialName struct {
	Name string `form:"name"`
}
type BillAmount struct {
	Name   string  `form:"name"`
	Amount float64 `json:"amount"`
}
type Totals struct {
	BillAmounts []BillAmount
	Total       float64 `json:"total"`
}

func main() {
	router := gin.Default()
	db, err := Database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	router.GET("/todaySpecials", func(context *gin.Context) {
		specials, err := getTodaySpecials(db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		context.JSON(http.StatusOK, gin.H{"data": specials})
	})
	router.POST("/addSpecial", func(context *gin.Context) {
		var insertData SpecialDish
		if err := context.BindJSON(&insertData); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		status, err := addSpecial(insertData, db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		context.JSON(http.StatusOK, gin.H{"status": status})
	})
	router.DELETE("/deleteSpecialOne", func(context *gin.Context) {
		var deleteName SpecialName
		if err := context.BindJSON(&deleteName); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		status, err := deleteSpecialOne(deleteName, db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		context.JSON(http.StatusOK, gin.H{"status": status})
	})
	router.DELETE("/deleteSpecialAll", func(context *gin.Context) {
		var targetTime SpecialTime
		if err := context.BindQuery(&targetTime); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		status, err := deleteSpecialAll(targetTime, db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		context.JSON(http.StatusOK, gin.H{"status": status})
	})
	router.GET("/BillAmount", func(context *gin.Context) {
		var targetTime SpecialTime
		if err := context.BindQuery(&targetTime); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		status, err := billAmount(db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		context.JSON(http.StatusOK, gin.H{"status": status})
	})
	router.Run(":8000")
}

// 添加一条特价菜品
func addSpecial(insertData SpecialDish, db *sql.DB) (string, error) {
	today := time.Now().Format("2006-01-02")
	_, err := db.Exec("insert into Special_dishes (name, quantily, price, date) values (?,?,?,?)", insertData.Name, insertData.Number, insertData.Price, today)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

// 删除指定名称的特价菜品
func deleteSpecialOne(deleteData SpecialName, db *sql.DB) (string, error) {
	_, err := db.Exec("delete from Special_dishes where name = ?", deleteData.Name)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "ok", nil
}

// 删除某天所有特价菜品
func deleteSpecialAll(time SpecialTime, db *sql.DB) (string, error) {
	_, err := db.Exec("delete from Special_dishes where date = ?", time.Time)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "ok", nil
}

// 获取今日特价菜品
func getTodaySpecials(db *sql.DB) ([]SpecialDish, error) {
	var specials []SpecialDish
	today := time.Now().Format("2006-01-02")
	rows, err := db.Query("select name, quantily, price from Special_dishes where date = ?", today)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var special SpecialDish
		if err := rows.Scan(&special.Name, &special.Number, &special.Price); err != nil {
			log.Fatal(err)
			return nil, err
		}
		specials = append(specials, special)
	}
	return specials, nil
}

// 账单金额结算
func billAmount(db *sql.DB) (Totals, error) {
	var billAmounts []BillAmount
	var total = 0.0
	rows, err := db.Query("select name, quantily, price from Special_dishes where date = ?", time.Now().Format("2006-01-02"))
	if err != nil {
		log.Fatal(err)
		return Totals{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var special SpecialDish
		if err := rows.Scan(&special.Name, &special.Number, &special.Price); err != nil {
			log.Fatal(err)
			return Totals{}, err
		}
		var billAmount BillAmount
		fmt.Printf("菜名: %s\t\t价格: %f\n", special.Name, special.Price*float64(special.Number))
		billAmount.Name = special.Name
		billAmount.Amount = special.Price * float64(special.Number)
		billAmounts = append(billAmounts, billAmount)
		total += billAmount.Amount
	}

	var totals Totals
	totals.BillAmounts = billAmounts
	totals.Total = total
	return totals, nil
}
