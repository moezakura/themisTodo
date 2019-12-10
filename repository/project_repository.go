package repository

import (
	"github.com/jinzhu/gorm"
	"themis.mox.si/themis/models"
	"themis.mox.si/themis/models/db"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (p *ProjectRepository) Add(name, description string) (uuid int, error error) {
	addData := db.Project{
		Name:        name,
		Description: description,
	}
	err := p.db.Create(&addData).Error
	if err != nil {
		return 0, err
	}

	return addData.Uuid, nil
}

func (p *ProjectRepository) AddUser(userId, projectId int) error {
	addData := db.UsersInProject{
		UserId:    userId,
		ProjectId: projectId,
		Enable:    true,
	}
	return p.db.Create(&addData).Error
}

func (p *ProjectRepository) GetProjectsByUserId(userId int) (projects []db.Project, err error) {
	projects = make([]db.Project, 0)

	subQuery := p.db.Table("users_in_projects").Select("project_id").Where("user_id = ?", userId).Order("user_id")
	err = p.db.Where("projects.uuid IN (?)", subQuery.QueryExpr()).Find(&projects).Error
	return projects, err
}

func (p *ProjectRepository) GetProjectById(id int) (project *db.Project, err error) {
	project = &db.Project{}
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

func (p *ProjectRepository) Delete(projectId int) (err error) {
	err = p.db.Exec("DELETE FROM todo_list WHERE project = ?;", projectId).Error

	if err != nil {
		return err
	}

	err = p.db.Exec("DELETE FROM users_in_projects WHERE project_id = ?;", projectId).Error

	if err != nil {
		return err
	}

	err = p.db.Exec("DELETE FROM projects WHERE uuid = ?;", projectId).Error

	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectRepository) Leave(projectId int, userId int) (err error) {
	return p.db.Delete(&db.UsersInProject{}, "project_id = ? user_id = ?", projectId, userId).Error
}
