package main

import (
	"errors"
	"fmt"
	"os/exec"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const usage string =
 	"Usage: migration [COMMAND] [OPTION...]\n\n" +
	"  init [FILE]           initialize database by [FILE]\n" +
	"  help | -h | --help    display this help and exit"

const migrationTable string =
	"CREATE TABLE IF NOT EXISTS `_migration` (" +
	"`app` varchar(16) NOT NULL," +
	"`name` varchar(16) NOT NULL," +
	"`modified_at` bigint(20) NOT NULL" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;"

var (
	mysql_username string
	mysql_password string
	mysql_db_name  string
)

func initEnv() error {
	mysql_username = os.Getenv("MYSQL_USERNAME")
	mysql_password = os.Getenv("MYSQL_PASSWORD")
	mysql_db_name = os.Getenv("MYSQL_DB_NAME")
	switch {
	case mysql_username == "":
		return errors.New("Environment variable `MYSQL_USERNAME` not found.")
	case mysql_db_name == "":
		return errors.New("Environment variable `MYSQL_DB_NAME` not found.")
	}
	return nil
}

// データベース初期化
// baseSql: 最初のデータベース定義をするフィアル
func initDatabase(baseSql string) error {
	var mysqlCommand string
	if mysql_password == "" {
		mysqlCommand = fmt.Sprintf("mysql -u%s < %s",
			mysql_username, baseSql)
	} else {
		mysqlCommand = fmt.Sprintf("mysql -u%s -p%s < %s",
			mysql_username, mysql_password, baseSql)

	}
	cmd := exec.Command("bash","-c", mysqlCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func getDbConn() (*sql.DB, error) {
	connectText := fmt.Sprintf("%s:%s@/%s", mysql_username, mysql_password, mysql_db_name)
	return sql.Open("mysql", connectText)
}

// マイグレーションを記録するDBの初期化
func initMigrationDatabase(db *sql.DB) error {
	if _, err := db.Exec(migrationTable); err != nil {
		return err
	}
	updatedAt := time.Now().Unix()
	query := "INSERT INTO `_migration` (`app`, `name`, `modified_at`) VALUES ('themis_todo', 'base', ?);"
	if _, err := db.Exec(query, updatedAt);
		err != nil {
		return err
	}
	return nil
}

func main() {
	if err := initEnv(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	switch command := os.Args[1]; command {
	case "init":
		if len(os.Args) < 3 { fmt.Println(usage); os.Exit(1) }
		baseSql := os.Args[2]
		if _, err := os.Stat(baseSql); err != nil {
			fmt.Printf("FILE: %s is not exist\n", baseSql)
			os.Exit(1)
		}
		if err := initDatabase(baseSql); err != nil {
			panic(err)
		}
		db, err := getDbConn()
		if err != nil { panic(err) }
		defer db.Close()
		if err := initMigrationDatabase(db); err != nil {
			panic(err)
		}
	case "-h", "--help", "help":
		fmt.Println(usage)
	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}
