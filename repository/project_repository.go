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

func (p *ProjectRepository) IsInBulk(userUuid int, projectId []int) (isIn bool, err error) {
	c := 0
	err = p.db.Where("user_id = ? AND project_id IN (?)", userUuid, projectId).Count(&c).Error

	return c > 0, err
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
