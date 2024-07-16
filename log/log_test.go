package log

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func reverseIterate(things []string, visitor func(string)) {
	for i := len(things) - 1; i >= 0; i-- {
		visitor(things[i])
	}
}

const name = "./test.log"

func TestOpenLogNewFile(t *testing.T) {
	// delete file if exists.
	if _, statErr := os.Stat(name); statErr == nil {
		t.Log("deleting test.log from previous test run")
		assert.Nil(t, os.Remove(name))
	}

	assertOpenLogNoError(t)

	// File should exist.
	_, statErr := os.Stat(name)
	assert.Nil(t, statErr)
}

func TestOpenLogExistingFile(t *testing.T) {
	_, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	assert.Nil(t, err)
	assertOpenLogNoError(t)
}

func TestLogSize(t *testing.T) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	assert.Nil(t, err)

	bytes := []byte("these are some characters")
	bytesWritten, err := file.Write(bytes)
	assert.Nil(t, err)
	assert.Equal(t, len(bytes), bytesWritten)

	assert.Nil(t, file.Close())

	log := openLog(name, t)

	logLen, err := log.size()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(bytes)), logLen)
}

func TestInputValidation(t *testing.T) {
	log := openLog(name, t)

	writeErr := log.Write(nil)
	assert.EqualError(t, writeErr, "illegal argument: nil")

	writeErr = log.Write([]byte{})
	assert.EqualError(t, writeErr, "illegal argument: empty")
}

func assertNextEquals(t *testing.T, logIterator Iterator, expected string) {
	next, err := logIterator.next()
	assert.Nil(t, err)
	// TODO: How to compare two byte slices
	assert.Equal(t, expected, string(next), "expected next line to be '%s', got '%s'", expected, next)
}

func assertNextEqualsBytes(t *testing.T, logIterator Iterator, expected []byte) {
	next, err := logIterator.next()
	require.Nil(t, err)
	// TODO: How to compare two byte slices
	require.Equal(t, expected, next, "expected next line to be '%s', got '%s'", string(expected), string(next))
}

func assertOpenLogNoError(t *testing.T) {
	filePath := "./test.log"
	l := openLog(filePath, t)
	defer closeLog(l, t)
}

func closeLog(log Log, t *testing.T) {
	fmt.Println("closing log")
	assert.Nil(t, log.Close(), "Close failure")
}

// Create a new log instance.
//
//	assert that Open() returns a non-nil log, nil err
func openLog(path string, t *testing.T) Log {
	log, err := Open(path)
	assert.Nilf(t, err, "shouldn't have failed to open the file.")
	assert.NotNil(t, log, "log instance should be non-nil")
	return log
}

func assertRemoveFile(t *testing.T, fileName string) {
	if err := os.Remove(fileName); err != nil {
		assert.FailNow(t, "failed to delete test.log")
	}
}

func assertLogWrite(t *testing.T, log Log, data string) {
	bytesToWrite := []byte(data)
	assert.Nil(t, log.Write(bytesToWrite), "write should not have returned err")
}

func assertLogWriteBytes(t *testing.T, log Log, bytes []byte) {
	require.Nil(t, log.Write(bytes), "write should not have returned err")
}
