package main

// @title 餐厅助手
// @version 1.0
// @description 一个在校大学生学习java语言时书上的案例，用go的web方法重构一遍，按照课程进度更新

import (
	"RestaurantAssistant/Database"
	"database/sql"
	"errors"
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
type User struct {
	Phone   string  `json:"phone"`
	Vip     bool    `json:"vip"`
	Deposit float64 `json:"deposit"`
}

func main() {
	router := gin.Default()
	db, err := Database.ConnectDB()
	if err != nil {
		log.Fatal(err)
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
	router.POST("/memberRecharge", func(context *gin.Context) {
		var userInfo User
		if err := context.BindJSON(&userInfo); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		err := memberRecharge(userInfo.Phone, userInfo.Deposit, db)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		context.JSON(http.StatusOK, gin.H{"status": userInfo.Phone + "充值成功!"})
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
	}
	defer rows.Close()
	for rows.Next() {
		var special SpecialDish
		if err := rows.Scan(&special.Name, &special.Number, &special.Price); err != nil {
			log.Fatal(err)
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
	}
	defer rows.Close()
	for rows.Next() {
		var special SpecialDish
		if err := rows.Scan(&special.Name, &special.Number, &special.Price); err != nil {
			log.Fatal(err)
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

// 会员充值系统
func memberRecharge(phoneNumber string, money float64, db *sql.DB) error {
	_ = register(phoneNumber, db)
	userInfo, err := db.Query("select * from User where phone = ?", phoneNumber)
	if err != nil {
		log.Fatal(err)
	}
	defer userInfo.Close()
	if money > 0.00 && money > 100.00 {
		err := updateUserInfo(phoneNumber, money, 0.00, db)
		if err != nil {
			log.Fatal(err)
		}
	} else if money >= 100.00 && money < 200.00 {
		err := updateUserInfo(phoneNumber, money, 10.00, db)
		if err != nil {
			log.Fatal(err)
		}
	} else if money >= 200 && money < 500 {
		err := updateUserInfo(phoneNumber, money, 30.00, db)
		if err != nil {
			log.Fatal(err)
		}
	} else if money >= 500 && money < 1000 {
		err := updateUserInfo(phoneNumber, money, 80.00, db)
		if err != nil {
			log.Fatal(err)
		}
	} else if money >= 1000 {
		err := updateUserInfo(phoneNumber, money, 200.00, db)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// updateUserInfo 更新用户充值信息
func updateUserInfo(phoneNumber string, money float64, addition float64, db *sql.DB) error {
	var vipStatus bool
	var rawMoney float64
	_ = register(phoneNumber, db)
	err := db.QueryRow("SELECT vip, deposit FROM User WHERE phone = ?", phoneNumber).Scan(&vipStatus, &rawMoney)
	if err != nil {
		log.Fatal(err)
	}
	if vipStatus == false && money >= 2000 {
		vipStatus = true
	} else {
		time.Sleep(1 * time.Second) // 开个小玩笑
	}
	if vipStatus == false {
		addition = 0.00
	}
	money += rawMoney
	_, err = db.Exec("update User set deposit=?, vip=? where phone = ? ", money+addition, vipStatus, phoneNumber)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 注册用户
func register(phoneNumber string, db *sql.DB) error {
	var existingPhone string
	err := db.QueryRow("SELECT phone FROM User WHERE phone = ?", phoneNumber).Scan(&existingPhone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := db.Exec("INSERT INTO User (phone) VALUES (?)", phoneNumber)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return errors.New("用户已经存在")
}
