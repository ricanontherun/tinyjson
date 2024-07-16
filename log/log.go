package log

import (
	"errors"
	"io"
	"os"
	"sync"
)

// Log is A delimited, Append-only log
type Log interface {
	Write([]byte) error
	Iterator() (Iterator, error)
	Close() error

	size() (int64, error)
	reader() io.ReaderAt
}

type logImpl struct {
	logPath    string
	fileHandle *os.File

	rootLock *sync.RWMutex
}

func (log *logImpl) size() (int64, error) {
	log.rootLock.Lock()
	defer log.rootLock.Unlock()
	stat, _ := log.fileHandle.Stat()
	return stat.Size(), nil
}

func Open(path string) (Log, error) {
	// Open file (research how file locking works)
	const openFileOptions = os.O_RDWR | os.O_APPEND | os.O_CREATE
	const openFilePerms = 0600
	if file, err := os.OpenFile(path, openFileOptions, openFilePerms); err != nil {
		return nil, errors.New("failed to open file: " + err.Error())
	} else {
		return &logImpl{path, file, &sync.RWMutex{}}, nil
	}
}

func (log *logImpl) Write(data []byte) error {
	if data == nil {
		return errors.New("illegal argument: nil")
	}

	if len(data) == 0 {
		return errors.New("illegal argument: empty")
	}

	data = append(data, '\n')

	log.rootLock.Lock()
	defer log.rootLock.Unlock()
	if _, err := log.fileHandle.Write(data); err == nil {
		return nil
	} else {
		if errors.Is(err, io.ErrShortWrite) {
			// Short write. Retry until succeeds OR exhausted?
		}
		return err
	}
}

func (log *logImpl) Iterator() (Iterator, error) {
	return newIterator(log, log.rootLock)
}

func (log *logImpl) reader() io.ReaderAt {
	return log.fileHandle
}

func (log *logImpl) Close() error {
	return log.fileHandle.Close()
}

// Implements a reverse Iterator.
//
//	Log - dumb byte appender/reverse-Iterator
