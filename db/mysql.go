package db

import (
	"log"

	"hscan/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//Database type
type Database struct {
	*gorm.DB
}

//IsRecordNotFoundError export
func IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}

func NewDB(conf config.MysqlConfig) *Database {
	log.Printf("mysql.go init")

	DB, err := gorm.Open("mysql", conf.MysqlRes)
	if err != nil {
		log.Printf("open database failed %s", err)
		panic(err)
	}
	log.Printf("connect to mysql success")

	return &Database{DB: DB}
}
