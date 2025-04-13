package usn

type reader struct {
}

func newReader() *reader {
	return &reader{}
}

type Record struct {
	FileID   uint64
	ParentID uint64
}

func (r *reader) ReadLogs(usnLogsPath string) ([]Record, error) {
	return nil, nil
}
