package models

import (
	"fmt"
	"time"
)

type FileType struct {
	ID        int64
	FileID    int64
	TypeID    int32
	UpdatedAt time.Time
}

// String returns the FileType struct in a printable format.
// Used for debugging mainly.
func (f *FileType) String() string {
	return fmt.Sprintf("FileType{ID: %d, FileID: %d, Type: %v, UpdatedAt: %s}", f.ID, f.FileID, f.TypeID, f.UpdatedAt.Format(time.RFC3339))
}
