package models

import (
	"fmt"
	"time"
)

type FileContent struct {
	ID        int64
	FileID    int64
	Text      string
	Bytes     []byte
	UpdatedAt time.Time
}

// String returns the FileContent struct in a printable format.
// Used for debugging mainly.
func (f *FileContent) String() string {
	return fmt.Sprintf("FileContent{ID: %d, FileID: %d, Text: %q, Bytes: %d bytes, UpdatedAt: %s}", f.ID, f.FileID, f.Text, len(f.Bytes), f.UpdatedAt.Format(time.RFC3339))
}
