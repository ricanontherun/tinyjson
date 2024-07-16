package log

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand/v2"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testLogFileName = "test.log"
)

var words [][]byte
var sentences [][]byte
var paragraphs [][]byte

var initOnce func() = sync.OnceFunc(func() {
	initializeTestData()
})

func initializeTestData() {
	wordsContent := readTestFile(t, "test-files/words.txt")
	words := bytes.Split(wordsContent, []byte{' '})
	require.NotEmpty(t, words, "words should have been loaded and parsed")
}

func TestInitializeEmptyFileAndIterate(t *testing.T) {
	assertRemoveFile(t, testLogFileName)

	l := openLog(testLogFileName, t)
	defer closeLog(l, t)

	logIterator := assertIterator(t, l)
	assertIteratorExhausted(t, logIterator)
}

func TestInitializeFileAndIteratorSingleSmallRecord(t *testing.T) {
	assertRemoveFile(t, testLogFileName)

	log := openLog(testLogFileName, t)
	defer closeLog(log, t)

	assertLogWrite(t, log, "Testing Log Line 100")

	logIterator := assertIterator(t, log)
	assertNextEquals(t, logIterator, "Testing Log Line 100")
	assertIteratorExhausted(t, logIterator)
}

func TestThing(t *testing.T) {
	lines := bytes.Split(readTestFile(t, "test-files/beyond-lies-the-wub.txt"), []byte("\n"))

	require.NotEmpty(t, lines, "test data from file should not be empty.")

	for _, line := range lines {
		t.Log(string(line))
	}
}

func TestInitializeFileAndIteratorMultipleSmallRecord(t *testing.T) {
	assertRemoveFile(t, testLogFileName)

	log := openLog(testLogFileName, t)
	defer closeLog(log, t)

	assertLogWrite(t, log, "Testing Log Line 0")
	assertLogWrite(t, log, "Testing Log Line 1")
	assertLogWrite(t, log, "Testing Log Line 2")

	logIterator := assertIterator(t, log)

	assertNextEquals(t, logIterator, "Testing Log Line 2")
	assertNextEquals(t, logIterator, "Testing Log Line 1")
	assertNextEquals(t, logIterator, "Testing Log Line 0")
}

func TestInitializeFileAndIteratorMultipleLongRecords(t *testing.T) {
	assertRemoveFile(t, testLogFileName)

	testLogLines := readTestLines(t)

	log := openLog(testLogFileName, t)
	defer closeLog(log, t)

	assertLogWriteBytes(t, log, testLogLines[0])
	assertLogWriteBytes(t, log, testLogLines[1])
	assertLogWriteBytes(t, log, testLogLines[2])

	logIterator := assertIterator(t, log)

	assertNextEqualsBytes(t, logIterator, testLogLines[2])
	assertNextEqualsBytes(t, logIterator, testLogLines[1])
	assertNextEqualsBytes(t, logIterator, testLogLines[0])
}

func readTestLines(t *testing.T) [][]byte {
	testDataFile, err := os.Open("long.lines")
	require.Nil(t, err)
	require.NotNil(t, testDataFile)
	testDataContent, err := io.ReadAll(testDataFile)
	require.Nil(t, err)
	require.NotEqual(t, 0, len(testDataContent))
	testLogLines := bytes.Split(testDataContent, []byte("\n"))
	require.NotEqual(t, 0, len(testLogLines))
	return testLogLines
}

func TestInitializeFileAndIterateMultipleLines(t *testing.T) {
	assertRemoveFile(t, testLogFileName)

	log := openLog(testLogFileName, t)
	defer closeLog(log, t)

	testDataFile, err := os.Open("long.lines")
	assert.Nil(t, err)
	assert.NotNil(t, testDataFile)

	testData, err := io.ReadAll(testDataFile)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(testData))

	testLogLines := strings.Split(string(testData), "\n")
	assert.NotEqual(t, 0, len(testLogLines))

	var writtenLines []string
	for _, line := range testLogLines {
		assertLogWrite(t, log, line)
		writtenLines = append(writtenLines, line)
	}

	logIterator := assertIterator(t, log)

	reverseIterate(writtenLines, func(data string) {
		assertNextEquals(t, logIterator, data)
	})

	// Iterator should be exhausted
	assertIteratorExhausted(t, logIterator)

	for i := range 10 {
		logLine := fmt.Sprintf("testing log line %d", 10+i)
		assert.Nil(t, log.Write([]byte(logLine)))
		writtenLines = append(writtenLines, logLine)
	}

	logIterator, err = log.Iterator()
	assert.NotNil(t, logIterator)
	assert.Nil(t, err)

	assert.Equal(t, 20, len(writtenLines))
	reverseIterate(writtenLines, func(data string) {
		assertNextEquals(t, logIterator, data)
	})
}

func TestWriteAndIteratorLongLines(t *testing.T) {
}

func generateBytes(min int, length int) []byte {
	lowerCaseLetters := "abcdefghijklmnopqrstuvwxyz"
	upperCaseLetters := strings.ToUpper(lowerCaseLetters)
	digits := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	var dictionary []byte

	dictionary = append(dictionary, []byte(lowerCaseLetters)...)
	dictionary = append(dictionary, digits...)
	dictionary = append(dictionary, []byte(upperCaseLetters)...)

	l := rand.IntN(length-min) + min
	byteBuffer := make([]byte, l)
	for i := range l {
		byteBuffer[i] = dictionary[rand.IntN(len(dictionary))]
	}

	return byteBuffer
}

func assertIteratorExhausted(t *testing.T, iterator Iterator) {
	readBytes, err := iterator.next()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "exhausted", "expected exhausted iterator")
	assert.Empty(t, readBytes)
}

func assertIterator(t *testing.T, l Log) Iterator {
	logIterator, err := l.Iterator()
	assert.Nil(t, err, "")
	assert.NotNil(t, logIterator)
	return logIterator
}

func readWordsFile(t *testing.T) [][]byte {
	content := readTestFile(t, "test-files/words.txt")

	words := bytes.Split(content, []byte{' '})
	require.NotEmpty(t, words)

	return words
}

func readTestFile(t *testing.T, filename string) []byte {
	content, err := os.ReadFile(filename)

	require.Nil(t, err)
	require.NotEmpty(t, content)

	return content
}
