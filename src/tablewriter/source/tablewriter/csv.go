package tablewriter

import (
	"encoding/csv"
	"io"
	"os"
)

func NewCSV(writer io.Writer, fileName string, hashHeader bool) (*Table, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return &Table{}, err
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	t, err := NewCSVReader(writer, csvReader, hashHeader)
	return t, err
}

func NewCSVReader(writer io.Writer, csvReader *csv.Reader, hashHeader bool) (*Table, error) {
	t := NewWriter(writer)
	if hashHeader {
		headers, err := csvReader.Read()
		if err != nil {
			return &Table{}, err
		}
		t.SetHeader(headers)
	}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return &Table{}, err
		}
		t.Append(record)
	}
	return t, nil
}
