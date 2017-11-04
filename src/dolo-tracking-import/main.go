package main

import (
	"bufio"
	"dolo-tracking-import/appconfig"
	"dolo-tracking-import/contact"
	"dolo-tracking-import/context"
	"dolo-tracking-import/hubspot"
	"dolo-tracking-import/logger"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
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

func processContactList(hubspotKey string, contactList []contact.Descriptor) error {
	var (
		err       error
		hsContact *hubspot.Contact
		hsCompany *hubspot.Company
	)

	for _, contact := range contactList {

		// Get or create contact
		if hsContact, err = hubspot.GetContact(hubspotKey, contact.Email); err != nil {
			return err
		}
		if hsContact == nil {
			logger.Debug(fmt.Sprintf("Creating new contact '%s'", contact.Email))
			if hsContact, err = hubspot.AddContact(hubspotKey, contact.Email); err != nil {
				return err
			}
		}

		// Get or create company
		if hsCompany, err = hubspot.GetCompany(hubspotKey, contact.DomainName, contact.CompanyName); err != nil {
			return err
		}
		if hsCompany == nil {
			logger.Debug(fmt.Sprintf("Creating company %s (%s)", contact.CompanyName, contact.DomainName))
			if hsCompany, err = hubspot.AddCompany(hubspotKey, contact.DomainName, contact.CompanyName); err != nil {
				return err
			}
		}

		// Ensure company has type radio and is associated with a company
		logger.Debug(fmt.Sprintf("Ensuring company %s (%s) has all properties", contact.CompanyName, contact.DomainName))
		if err = hubspot.UpdateCompany(hubspotKey, hsCompany.CompanyID); err != nil {
			return err
		}
		if err = hubspot.AddCompanyContact(hubspotKey, hsCompany.CompanyID, hsContact.ContactID); err != nil {
			return err
		}

		// Hubspot only allows 10 request per second. We are making 6, but you know..
		time.Sleep(time.Second)
	}

	return nil
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

	if err = processContactList(*hubspotKey, contactList); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
