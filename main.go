package main

import (
	"./routers"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"
	"time"
)

func main() {
	os.Mkdir("www/assets/accountIcon", 0777)
	os.Mkdir("data", 0777)

	connectText := fmt.Sprintf("%s:%s@tcp(%s)/%s", MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_DB_NAME)
	db, err := sql.Open("mysql", connectText)
	if err != nil {
		panic(err.Error())
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(15000)
	db.SetMaxIdleConns(3000)
	defer db.Close()

	r := routers.Init(db)
	r.Run(":31204")
}
