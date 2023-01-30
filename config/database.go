package config

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"

	"gorm.io/gorm"
)

func DBconn() (db *sql.DB, err error) {

	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "12345678"
	dbName := "crud_terbaru"

	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	return
}

var DB *gorm.DB

// func ConnectDatabase() {
// 	database, err := gorm.Open(mysql.Open("root:12345678@tcp(localhost:3306)/crud_terbaru"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	database.AutoMigrate(&entities.User{})

// 	DB = database
// }
