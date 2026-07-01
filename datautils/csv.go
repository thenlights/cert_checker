package datautils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func ReadKnowledgeCsv(filename string, separator string) (KnowledgeBase, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	if separator == "" {
		reader.Comma = ','
	} else {
		reader.Comma = rune(separator[0])
	}

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read knowledge csv header: %w", err)
	}

	data := make(KnowledgeBase)

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("read knowledge csv row: %w", err)
		}

		row := make(map[string]string)
		for i, header := range headers {
			if i < len(record) {
				row[header] = record[i]
			}
		}
		data[row["domain"]] = row
	}

	return data, nil
}

func ReadHostsCsv(filename string, separator string) (IpToHostName, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	if separator == "" {
		reader.Comma = ','
	} else {
		reader.Comma = rune(separator[0])
	}

	data := make(IpToHostName)

	row := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("read hosts csv row: %w", err)
		}

		row++
		if len(record) != 2 {
			return nil, fmt.Errorf("row %d: invalid ip-to-host record (expected 2 columns, got %d)", row, len(record))
		}
		data[record[0]] = record[1]
	}

	return data, nil
}
