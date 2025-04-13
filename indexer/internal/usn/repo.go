package usn

type Repo struct {
	Executor interface {
		ExecuteReadUSNJournal() error
		ExecuteQueryUSNJournal() error
	}

	Reader interface {
		ReadLogs(string) ([]Record, error)
	}
}

func NewRepo(executorConfig ExecutorConfig) Repo {
	return Repo{
		Executor: newExecutor(executorConfig),
		Reader:   newReader(),
	}
}
