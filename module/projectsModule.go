package module

import (
	"database/sql"
	"log"
)

type ProjectsModule struct {
	db *sql.DB
}

func NewProjectsModule(db *sql.DB) *ProjectsModule {
	return &ProjectsModule{db}
}

func (self *ProjectsModule) Add(name, description string) (error bool, uuid int) {
	stmt, err := self.db.Prepare("INSERT INTO `projects` (`name`, `description`) VALUES(?, ?);")

	if err != nil{
		log.Printf("ProjectsModule.Add Error: %+v", err)
		return true, -1
	}

	defer stmt.Close()

	ins, err := stmt.Exec(name, description)

	if err != nil{
		log.Printf("ProjectsModule.Add Error: %+v", err)
		return true, -1
	}

	insId, err := ins.LastInsertId()

	if err != nil{
		log.Printf("ProjectsModule.Add Error: %+v", err)
		return true, -1
	}

	return false , int(insId)
}
