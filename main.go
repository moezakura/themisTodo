package main

import (
	"./routers"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"os"
	"time"
)

func main() {
	os.Mkdir("data", 0777)
	os.Mkdir("data/account_icon", 0777)

	connectText := fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true",
		MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_DB_NAME)
	db, err := sql.Open("mysql", connectText)
	if err != nil {
		panic(err.Error())
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(15000)
	db.SetMaxIdleConns(3000)

	schema_rows, err := db.Query("SHOW TABLES LIKE 'schema_migrations';")
	if err != nil {
		panic(err)
	}
	if !schema_rows.Next() {
		db.Exec("CREATE TABLE `schema_migrations` (`version` bigint(20) NOT NULL,`dirty` tinyint(1) NOT NULL) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;")
		db.Exec("INSERT INTO `schema_migrations` (`version`, `dirty`) VALUES (2, 0);")
		db.Exec("ALTER TABLE `schema_migrations` ADD PRIMARY KEY (`version`); COMMIT;")
	}

	db_migrate(db)

	defer db.Close()

	r := routers.Init(db)
	r.Run(":31204")
}

func db_migrate(db *sql.DB) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./database",
		"mysql",
		driver,
	)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			panic(err)
		}
	}
}
