package log

//

type concurrentFile interface {
	write(data []byte) error
}

type appendOnlyHandle interface {
}

func NewHandle() (concurrentFile, error) {
	return nil, nil
}
