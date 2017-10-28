package contact

import (
	"dolo-tracking-import/logger"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/govalidator"
)

// ParseFromCSV reads csv in ctx.Config.CSVFile and returns a list of Descriptor
func ParseFromCSV(reader *csv.Reader) ([]Descriptor, error) {
	var (
		err                error
		descriptorList     []Descriptor
		record             []string
		emailDoubleChecker map[string]bool
	)

	emailDoubleChecker = map[string]bool{}

	i := 0
	added := 0
	for {
		i++

		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return descriptorList, fmt.Errorf("Error on line %d %s", i, err.Error())
		}

		if len(record) < 3 {
			return descriptorList, fmt.Errorf("Invalid record on line %d", i)
		}

		email := strings.Trim(record[1], " ")
		if !govalidator.IsEmail(email) {
			logger.Warn(fmt.Sprintf("Email '%s' is not an email on line %d", email, i))
			continue
		}

		if emailDoubleChecker[email] {
			logger.Warn(fmt.Sprintf("Email '%s' was found in double on line: %d", email, i))
			continue
		}
		emailDoubleChecker[email] = true

		descriptorList = append(descriptorList, NewDescriptor(record[0], record[2], email))
		added++
	}

	logger.Debug(fmt.Sprintf("Parsed %d records out of %d", added, i))
	return descriptorList, nil
}