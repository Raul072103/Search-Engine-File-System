package usn

import (
	"errors"
	"math"
)

var (
	ErrFileIDTooBig = errors.New("file ID too big")
)

type differenceFinder struct {
}

func newDifferenceFinder() *differenceFinder {
	return &differenceFinder{}
}

func (d *differenceFinder) FindUpdatedDirectories(records []Record, parents map[int64]any) ([]int64, error) {
	var differentDirectories = make([]int64, 0)

	for _, record := range records {
		if record.ParentID > math.MaxInt64 {
			return nil, ErrFileIDTooBig
		}

		parentID := int64(record.ParentID)
		if _, exists := parents[parentID]; exists {
			differentDirectories = append(differentDirectories, parentID)
			delete(parents, parentID)
		}
	}

	return differentDirectories, nil
}
