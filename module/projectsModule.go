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

	defer rows.Close()

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

	defer rows.Close()

	if !rows.Next() {
		return true, nil
	}

	if err := rows.Scan(&project.Uuid, &project.Name, &project.Description); err != nil {
		log.Printf("ProjectsModule.GetProject Error: %+v\n", err)
		return true, nil
	}

	return false, project
}

func (self *ProjectsModule) GetUser(projectId int) (error bool, accounts []models.Account) {
	accounts = make([]models.Account, 0)

	rows, err := self.db.Query("SELECT `uuid`, `name`, `displayName` FROM `users` WHERE `users`.`uuid` IN (SELECT `user_id` FROM `users_in_projects` WHERE `project_id` = ?);",
		projectId)

	if err != nil {
		return true, nil
	}

	defer rows.Close()

	for rows.Next() {
		accountOne := models.Account{}
		if err := rows.Scan(&accountOne.Uuid, &accountOne.Name, &accountOne.DisplayName); err != nil {
			log.Printf("ProjectsModule.GetUser Error: %+v\n", err)
			return true, nil
		}
		accounts = append(accounts, accountOne)
	}

	return false, accounts
}

func (self *ProjectsModule) Update(project *models.Project) bool {
	result, err := self.db.Exec("UPDATE `projects` SET `name` = ?, `description` = ? WHERE `uuid` = ?;",
		project.Name, project.Description, project.Uuid)

	if err != nil {
		log.Printf("ProjectsModule.Update Error: %+v\n", err)
		return true
	}

	_, err = result.RowsAffected()
	if err != nil {
		log.Printf("ProjectsModule.Update Error: %+v\n", err)
		return true
	}

	return false
}

func (self *ProjectsModule) IsIn(userUuid, projectId int) (isIn bool) {
	rows, err := self.db.Query("SELECt count(`user_id`) FROM `users_in_projects` WHERE `user_id` = ? AND `project_id` = ?;",
		userUuid, projectId)

	if err != nil {
		return false
	}

	defer rows.Close()

	for rows.Next() {
		var inCount int
		if err := rows.Scan(&inCount); err != nil {
			log.Printf("ProjectsModule.IsIn Error: %+v\n", err)
			return false
		}
		return inCount > 0
	}

	return false
}

func (self *ProjectsModule) Delete(projectId int) (isError bool) {
	_, err := self.db.Exec("DELETE FROM todo_list WHERE project = ?;", projectId)

	if err != nil {
		log.Printf("ProjectsModule.Delete Error (todo_list): %+v\n", err)
		return true
	}

	_, err = self.db.Exec("DELETE FROM users_in_projects WHERE project_id = ?;", projectId)

	if err != nil {
		log.Printf("ProjectsModule.Delete Error (users_in_projects): %+v\n", err)
		return true
	}

	_, err = self.db.Exec("DELETE FROM projects WHERE uuid = ?;", projectId)

	if err != nil {
		log.Printf("ProjectsModule.Delete Error (projects): %+v\n", err)
		return true
	}

	return false
}
