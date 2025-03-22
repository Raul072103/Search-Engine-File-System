package database

import (
	"MyFileExporer/common/models"
	"context"
	"database/sql"
)

type Repo struct {
	Files interface {
		Search(ctx context.Context, searchRequest FileSearchRequest) ([]models.File, error)
	}
}

func NewRepo(db *sql.DB) Repo {
	return Repo{Files: &fileRepo{db: db}}
}
