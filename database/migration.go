package main

import (
	"errors"
	"fmt"
	"os/exec"
	"os"
)

const usage string =
 	"Usage: migration [COMMAND] [OPTION...]\n\n" +
	"  init [FILE]           initialize database by [FILE]\n" +
	"  help | -h | --help    display this help and exit"

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
	case "-h", "--help", "help":
		fmt.Println(usage)
	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}
