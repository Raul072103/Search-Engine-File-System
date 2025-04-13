package usn

type Repo struct {
	Executor interface {
		ExecuteReadUSNJournal() error
		ExecuteQueryUSNJournal() error
	}

	Parser interface {
		ReadLogs(string) ([]Record, error)
	}

	DifferenceFinder interface {
		FindUpdatedDirectories([]Record, map[int64]any) ([]int64, error)
	}
}

func NewRepo(executorConfig ExecutorConfig) Repo {
	return Repo{
		Executor:         newExecutor(executorConfig),
		Parser:           newParser(),
		DifferenceFinder: newDifferenceFinder(),
	}
}
