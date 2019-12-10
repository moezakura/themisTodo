package repository

import (
	"github.com/jinzhu/gorm"
	"log"
	"themis.mox.si/themis/models"
	"themis.mox.si/themis/models/db"
)

type ProjectRepository struct {
	db *gorm.DB
}

func (p *ProjectRepository) Add(name, description string) (uuid int, error error) {
	addData := db.Project{
		Uuid:        0,
		Name:        name,
		Description: description,
	}
	err := p.db.Save(&addData).Error
	if err != nil {
		return 0, err
	}

	return addData.Uuid, nil
}

func (p *ProjectRepository) AddUser(userId, projectId int) bool {


	stmt, err := p.db.Prepare("INSERT INTO `users_in_projects` (`user_id`, `project_id`, `enable`, `expiration`) VALUES(?, ?, TRUE, NULL);")

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

func (p *ProjectRepository) GetProjects(userId int) (error bool, project []models.Project) {
	project = []models.Project{}

	rows, err := p.db.Query("SELECT `uuid`,`name`,`description` FROM `projects` WHERE `projects`.`uuid` IN (SELECT `project_id` FROM `users_in_projects` WHERE `user_id` = ? ORDER BY `user_id`);", userId)

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

func (p *ProjectRepository) GetProject(userId int) (error bool, project *models.Project) {
	project = &models.Project{}

	rows, err := p.db.Query("SELECT `uuid`,`name`,`description` FROM `projects` WHERE `uuid` = ?;", userId)

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

func (p *ProjectRepository) GetUser(projectId int) (error bool, accounts []models.Account) {
	accounts = make([]models.Account, 0)

	rows, err := p.db.Query("SELECT `uuid`, `name`, `displayName`, `icon_path` FROM `users` WHERE `users`.`uuid` IN (SELECT `user_id` FROM `users_in_projects` WHERE `project_id` = ?);",
		projectId)

	if err != nil {
		return true, nil
	}

	defer rows.Close()

	for rows.Next() {
		accountOne := models.Account{}
		if err := rows.Scan(&accountOne.Uuid, &accountOne.Name, &accountOne.DisplayName, &accountOne.IconPath); err != nil {
			log.Printf("ProjectsModule.GetUser Error: %+v\n", err)
			return true, nil
		}
		accounts = append(accounts, accountOne)
	}

	return false, accounts
}

func (p *ProjectRepository) Update(project *models.Project) bool {
	result, err := p.db.Exec("UPDATE `projects` SET `name` = ?, `description` = ? WHERE `uuid` = ?;",
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

func (p *ProjectRepository) IsIn(userUuid, projectId int) (isIn bool) {
	rows, err := p.db.Query("SELECt count(`user_id`) FROM `users_in_projects` WHERE `user_id` = ? AND `project_id` = ?;",
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

func (p *ProjectRepository) IsInBulk(userUuid, projectId []int) (isIn bool) {
	queryString := ""
	queryArray := make([]interface{}, 0)

	queryString += "("
	queryStringUser := ""
	for _, user := range userUuid {
		if len(queryStringUser) > 0 {
			queryStringUser += " OR "
		}
		queryStringUser += " `user_id` = ? "
		queryArray = append(queryArray, user)
	}
	queryString += queryStringUser + ") AND ("
	queryStringProject := ""
	for _, project := range projectId {
		if len(queryStringProject) > 0 {
			queryStringProject += " OR "
		}
		queryStringProject += " `project_id` = ? "
		queryArray = append(queryArray, project)
	}
	queryString += queryStringProject + ")"

	rows, err := p.db.Query("SELECT count(`user_id`) FROM `users_in_projects` WHERE "+queryString+";",
		queryArray...)

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

func (p *ProjectRepository) Delete(projectId int) (isError bool) {
	_, err := p.db.Exec("DELETE FROM todo_list WHERE project = ?;", projectId)

	if err != nil {
		log.Printf("ProjectsModule.Delete Error (todo_list): %+v\n", err)
		return true
	}

	_, err = p.db.Exec("DELETE FROM users_in_projects WHERE project_id = ?;", projectId)

	if err != nil {
		log.Printf("ProjectsModule.Delete Error (users_in_projects): %+v\n", err)
		return true
	}

	_, err = p.db.Exec("DELETE FROM projects WHERE uuid = ?;", projectId)

	if err != nil {
		log.Printf("ProjectsModule.Delete Error (projects): %+v\n", err)
		return true
	}

	return false
}

func (p *ProjectRepository) Leave(projectId int, userId int) (isError bool) {
	_, err := p.db.Exec("DELETE FROM users_in_projects WHERE project_id = ? AND user_id = ?;", projectId, userId)

	if err != nil {
		log.Printf("ProjectsModule.Delete Error (users_in_projects): %+v\n", err)
		return true
	}
	return false
}
