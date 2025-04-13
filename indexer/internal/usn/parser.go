package usn

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ErrUSNJournalFormatChanged = errors.New("USN journal format changed")
)

type parser struct {
}

func newParser() *parser {
	return &parser{}
}

type Record struct {
	FileID   uint64
	ParentID uint64
}

func (r Record) String() string {
	return fmt.Sprintf("Record{FileID: %d, ParentID: %d}", r.FileID, r.ParentID)
}

func (r *parser) ReadLogs(usnLogsPath string) ([]Record, error) {
	bytesContent, err := os.ReadFile(usnLogsPath)
	if err != nil {
		return nil, err
	}

	var lines []string
	lines = strings.Split(string(bytesContent), "\n")
	if len(lines) < 8 {
		return nil, err
	}

	// Skip first 7 lines (contains journal query information)
	lines = lines[7:]

	// File ID           : 00000000000000000007000000018635
	// Parent file ID    : 0000000000000000000400000001a0f8

	i := 0
	records := make([]Record, 0)

	for i < len(lines) {
		paragraph := lines[i : i+13]

		fileIDLine := paragraph[6]
		parentIDLine := paragraph[7]

		if strings.HasPrefix(fileIDLine, "File ID") && strings.HasPrefix(parentIDLine, "Parent file ID") {
			fileIDParts := strings.Split(fileIDLine, ":")
			parentIDParts := strings.Split(parentIDLine, ":")

			if len(fileIDParts) < 2 || len(parentIDParts) < 2 {
				return nil, ErrUSNJournalFormatChanged
			}

			fileID, err := parseUSNHexString(strings.TrimSpace(fileIDParts[1]))
			if err != nil {
				return nil, err
			}

			parentID, err := parseUSNHexString(strings.TrimSpace(parentIDParts[1]))
			if err != nil {
				return nil, err
			}

			record := Record{
				FileID:   fileID,
				ParentID: parentID,
			}

			records = append(records, record)
		} else {
			return nil, ErrUSNJournalFormatChanged
		}

		i += 14
	}

	return records, nil
}

func parseUSNHexString(usnStr string) (uint64, error) {
	return strconv.ParseUint(usnStr, 16, 64)
}
