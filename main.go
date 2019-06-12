package main

import (
	"./module/database"
	"./routers"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"log"
	"net/url"
	"os"
	"time"
)

func main() {
	os.Mkdir("data", 0777)
	os.Mkdir("data/account_icon", 0777)

	var (
		db    *sql.DB
		dbErr error
	)
	for i := 0; i < 30; i++ {
		connectText := fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true&loc=%s",
			MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_DB_NAME, url.QueryEscape("Asia/Tokyo"))
		_db, dbErr := sql.Open("mysql", connectText)
		err := _db.Ping()
		if dbErr != nil || err != nil {
			if dbErr == nil && err != nil {
				dbErr = err
			}
			log.Printf("DB ERROR %+v \n", dbErr)
			time.Sleep(1 * time.Second)
		} else {
			db = _db
			break
		}
	}
	if db == nil {
		panic(dbErr.Error())
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

	database.Migrate(db, m)

	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			panic(err)
		}
	}
	log.Println("migration complete.")
}
