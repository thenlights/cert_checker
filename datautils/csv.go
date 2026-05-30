package datautils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func ReadCsv(filename string, separator string) KnowledgeBase {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

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
			row[header] = record[i]
		}
		data[row["domain"]] = row
	}

	return data
}
