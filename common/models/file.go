package models

import (
	"fmt"
	"time"
)

type File struct {
	ID        int64
	Path      string
	Name      string
	Size      int64
	IsDir     bool
	Mode      uint32
	Extension string
	Content   string
	UpdatedAt time.Time
}

// String returns the File struct in a printable format.
// Used for debugging mainly.
func (f *File) String() string {
	var contentStr string

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
