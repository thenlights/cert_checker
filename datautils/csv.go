package datautils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func ReadKnowledgeCsv(filename string, separator string) KnowledgeBase {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if separator == "" {
		reader.Comma = ','
	} else {
		reader.Comma = rune(separator[0])
	}

	// Read the header of the file to be used as map keys
	headers, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	data := make(KnowledgeBase)

	// Read all the other rows
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		// Create map  CSV header => row value
		row := make(map[string]string)
		for i, header := range headers {
			if i < len(record) {
				row[header] = record[i]
			}
		}
		data[row["domain"]] = row
	}

	return data
}

func ReadHostsCsv(filename string, separator string) IpToHostName {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	row := 0

	reader := csv.NewReader(file)
	if separator == "" {
		reader.Comma = ','
	} else {
		reader.Comma = rune(separator[0])
	}

	data := make(IpToHostName)

	// Read all the other rows
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		row += 1

		if len(record) == 2 {
			data[record[0]] = record[1]
		} else {
			log.Fatalf("Row: %i - Invalid record for ip to host in csv, skipping", row)
		}
	}

	return data
}
