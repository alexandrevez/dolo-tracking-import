package main

import (
	"bufio"
	"dolo-tracking-import/appconfig"
	"dolo-tracking-import/contact"
	"dolo-tracking-import/context"
	"dolo-tracking-import/logger"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func newConfiguration(csvFile string, hubspotKey string) (*context.Configuration, error) {
	csvFileFullPath := appconfig.GetAppPath() + csvFile
	if _, err := os.Stat(csvFileFullPath); err != nil {
		return nil, err
	}

	return &context.Configuration{
		CSVFile: csvFileFullPath,
		Hubspot: context.HubspotConfig{
			APIKey: hubspotKey,
		},
	}, nil
}

func buildContext(csvFile string, hubspotKey string) (*context.App, error) {
	var (
		err    error
		config *context.Configuration
		ctx    *context.App
	)

	if config, err = newConfiguration(csvFile, hubspotKey); err != nil {
		return ctx, err
	}

	ctx = &context.App{
		Config: *config,
	}

	return ctx, nil
}

func main() {
	var (
		err         error
		ctx         *context.App
		file        *os.File
		contactList []contact.Descriptor
	)

	csvFile := flag.String("file", "", "CSV file. Comma separated list of contacts with format:\n\t<company_name>,<email>,<domain>")
	hubspotKey := flag.String("hubspot", "", "Hubspot API key")
	flag.Parse()

	fmt.Println(*csvFile, " ", *hubspotKey)

	if *csvFile == "" || *hubspotKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Build the app context
	if ctx, err = buildContext(*csvFile, *hubspotKey); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Parse the CSV file
	if file, err = os.Open(ctx.Config.CSVFile); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	if contactList, err = contact.ParseFromCSV(reader); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// TODO add records to Hubspot (if not exist)
	// -1 Add contact if does not exist
	// -2 Add company (and set filters)
	fmt.Println("TODO: shits", len(contactList))
}
