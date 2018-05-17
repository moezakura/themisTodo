package main

import (
	"./routers"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"
)

func main() {
	os.Mkdir("www/assets/accountIcon", 0777)

	connectText := fmt.Sprintf("%s:%s@/%s", MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_DB_NAME)
	db, err := sql.Open("mysql", connectText)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	r := routers.Init(db)
	r.Run(":31204")
}
