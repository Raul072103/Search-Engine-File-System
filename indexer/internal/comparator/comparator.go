package comparator

type Directory interface {
	CompareDirectory()
}

type directory struct {
}

func New() Directory {
	return &directory{}
}
