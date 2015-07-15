package main

import (
	"encoding/csv"
	"errors"
	"io"
	"math/rand"
	"os"
	"strings"
)

var (
	missingEmailField = errors.New("Email field missing in header.")
)

func readCSV(path string, randDesc []RandTempDesc) (*[]Recipient, *string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var (
		header     []string
		headerRead bool
		emailField string
		recipients []Recipient
	)

	reader := csv.NewReader(file)
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}

		if headerRead {
			recipient := make(Recipient)

			for i, key := range header {
				recipient[key] = fields[i]
				for _, randColum := range randDesc {
					if strings.ToLower(key) == strings.ToLower(randColum.Name) {
						recipient[key] = randColum.Items[rand.Intn(len(randColum.Items))]
					}
					break
				}
			}

			recipients = append(recipients, recipient)
		} else {
			header = fields

			for _, v := range header {
				if strings.ToLower(v) == "email" {
					emailField = v
				}
			}

			if emailField == "" {
				return nil, nil, missingEmailField
			}

			headerRead = true
		}
	}

	return &recipients, &emailField, nil
}
