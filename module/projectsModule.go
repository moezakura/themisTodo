package module

import (
	"database/sql"
	"log"
	"../models"
)

type ProjectsModule struct {
	db *sql.DB
}

func NewProjectsModule(db *sql.DB) *ProjectsModule {
	return &ProjectsModule{db}
}

func (self *ProjectsModule) Add(name, description string) (error bool, uuid int) {
	stmt, err := self.db.Prepare("INSERT INTO `projects` (`name`, `description`) VALUES(?, ?);")

	if err != nil {
		log.Printf("ProjectsModule.Add Error: %+v", err)
		return true, -1
	}

	defer stmt.Close()

	ins, err := stmt.Exec(name, description)

	if err != nil {
		log.Printf("ProjectsModule.Add Error: %+v", err)
		return true, -1
	}

	insId, err := ins.LastInsertId()

	if err != nil {
		log.Printf("ProjectsModule.Add Error: %+v", err)
		return true, -1
	}

	return false, int(insId)
}

func (self *ProjectsModule) AddUser(userId, projectId int) bool {
	stmt, err := self.db.Prepare("INSERT INTO `users_in_projects` (`user_id`, `project_id`, `enable`, `expiration`) VALUES(?, ?, TRUE, NULL);")

	if err != nil {
		log.Printf("ProjectsModule.AddUser Error: %+v", err)
		return true
	}

	defer stmt.Close()

	_, err = stmt.Exec(userId, projectId)
	if err != nil {
		log.Printf("ProjectsModule.AddUser Error: %+v", err)
		return true
	}

	return false
}

func (self *ProjectsModule) GetProjects(userId int) (error bool, project []models.Project) {
	project = []models.Project{}

	rows, err := self.db.Query("SELECT `uuid`,`name`,`description` FROM `projects` WHERE `projects`.`uuid` IN (SELECT `project_id` FROM `users_in_projects` WHERE `user_id` = ? ORDER BY `user_id`);", userId)

	if err != nil {
		return true, nil
	}

	for rows.Next() {
		projectOne := models.Project{}
		if err := rows.Scan(&projectOne.Uuid, &projectOne.Name, &projectOne.Description); err != nil {
			log.Printf("ProjectsModule.GetProject Error: %+v\n", err)
			return true, nil
		}
		project = append(project, projectOne)
	}

	return false, project
}

func (self *ProjectsModule) GetProject(userId int) (error bool, project *models.Project) {
	project = &models.Project{}

	rows, err := self.db.Query("SELECT `uuid`,`name`,`description` FROM `projects` WHERE `uuid` = ?;", userId)

	if err != nil {
		return true, nil
	}

	if !rows.Next() {
		return true, nil
	}

	if err := rows.Scan(&project.Uuid, &project.Name, &project.Description); err != nil {
		log.Printf("ProjectsModule.GetProject Error: %+v\n", err)
		return true, nil
	}

	return false, project
}
