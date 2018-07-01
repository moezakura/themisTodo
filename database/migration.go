package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const usage string = "Usage: migration [COMMAND] [OPTION...]\n\n" +
	"  init                  initialize database\n" +
	"  up                    upgrade database\n" +
	"  down                  downgrade database\n" +
	"  status                display migration status\n" +
	"  help | -h | --help    display this help and exit"

const migrationTable string = "CREATE TABLE IF NOT EXISTS `_migration` (" +
	"`grade` int NOT NULL," +
	"`name` varchar(32) NOT NULL," +
	"`modified_at` bigint(20) NOT NULL" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;"

var (
	mysql_username string
	mysql_password string
	mysql_db_name  string
)

type Generation struct {
	Grade    int
	UpFile   string
	DownFile string
	Name     string
}

var generationZero = Generation{0, "", "", "zero"}

type AppDB struct {
	Generations []Generation
	Current     Generation
	UserName    string
	Password    string
	Name        string
	db          *sql.DB
}

// マイグレーションファイルの命名規則
// <世代: Int>-<名前: String>.[up|down].sql
var upFileRe = regexp.MustCompile(`/([0-9]{3})-.*?\.up\.sql$`)
var downFileRe = regexp.MustCompile(`/([0-9]{3})-.*?\.down\.sql$`)

func getGenerationFiles(absDir string, re *regexp.Regexp) map[int]string {
	files := make(map[int]string)
	sqlFiles, err := filepath.Glob(filepath.Join(absDir, "*.sql"))
	if err != nil {
		return files
	}
	for _, file := range sqlFiles {
		tmp := re.FindStringSubmatch(file)
		if len(tmp) >= 2 {
			grade, _ := strconv.Atoi(tmp[1])
			files[grade] = file
		}
	}
	return files
}

var nameRe = regexp.MustCompile(`/[0-9]{3}-(.*?)\.up\.sql$`)

func getGenerationName(generationFileName string) (string, error) {
	tmp := nameRe.FindStringSubmatch(generationFileName)
	if len(tmp) >= 2 {
		return tmp[1], nil
	}
	return "", errors.New("Name not found")
}

func getGenerations(dir string) ([]Generation, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	upFiles := getGenerationFiles(absDir, upFileRe)
	downFiles := getGenerationFiles(absDir, downFileRe)
	generations := []Generation{generationZero}
	for grade := 1; grade <= len(upFiles); grade++ {
		upFile, upFileExists := upFiles[grade]
		downFile, downFileExists := downFiles[grade]
		if !upFileExists || !downFileExists {
			break
		}
		name, _ := getGenerationName(upFile)
		generations = append(generations, Generation{grade, upFile, downFile, name})
	}
	return generations, nil
}

func getCurrentGeneration(db *sql.DB) (Generation, error) {
	query := "SELECT grade FROM _migration"
	rows, err := db.Query(query)
	if err != nil {
		return Generation{}, err
	}
	defer rows.Close()
	rows.Next()
	var grade int
	if err := rows.Scan(&grade); err != nil {
		return Generation{}, err
	}
	generations, err := getGenerations("./")
	if err != nil || len(generations) < (grade+1) {
		return Generation{}, errors.New("Generation files are not found or less than current grade")
	}
	return generations[grade], nil
}

func NewNonInitializedAppDB(dir string, db *sql.DB) (*AppDB, error) {
	generations, err := getGenerations(dir)
	if err != nil {
		return nil, errors.New("Can not get the migration files")
	}
	return &AppDB{
		Generations: generations,
		Current:     generationZero,
		UserName:    mysql_username,
		Password:    mysql_password,
		Name:        mysql_db_name,
		db:          db,
	}, nil
}

func NewAppDB(dir string, db *sql.DB) (*AppDB, error) {
	generations, err := getGenerations(dir)
	if err != nil {
		return nil, errors.New("Can not get the migration files")
	}
	current, err := getCurrentGeneration(db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can not get the current generation: %v", err))
	}
	return &AppDB{
		Generations: generations,
		Current:     current,
		UserName:    mysql_username,
		Password:    mysql_password,
		Name:        mysql_db_name,
		db:          db,
	}, nil
}

func (app *AppDB) isUpgradable() bool {
	return (app.Current.Grade + 1) < len(app.Generations)
}

func (app *AppDB) isDowngradable() bool {
	return app.Current != generationZero
}

func (app *AppDB) execSql(sqlFile string) error {
	var cmd string
	if app.Password == "" {
		cmd = fmt.Sprintf("mysql -u%s %s < %s",
			app.UserName, app.Name, sqlFile)
	} else {
		cmd = fmt.Sprintf("mysql -u%s -p%s %s < %s",
			app.UserName, app.Password, app.Name, sqlFile)
	}
	return exec.Command("bash", "-c", cmd).Run()
}

// マイグレーションを記録するDBの初期化
func (app *AppDB) initMigrationDatabase() error {
	if _, err := app.db.Exec(migrationTable); err != nil {
		return err
	}
	updatedAt := time.Now().Unix()
	query := "INSERT INTO `_migration` (`grade`, `name`, `modified_at`) VALUES (0, '', ?);"
	if _, err := app.db.Exec(query, updatedAt); err != nil {
		return err
	}
	return nil
}

func (app *AppDB) Init() error {
	if len(app.Generations) < 2 {
		return errors.New("Could not find migration file")
	}
	if err := app.initMigrationDatabase(); err != nil {
		return errors.New(fmt.Sprintf("fialued to create migration table: %v", err))
	}
	return nil
}

func (app *AppDB) nextGen() Generation {
	if app.isUpgradable() {
		return app.Generations[app.Current.Grade+1]
	}
	return Generation{}
}

func (app *AppDB) Up() error {
	if !app.isUpgradable() {
		return errors.New("Aleady up to date.")
	}
	next := app.nextGen()
	if err := app.execSql(next.UpFile); err != nil {
		return err
	}
	app.db.Exec("UPDATE _migration SET grade=?, name=?, modified_at=?",
		next.Grade, next.Name, time.Now().Unix())
	return nil
}

func (app *AppDB) prevGen() Generation {
	if app.isDowngradable() {
		return app.Generations[app.Current.Grade-1]
	}
	return Generation{}
}

func (app *AppDB) Down() error {
	return errors.New("Not implementation error") // TODO
}

func (app *AppDB) Status() {
	status := "grade | status\n------+-----------\n"
	for _, generation := range app.Generations[1:] {
		point := ""
		if generation.Grade == app.Current.Grade {
			point = "<- current"
		}
		status += fmt.Sprintf("%03d   | %s\n",
			generation.Grade, point)
	}
	fmt.Println(status)
}

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

func getDbConn() (*sql.DB, error) {
	connectText := fmt.Sprintf("%s:%s@/%s", mysql_username, mysql_password, mysql_db_name)
	return sql.Open("mysql", connectText)
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

	command := os.Args[1]
	switch command {
	case "help", "--help", "-h":
		fmt.Println(usage)
		os.Exit(0)
	case "init":
		db, err := getDbConn()
		if err != nil {
			panic(err)
		}
		defer db.Close()
		app, err := NewNonInitializedAppDB("./", db)
		if err != nil {
			panic(err)
		}
		if err := app.Init(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	db, err := getDbConn()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app, err := NewAppDB("./", db)
	if err != nil {
		fmt.Println("Database is not initialzed. Please execute `migration init` on your shell.")
		os.Exit(1)
	}

	switch command {
	case "up":
		if err := app.Up(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "down":
		if err := app.Down(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "status":
		app.Status()
	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}
