package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

// reverse iterators emit records in reverse order
// init:
//	const defaultReadSizeBytes = 1024
//
//	set end = size(file) - 1
//	set start = Match.max(end - defaultReadSizeBytes, 0)

//	next:
//	set shouldLoop = true
// 	loop shouldLoop:
//		readSizes bytes from file into temp readSizes buffer
//			if EOF: # hit end of file, end loop.
//				shouldLoop = false
//
//			if temp buffer has DELIM:
//				# only take what we need.
//				append record buffer with 0 -> index of DELIM
//				EMIT record buffer
//				# setup the subsequent next() call
//				# start the next readSizes at the byte directly AFTER the DELIM position
//				start = (index of delim) + 1
//				end =
//
//			else:
//				append record buffer with everything in readSizes buffer.

const (
	defaultReadSizeBytes = 1024
)

func back(cursor, amount int64) int64 {
	return cursor - amount
}

type reverseIterator struct {
	log      Log
	fileSize int64

	rootLock *sync.RWMutex

	queuedReads   [][]byte
	nextQueuedIdx int
	firstRead     bool
	exhausted     bool
	readCursor    int64
	// Each time next() is called, READ_BUFFER_SIZE bytes are readSizes from the file into readBuffer, starting at readCursor, which is relative to the end of the file.
	readBuffer *bytes.Buffer

	nextRecord string
}

type logReader interface {
	Read() ([]byte, error)
}

type Iterator interface {
	next() ([]byte, error)
}

func newIterator(log Log, rootLock *sync.RWMutex) (Iterator, error) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	var readStart int64 = 0

	fileSize, err := log.size()
	if err != nil {
		return nil, err
	} else {
		if fileSize == 0 {
			fmt.Println("file is empty, returning pre-exhausted iterator. ")
			return exhaustedIterator(), nil
		}

		// Initialize the buffer bounds
		readStart = fileSize - 2 // ignore the trailing newline.

		fmt.Printf("iterator initialized, readCursor=%d, fileSize=%d\n", readStart, fileSize)
	}

	buffer := bytes.NewBuffer([]byte{})
	return &reverseIterator{log, fileSize, rootLock, [][]byte{}, -1, true, false, readStart, buffer, ""}, nil
}

func (iterator *reverseIterator) next() ([]byte, error) {
	if iterator.exhausted {
		slog.Debug("iterator is exhausted")
		return nil, errors.New("iterator is exhausted")
	}

	var workingRecord []byte

	slog.Info("Start next read loop")
	for {
		var readSize int64 = defaultReadSizeBytes

		// We're at the end of the file, read fewer bytes
		readStart := iterator.readCursor - readSize
		if readStart < 0 {
			slog.Debug("read start out of bounds", "readStart", readStart, "readSize", readSize)
			readSize += readStart
			readStart = 0
			slog.Debug("truncated readStart", "readStart", readStart, "readSize", readSize)
		}

		if readSize == -1 {
			iterator.exhausted = true
			return reverseBytes(workingRecord), nil
		}

		// Read readSize bytes starting at readStart
		slog.Debug("reading bytes", "readStart", readStart, "readSize", readSize)
		readBuffer := make([]byte, readSize+1)
		bytesRead, err := iterator.log.reader().ReadAt(readBuffer, readStart)
		slog.Debug("read bytes", "readStart", readStart, "readSize", readSize, "bytesRead", bytesRead)
		hadReadErr := err != nil
		if hadReadErr {
			isNotEndOfFile := !errors.Is(err, io.EOF)
			if isNotEndOfFile { // unexpected
				// Fewer bytes were read (not an EOF)
				// TODO: How do we handle this? Retry from attempted byte?
				readFewerBytesThanExpected := bytesRead != len(readBuffer)
				if readFewerBytesThanExpected {
					panic(errors.New("inconsistent read"))
				}

				panic(err)
			}
		}

		indexOfLastDelimiter := bytes.LastIndex(readBuffer, []byte("\n"))
		hasDelimiter := indexOfLastDelimiter != -1
		if hasDelimiter {
			workingRecord = append(workingRecord, reverseBytes(readBuffer[indexOfLastDelimiter+1:])...)
			iterator.readCursor -= int64(bytesRead - indexOfLastDelimiter)
			return reverseBytes(workingRecord), nil
		} else {
			// No delim, add to working buffer and move to next loop iteration.
			workingRecord = append(workingRecord, reverseBytes(readBuffer)...)
			iterator.readCursor -= int64(bytesRead)
			continue
		}
	}

	return nil, nil
}

func reverseBytes(source []byte) []byte {
	reversed := make([]byte, len(source))
	for write, read := 0, len(source)-1; read >= 0; write, read = write+1, read-1 {
		reversed[write] = source[read]
	}
	return reversed
}

// record1\n
// record2\n
// record3EOF
//	NO delimiter after last record.

type iteratorState interface {
}

type readSizes struct {
	readStart int64
	// When this is zero, we've hit the end of the file.
	readSize int
}

type Cursor interface {
	Next() readSizes
}

type cursorImpl struct {
	fileSize int64
}

func (c cursorImpl) Next() readSizes {
	return readSizes{}
}

func newCursor(fileSize int64) Cursor {
	return cursorImpl{fileSize}
}

func exhaustedIterator() Iterator {
	return &reverseIterator{
		log:           nil,
		fileSize:      0,
		rootLock:      nil,
		nextQueuedIdx: -1,
		queuedReads:   nil,
		firstRead:     false,
		exhausted:     true,
		readCursor:    0,
		readBuffer:    nil,
		nextRecord:    "",
	}
}

func findDelimiters(s []byte) []int {
	delims := []int{}

	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			delims = append(delims, i)
		}
	}

	return delims
}
