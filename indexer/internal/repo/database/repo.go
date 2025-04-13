package database

import (
	"MyFileExporer/common/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	InsertFileTimeoutDuration           = time.Second * 120
	DeleteRecursiveFilesTimeoutDuration = time.Second * 360

	ErrConflict = errors.New("file already exists in database")
)

type Repo struct {
	Files interface {
		Insert(ctx context.Context, file *models.File) error
		Update(ctx context.Context, file *models.File) error
		Delete(ctx context.Context, file *models.File) error
		DeleteAllUnderDirectory(ctx context.Context, directory *models.File) error
		GetAllDirectoriesFileIDs(ctx context.Context) ([]int64, error)
		GetFileByWindowsFileID(ctx context.Context, fileID int64) (models.File, error)
		GetAllFilesWithParent(ctx context.Context, parentID int64) ([]models.File, error)
	}
}

func NewRepo(db *sql.DB, typesConfig models.FileTypesConfig) Repo {
	return Repo{
		Files: &fileRepo{db: db, typesMap: typesConfig},
	}
}

func withTransaction(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
