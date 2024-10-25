package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
type BillAmounts struct {
	Name   string  `form:"name"`
	Amount float64 `json:"amount"`
}
type Totals struct {
	BillAmounts []BillAmounts
	Total       float64 `json:"total"`
}
type User struct {
	Phone   string  `json:"phone"`
	Vip     bool    `json:"vip"`
	Deposit float64 `json:"deposit"`
}

// AddSpecial 添加一条特价菜品
func AddSpecial(insertData SpecialDish, db *sql.DB) (string, error) {
	today := time.Now().Format("2006-01-02")
	_, err := db.Exec("insert into Special_dishes (name, quantily, price, date) values (?,?,?,?)", insertData.Name, insertData.Number, insertData.Price, today)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

// DeleteSpecialOne 删除指定名称的特价菜品
func DeleteSpecialOne(deleteData SpecialName, db *sql.DB) (string, error) {
	_, err := db.Exec("delete from Special_dishes where name = ?", deleteData.Name)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "ok", nil
}

// DeleteSpecialAll 删除某天所有特价菜品
func DeleteSpecialAll(time SpecialTime, db *sql.DB) (string, error) {
	_, err := db.Exec("delete from Special_dishes where date = ?", time.Time)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "ok", nil
}

// GetTodaySpecials 获取今日特价菜品
func GetTodaySpecials(db *sql.DB) ([]SpecialDish, error) {
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

// BillAmount 账单金额结算
func BillAmount(db *sql.DB) (Totals, error) {
	var billAmounts []BillAmounts
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
		var billAmount BillAmounts
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

// MemberRecharge 会员充值系统
func MemberRecharge(phoneNumber string, money float64, db *sql.DB) error {
	_ = Register(phoneNumber, db)
	userInfo, err := db.Query("select * from User where phone = ?", phoneNumber)
	if err != nil {
		log.Fatal(err)
	}
	defer userInfo.Close()
	if money > 0.00 && money > 100.00 {
		err := UpdateUserInfo(phoneNumber, money, 0.00, db)
		if err != nil {
			log.Fatal(err)
		}
	} else if money >= 100.00 && money < 200.00 {
		err := UpdateUserInfo(phoneNumber, money, 10.00, db)
		if err != nil {
			log.Fatal(err)
		}
	} else if money >= 200 && money < 500 {
		err := UpdateUserInfo(phoneNumber, money, 30.00, db)
		if err != nil {
			log.Fatal(err)
		}
	} else if money >= 500 && money < 1000 {
		err := UpdateUserInfo(phoneNumber, money, 80.00, db)
		if err != nil {
			log.Fatal(err)
		}
	} else if money >= 1000 {
		err := UpdateUserInfo(phoneNumber, money, 200.00, db)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// UpdateUserInfo 更新用户充值信息
func UpdateUserInfo(phoneNumber string, money float64, addition float64, db *sql.DB) error {
	var vipStatus bool
	var rawMoney float64
	_ = Register(phoneNumber, db)
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

// Register 注册用户
func Register(phoneNumber string, db *sql.DB) error {
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
