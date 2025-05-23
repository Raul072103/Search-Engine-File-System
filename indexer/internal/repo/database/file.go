package database

import (
	"MyFileExporer/common/models"
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

const (
	databaseTextCache = 200
)

type fileRepo struct {
	db       *sql.DB
	typesMap models.FileTypesConfig
}

// Insert inserts a file into the database.
// It only inserts the first 200 characters into the content field, if there aren't that many it will insert only the
// content passed to the File instance.
// If a resource with the same path is trying to be inserted an error of type ErrConflict will be thrown.
func (r *fileRepo) Insert(ctx context.Context, file *models.File) error {
	return withTransaction(r.db, ctx, func(tx *sql.Tx) error {
		query := `
		INSERT INTO files (path, name, size, mode, extension, file_id, parent_id, rank, hash, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id
	`

		ctx, cancel := context.WithTimeout(ctx, InsertFileTimeoutDuration)
		defer cancel()

		err := tx.QueryRowContext(
			ctx,
			query,
			file.Path,
			file.Name,
			file.Size,
			file.Mode,
			file.Extension,
			file.WindowsFileID,
			file.ParentFileID,
			file.Rank,
			file.Hash,
			file.UpdatedAt,
		).Scan(&file.ID)

		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				return ErrConflict
			}
			return err
		}

		err = r.insertFileType(ctx, tx, file)
		if err != nil {
			return err
		}

		if r.typesMap.TypesMapping[file.Type.TypeID] == "txt" ||
			r.typesMap.TypesMapping[file.Type.TypeID] == "code" {
			err = r.insertTxtFileContent(ctx, tx, file)
			if err != nil {
				return err
			}
		} else {
			err = r.insertFileContentForPreview(ctx, tx, file)
			if err != nil {
				return err
			}
		}

		return err
	})
}

// Update updates a file with the given WindowsFileID based on the non-zero/non-empty fields set in the File instance.
func (r *fileRepo) Update(ctx context.Context, file *models.File) error {
	return withTransaction(r.db, ctx, func(tx *sql.Tx) error {
		fileQuery := `
			UPDATE files
			SET path = COALESCE(NULLIF($1::TEXT, ''), path),
				name = COALESCE(NULLIF($2::TEXT, ''), name),
				size = COALESCE(NULLIF($3::BIGINT, 0), size),
				mode = COALESCE(NULLIF($4::BIGINT, 0), mode),
				extension = COALESCE(NULLIF($5::TEXT, ''), extension),
				parent_id = COALESCE(NULLIF($6::BIGINT, 0), parent_id),
				rank = COALESCE(NULLIF($7::INT, 0), rank),
				hash = COALESCE(NULLIF($8::TEXT, ''), hash),
				updated_at = COALESCE(NULLIF($9, to_timestamp(0)), updated_at)
			WHERE file_id = $10
			RETURNING id
		`

		err := tx.QueryRowContext(ctx, fileQuery,
			file.Path,
			file.Name,
			file.Size,
			file.Mode,
			file.Extension,
			file.ParentFileID,
			file.Rank,
			file.Hash,
			file.UpdatedAt,
			file.WindowsFileID,
		).Scan(&file.ID)
		if err != nil {
			return err
		}

		if r.typesMap.TypesMapping[file.Type.TypeID] == "txt" {
			if err := r.updateFileContent(ctx, tx, file); err != nil {
				return err
			}
		}

		if err := r.updateFileType(ctx, tx, file); err != nil {
			return err
		}

		return nil
	})
}

