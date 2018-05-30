package logger

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLog(t *testing.T) {
	f, err := testTempFileName("testdata")
	if err != nil {
		t.Fatalf("error opening temp file: %s\n", err)
	}

	defer func(f *os.File) {
		os.Remove(f.Name())
	}(f)

	l := New(os.Stdout)
	l.Info("test info", "addr", 80)
	l.Error("test error", "test", "abc")

	f.Close()

	logFile, err := os.Open(f.Name())
	if err != nil {
		t.Fatalf("error opening log file for reading: %s\n", err)
	}
	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)
	scanner.Scan()
	line := scanner.Text()

	if !strings.Contains(line, "\"info\":\"test info\"") {
		t.Fatalf("missing info level test info element\n")
	}
	if !strings.Contains(line, "\"addr\":80") {
		t.Fatalf("missing addr 80 element\n")
	}

	scanner.Scan()
	line = scanner.Text()
	if !strings.Contains(line, "\"error\":\"test error\"") {
		t.Fatalf("missing error level test error element\n")
	}
	if !strings.Contains(line, "\"test\":\"abc\"") {
		t.Fatalf("missing test abc element\n")
	}
}

func testTempFileName(dir string) (*os.File, error) {
	f, err := ioutil.TempFile(dir, "test_")
	if err != nil {
		return nil, err
	}
	return f, nil
}
