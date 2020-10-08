package main

import (
	"bufio"
	"io"

	"github.com/pkg/errors"
)

// CSVIOReader implements the CSVIOReader interface
type CSVIOReader struct {
}

// NewCSVIOReader initializes a new CSV file reader
func NewCSVIOReader() CSVIOReader {
	return CSVIOReader{}
}

// ReadAll returns the contents of a CSV file as struct
func (reader CSVIOReader) ReadAll(input io.Reader) (records [][]string, err error) {
	r := bufio.NewReader(input)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return records, errors.Wrapf(err, "could not read line %d from the CSV data %s", len(records)+1, input)
		}

		records = append(records, reader.fetchRecordsFromLine(line))
	}

	return records, err
}

// basic CSV file parser
func (reader CSVIOReader) fetchRecordsFromLine(line string) (result []string) {
	prefix := ""
	prefixIsOpen := false
	for _, char := range line {
		if char == ';' || char == '\n' || char == '\r' {
			continue
		}
		if char == '"' && !prefixIsOpen {
			prefixIsOpen = true
			continue
		}

		if char == '"' && prefixIsOpen {
			result = append(result, prefix)
			prefixIsOpen = false
			prefix = ""
			continue
		}

		prefix += string(char)
	}

	return result
}