// Delete permanently deletes a file with the given path.
func (r *fileRepo) Delete(ctx context.Context, file *models.File) error {
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

// GetAllDirectoriesFileIDs returns a list with all directories File IDs from the database
func (r *fileRepo) GetAllDirectoriesFileIDs(ctx context.Context) ([]int64, error) {
	query := `
		SELECT file_id
		FROM files 
		WHERE extension = '';
	`

	ctx, cancel := context.WithTimeout(ctx, InsertFileTimeoutDuration)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var directoriesIDs = make([]int64, 0)
	for rows.Next() {
		var directoryID int64
		err := rows.Scan(&directoryID)
		if err != nil {
			return nil, err
		}

		directoriesIDs = append(directoriesIDs, directoryID)
	}

	return directoriesIDs, nil
}

// GetFileByWindowsFileID returns the file with the given Windows File ID
func (r *fileRepo) GetFileByWindowsFileID(ctx context.Context, fileID int64) (models.File, error) {
	query := `
		SELECT id, path, name, size, mode, extension, file_id, parent_id, rank, hash, updated_at
		FROM files
		WHERE file_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, InsertFileTimeoutDuration)
	defer cancel()

	var file models.File
	err := r.db.QueryRowContext(ctx, query, fileID).Scan(
		&file.ID,
		&file.Path,
		&file.Name,
		&file.Size,
		&file.Mode,
		&file.Extension,
		&file.WindowsFileID,
		&file.ParentFileID,
		&file.Rank,
		&file.Hash,
		&file.UpdatedAt,
	)

	return file, err
}

// GetAllFilesWithParent returns all files that have as parent the given Windows File ID
func (r *fileRepo) GetAllFilesWithParent(ctx context.Context, parentID int64) ([]models.File, error) {
	query := `
		SELECT id, path, name, size, mode, extension, file_id, parent_id, rank, hash, updated_at
		FROM files
		WHERE parent_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, InsertFileTimeoutDuration)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		err := rows.Scan(
			&file.ID,
			&file.Path,
			&file.Name,
			&file.Size,
			&file.Mode,
			&file.Extension,
			&file.WindowsFileID,
			&file.ParentFileID,
			&file.Rank,
			&file.Hash,
			&file.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// DeleteAllUnderDirectory deletes all descendants of directory, including directory itself
func (r *fileRepo) DeleteAllUnderDirectory(ctx context.Context, directory *models.File) error {
	query := `
		WITH RECURSIVE descendants AS (
		  SELECT id, file_id
		  FROM files
		  WHERE file_id = $1  -- root directory's file_id
		
		  UNION ALL
		
		  SELECT f.id, f.file_id
		  FROM files f
		  INNER JOIN descendants d ON f.parent_id = d.file_id
		)
		DELETE FROM files
		WHERE file_id IN (SELECT file_id FROM descendants);
	`

	ctx, cancel := context.WithTimeout(ctx, DeleteRecursiveFilesTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, directory.WindowsFileID)
	if err != nil {
		return err
	}

	return nil
}

// insertTxtFileContent helper method to insert the content of the file in "contents" table5
func (r *fileRepo) insertTxtFileContent(ctx context.Context, tx *sql.Tx, file *models.File) error {
	fileContent := &file.Content

	query := `
		INSERT INTO contents (file_id, content_text, searchable_tsv, updated_at)
		VALUES ($1, $2::TEXT, to_tsvector('english', $3), $4) RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, InsertFileTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		file.ID,
		fileContent.Text[:min(databaseTextCache, len(fileContent.Text))],
		fileContent.Text,
		file.UpdatedAt,
	).Scan(&fileContent.ID)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return err
}

// insertFileContentForPreview helper method to insert the content of the file in "contents" table.
// This function is specific for files that tsvector cannot perform search on them but the user should be able
// to preview them
func (r *fileRepo) insertFileContentForPreview(ctx context.Context, tx *sql.Tx, file *models.File) error {
	fileContent := &file.Content

	query := `
		INSERT INTO contents (file_id, content_text, updated_at)
		VALUES ($1, $2::TEXT, $3) RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, InsertFileTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		file.ID,
		fileContent.Text[:min(databaseTextCache, len(fileContent.Text))],
		file.UpdatedAt,
	).Scan(&fileContent.ID)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return err
}

// insertFileType helper method to insert the type of the file in "types" table
func (r *fileRepo) insertFileType(ctx context.Context, tx *sql.Tx, file *models.File) error {
	fileType := &file.Type

	query := `
		INSERT INTO types (file_id, type, updated_at)
		VALUES ($1, $2, $3) RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, InsertFileTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		file.ID,
		fileType.TypeID,
		file.UpdatedAt,
	).Scan(&fileType.ID)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return err
}

// updateFileContent updates the content and searchable_tsv in the contents table.
func (r *fileRepo) updateFileContent(ctx context.Context, tx *sql.Tx, file *models.File) error {
	contentQuery := `
	UPDATE contents
	SET content_text = COALESCE(NULLIF($1, ''), contents.content_text),
	    searchable_tsv = CASE 
	                        WHEN (NULLIF($2::TEXT, '')) IS NOT NULL THEN to_tsvector('english', $3)
	                        ELSE contents.searchable_tsv 
	                    END,
	    content_bytes = COALESCE(NULLIF($3::bytea, ''), contents.content_bytes),
	    updated_at = $4	
	WHERE file_id = $5
	`
	_, err := tx.ExecContext(ctx, contentQuery,
		file.Content.Text[:min(databaseTextCache, len(file.Content.Text))],
		file.Content.Text,
		file.Content.Bytes,
		file.UpdatedAt,
		file.ID,
	)
	return err
}

// updateFileType updates the file type in the types table.
func (r *fileRepo) updateFileType(ctx context.Context, tx *sql.Tx, file *models.File) error {
	typeQuery := `
	UPDATE types
	SET type = COALESCE(NULLIF($1, -1), types.type),
	    updated_at = $2
	WHERE file_id = $3
	`
	_, err := tx.ExecContext(ctx, typeQuery,
		file.Type.TypeID,
		file.UpdatedAt,
		file.ID,
	)
	return err
}
