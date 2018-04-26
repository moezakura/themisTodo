package main

import (
	"./routers"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

func main() {
	connectText := fmt.Sprintf("%s:%s@/%s", MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_DB_NAME)
	db, err := sql.Open("mysql", connectText)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	r := routers.Init(db)
	r.Run(":31204")
}
