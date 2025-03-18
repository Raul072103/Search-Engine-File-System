package database

import (
	"MyFileExporer/common/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"
)

var (
	InsertFileTimeoutDuration = time.Second * 120

	ErrConflict = errors.New("file already exists in database")
)

type Repo interface {
	InsertFile(ctx context.Context, file *models.File) error
	UpdateFile(ctx context.Context, file *models.File) error
	DeleteFile(ctx context.Context, file *models.File) error
}

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repo {
	return &repo{db: db}
}

// InsertFile inserts a file into the database.
// It only inserts the first 200 characters into the content field, if there aren't that many it will insert only the
// content passed to the File instance.
// If a resource with the same path is trying to be inserted an error of type ErrConflict will be thrown.
func (r *repo) InsertFile(ctx context.Context, file *models.File) error {
	query := `
		INSERT INTO files (path, name, size, mode, extension, updated_at, content)
		VALUES ($1, $2, $3, $4, $5, $6, $7,
		        COALESCE($8::TEXT),
		        COALESCE(to_tsvector('english', $9), NULL))
		        RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, InsertFileTimeoutDuration)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx,
		query,
		file.Path,
		file.Name,
		file.Size,
		file.IsDir,
		file.Mode,
		file.Extension,
		file.UpdatedAt,
		file.Content[:min(len(file.Content), 200)],
		file.Content,
	).Scan(&file.ID)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return err
}

// UpdateFile updates a file with the given path based on the "non-null" fields set in the File instance.
func (r *repo) UpdateFile(ctx context.Context, file *models.File) error {
	query := `UPDATE files SET `
	var args []interface{}
	argPos := 1

	if file.Name != "" {
		query += fmt.Sprintf("name = $%d, ", argPos)
		args = append(args, file.Name)
		argPos++
	}

	if file.Size > 0 {
		query += fmt.Sprintf("size = $%d, ", argPos)
		args = append(args, file.Size)
		argPos++
	}

	if file.IsDir {
		query += fmt.Sprintf("is_dir = $%d, ", argPos)
		args = append(args, file.IsDir)
		argPos++
	}

	if file.Mode != 0 {
		query += fmt.Sprintf("mode = $%d, ", argPos)
		args = append(args, file.Mode)
		argPos++
	}

	if file.Extension != "" {
		query += fmt.Sprintf("extension = $%d, ", argPos)
		args = append(args, file.Extension)
		argPos++
	}

	if !file.UpdatedAt.IsZero() {
		query += fmt.Sprintf("updated_at = $%d, ", argPos)
		args = append(args, file.UpdatedAt)
		argPos++
	}

	if file.Content != "" {
		query += fmt.Sprintf("content = $%d, ", argPos)
		args = append(args, file.Content[:min(len(file.Content), 200)])
		argPos++
	}

	if file.Content != "" {
		query += fmt.Sprintf("searchable_tsv = $%d ", argPos)
		args = append(args, toTsvector(file.Content)) // Use a helper function for tsvector conversion
		argPos++
	} else {
		query += fmt.Sprintf("searchable_tsv = NULL ")
	}

	query += fmt.Sprintf("WHERE path = $%d RETURNING id", argPos)
	args = append(args, file.Path)

	ctx, cancel := context.WithTimeout(ctx, InsertFileTimeoutDuration)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&file.ID)

	return err
}

// DeleteFile permanently deletes a file with the given path.
func (r *repo) DeleteFile(ctx context.Context, file *models.File) error {
	query := `
		DELETE 
		FROM files
		WHERE path = $1
	`

	ctx, cancel := context.WithTimeout(ctx, InsertFileTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		query,
		file.Path)

	return err
}

// Helper function to create the tsvector
func toTsvector(content string) string {
	return fmt.Sprintf("to_tsvector('english', %s)", content)
}
