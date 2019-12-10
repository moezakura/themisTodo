package module

import (
	"themis.mox.si/themis/models"
	"database/sql"
	"log"
)

type ProjectsModule struct {
	db *sql.DB
}

func NewProjectsModule(db *sql.DB) *ProjectsModule {
	return &ProjectsModule{db}
}


