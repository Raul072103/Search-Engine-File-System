package database

import (
	"MyFileExporer/common/models"
	"context"
	"database/sql"
	"fmt"
)

type FileSearchRequest struct {
	Words     *[]string
	Extension *[]string
	Name      *string
}

type fileRepo struct {
	db *sql.DB
}

// Search searches for the files that match the FileSearchRequest
// TODO() find another way since this can be used to inject SQL querie
// TODO() find a way to regex match the names
func (r *fileRepo) Search(ctx context.Context, searchRequest FileSearchRequest) ([]models.File, error) {
	query := `
		SELECT 
		    f.id,
		    f.path, 
		    f.name, 
		    f.size,
		    f.mode, 
		    f.extension, 
		    f.updated_at
		FROM files AS f
		JOIN types AS t ON f.id = t.file_id
		JOIN contents AS c ON c.file_id = f.id
	`

	var argIdx = 0
	var args = make([]any, 0)

	if searchRequest.Extension != nil {
		extensionQueryCondition := " AND ("

		for i, extension := range *searchRequest.Extension {
			argIdx += 1
			args = append(args, extension)
			extensionQueryCondition += fmt.Sprintf(" f.extension = $%d", argIdx)

			if i < len(*searchRequest.Extension)-1 {
				extensionQueryCondition += " OR"
			} else {
				extensionQueryCondition += ")\n"
			}
		}

		query += extensionQueryCondition
	}

	if searchRequest.Name != nil {
		argIdx += 1
		args = append(args, searchRequest.Name)
		nameQueryCondition := fmt.Sprintf(" AND f.name LIKE $%d || '%%'\n", argIdx)
		query += nameQueryCondition
	}

	if searchRequest.Words != nil {
		wordsQueryCondition := " AND (c.searchable_tsv @@ to_tsquery('english', "

		for i, word := range *searchRequest.Words {
			argIdx += 1
			args = append(args, word)
			wordsQueryCondition += fmt.Sprintf("$%d", argIdx)

			if i < len(*searchRequest.Words)-1 {
				wordsQueryCondition += " || ' & ' || "
			} else {
				wordsQueryCondition += "))\n"
			}
		}
		query += wordsQueryCondition
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var files = make([]models.File, 0)

	for rows.Next() {
		var file models.File

		err := rows.Scan(
			&file.ID,
			&file.Path,
			&file.Name,
			&file.Size,
			&file.Mode,
			&file.Extension,
			&file.UpdatedAt)

		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}
