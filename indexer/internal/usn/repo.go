package usn

type Repo struct {
	Executor interface {
		ExecuteReadUSNJournal() error
		ExecuteQueryUSNJournal() error
	}

	Parser interface {
		ReadLogs(string) ([]Record, error)
	}
}

func NewRepo(executorConfig ExecutorConfig) Repo {
	return Repo{
		Executor: newExecutor(executorConfig),
		Parser:   newParser(),
	}
}
