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

func (p *ProjectRepository) AddUser(userId, projectId int) error {
	addData := db.UsersInProject{
		UserId:     userId,
		ProjectId:  projectId,
		Enable:     true,
		Expiration: nil,
	}
	return p.db.Save(&addData).Error
}

func (p *ProjectRepository) GetProjectsByUserId(userId int) (projects []db.Project, err error) {
	projects = make([]db.Project, 0)

	subQuery := p.db.Table("users_in_projects").Select("project_id").Where("user_id = ?", userId).Order("user_id")
	err = p.db.Where("uuid = ?", subQuery.QueryExpr()).Find(&projects).Error
	return projects, err
}

func (p *ProjectRepository) GetProjectById(id int) (project *db.Project, err error) {
	err = p.db.First(project, "uuid = ?", id).Error
	return project, err
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
