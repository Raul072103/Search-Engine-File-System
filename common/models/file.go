package models

import (
	"fmt"
	"strings"
	"time"
)

type File struct {
	ID        int64
	Path      string
	Name      string
	Size      int64
	IsDir     bool
	Mode      int32
	Extension string
	Content   []string
	UpdatedAt time.Time
}

// String returns the File struct in a printable format.
// Used for debugging mainly.
func (f *File) String() string {
	var contentStr string
	if f.Content != nil {
		contentStr = strings.Join(f.Content, ", ")
	} else {
		contentStr = "nil"
	}

	return fmt.Sprintf(
		"File{\n"+
			"ID: %d\n"+
			"Path: %q\n"+
			"Name: %q\n"+
			"Size: %d\n"+
			"IsDir: %t\n"+
			"Mode: %o\n"+
			"Extension: %q\n"+
			"Content: [%s]\n"+
			"UpdatedAt: %s\n"+
			"}",
		f.ID,
		f.Path,
		f.Name,
		f.Size,
		f.IsDir,
		f.Mode,
		f.Extension,
		contentStr,
		f.UpdatedAt.Format(time.RFC3339),
	)
}
