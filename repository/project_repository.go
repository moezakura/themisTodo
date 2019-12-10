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
	err = p.db.Where("projects.uuid = ?", subQuery.QueryExpr()).Find(&projects).Error
	return projects, err
}

func (p *ProjectRepository) GetProjectById(id int) (project *db.Project, err error) {
	err = p.db.First(project, "uuid = ?", id).Error
	return project, err
}

func (p *ProjectRepository) GetUserListInProject(projectId int) (accounts []models.Account, err error) {
	accounts = make([]models.Account, 0)

	subQuery := p.db.Table("users_in_projects").Select("user_id").Where("project = ?", projectId)
	err = p.db.Where("users.uuid IN (?)", subQuery.QueryExpr()).Find(&accounts).Error
	return accounts, err
}

func (p *ProjectRepository) Update(project *db.Project) error {
	return p.db.Save(project).Error
}

func (p *ProjectRepository) IsIn(userUuid, projectId int) (isIn bool, err error) {
	c := 0
	err = p.db.Model(&db.UsersInProject{}).Where("`user_id` = ? AND `project_id` = ?", userUuid, projectId).Count(&c).Error

	return c > 0, err
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
