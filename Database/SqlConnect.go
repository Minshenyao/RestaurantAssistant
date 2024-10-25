package Database

import (
	"RestaurantAssistant/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB
var err error

func ConnectDB() (*sql.DB, error) {
	config := utils.AppConfig.Database
	sqlConnect := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password, config.Host, config.Port, config.Dbname)
	Db, err = sql.Open("mysql", sqlConnect)
	if err != nil {
		return nil, err
	}
	if err = Db.Ping(); err != nil {
		return nil, err
	}
	return Db, nil
}
