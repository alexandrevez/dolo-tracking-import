package hubspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// CompanyTypeRadio is the comapny type used for radio tracking
const CompanyTypeRadio = "RADIO"

// Company represents a company in Hubspot
type Company struct {
	CompanyID int
	Domain    string
	Name      string
	Type      string
}

// SearchRequest contains the properties as represented in the Hubspot API
type SearchRequest struct {
	Limit   int                   `json:"limit"`
	Options CompanyRequestOptions `json:"requestOptions"`
	Offset  CompanyOffset         `json:"offset"`
}

// CompanyRequestOptions is the list of parameters you need
type CompanyRequestOptions struct {
	Properties []string `json:"properties"`
}

// CompanyOffset will offset the result from CompanyID
// CompanyID can be found in a request result
type CompanyOffset struct {
	IsPrimary bool `json:"isPrimary"`
	CompanyID int  `json:"companyId"`
}

// SearchResponse is the response we receive from hubspot
type SearchResponse struct {
	Results []CompanyResponseResult `json:"results"`
	HasMore bool                    `json:"hasMore"`
	Offset  CompanyOffset           `json:"offset"`
}

// CompanyResponseResult is an actual search result
type CompanyResponseResult struct {
	CompanyID  int                             `json:"companyId"`
	Properties CompanyResponseResultProperties `json:"properties"`
}

// CompanyResponseResultProperties properties are wrapped... thanks hubspot
type CompanyResponseResultProperties struct {
	Name SearchProperty `json:"name"`
	Type SearchProperty `json:"type"`
}

// SearchProperty is a search result property. May contain other information, but not used for now
type SearchProperty struct {
	Value string `json:"value"`
}

// GetCompany will try to find a company with a domain name and name matching
// if it has more than 1 result, it will return the first one
func GetCompany(apiKey string, domain string, name string) (*Company, error) {
	const (
		hubspotURL = "https://api.hubapi.com/companies/v2/domains/%s/companies?hapikey=%s"
	)
	var (
		err           error
		url           string
		httpClient    http.Client
		searchReq     SearchRequest
		req           *http.Request
		resp          SearchResponse
		respRaw       *http.Response
		respBodyBytes []byte
		reqBodyBytes  []byte
	)

	httpClient = http.Client{}
	url = fmt.Sprintf(hubspotURL, domain, apiKey)

	searchReq = SearchRequest{
		Limit: 100,
		Offset: CompanyOffset{
			IsPrimary: true,
			CompanyID: 0,
		},
		Options: CompanyRequestOptions{
			Properties: []string{
				"domain",
				"name",
				"type",
			},
		},
	}

	for {
		if reqBodyBytes, err = json.Marshal(searchReq); err != nil {
			return nil, err
		}
		if req, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBodyBytes)); err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		if respRaw, err = httpClient.Do(req); err != nil {
			return nil, err
		}
		defer respRaw.Body.Close()

		if respBodyBytes, err = ioutil.ReadAll(respRaw.Body); err != nil {
			return nil, err
		}

		if respRaw.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error: %s \n%s", respRaw.Status, string(respBodyBytes))
		}

		if err = json.Unmarshal(respBodyBytes, &resp); err != nil {
			return nil, err
		}

		for _, companyResult := range resp.Results {
			if companyResult.Properties.Name.Value == name {
				return &Company{
					CompanyID: companyResult.CompanyID,
					Name:      name,
					Domain:    domain,
					Type:      companyResult.Properties.Type.Value,
				}, nil
			}
		}

		if !resp.HasMore {
			break
		}

		searchReq.Offset = resp.Offset
	}
	return nil, nil
}

// AddCompanyRequest is the request used to add a company
type AddCompanyRequest struct {
	Properties []Property `json:"properties"`
}

// AddCompany will add a company in hubspot
func AddCompany(apiKey string, domain string, name string) (*Company, error) {
	const (
		hubspotURL = "https://api.hubapi.com/companies/v2/companies?hapikey=%s"
	)
	var (
		err        error
		url        string
		req        AddCompanyRequest
		reqBytes   []byte
		httpClient http.Client
		httpReq    *http.Request
		respRaw    *http.Response
		bodyBytes  []byte
		resp       CompanyResponseResult
	)

	req = AddCompanyRequest{
		Properties: []Property{
			Property{
				Name:  "name",
				Value: name,
			},
			Property{
				Name:  "domain",
				Value: domain,
			},
			Property{
				Name:  "type",
				Value: CompanyTypeRadio,
			},
		},
	}

	if reqBytes, err = json.Marshal(req); err != nil {
		return nil, err
	}
	url = fmt.Sprintf(hubspotURL, apiKey)

	if httpReq, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBytes)); err != nil {
		return nil, err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient = http.Client{}
	if respRaw, err = httpClient.Do(httpReq); err != nil {
		return nil, err
	}
	defer respRaw.Body.Close()

	if bodyBytes, err = ioutil.ReadAll(respRaw.Body); err != nil {
		return nil, err
	}

	if respRaw.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %s \n%s", respRaw.Status, string(bodyBytes))
	}

	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, err
	}

	return &Company{
		CompanyID: resp.CompanyID,
		Name:      name,
		Domain:    domain,
		Type:      resp.Properties.Type.Value,
	}, nil
}

// UpdateCompany will set the company type to radio and contactId to contact list
func UpdateCompany(apiKey string, companyID int) error {

	const (
		hubspotURL = "https://api.hubapi.com/companies/v2/companies/%d?hapikey=%s"
	)
	var (
		err        error
		url        string
		req        AddCompanyRequest
		reqBytes   []byte
		httpClient http.Client
		httpReq    *http.Request
		respRaw    *http.Response
		bodyBytes  []byte
	)

	req = AddCompanyRequest{
		Properties: []Property{
			Property{
				Name:  "type",
				Value: CompanyTypeRadio,
			},
		},
	}

	if reqBytes, err = json.Marshal(req); err != nil {
		return err
	}
	url = fmt.Sprintf(hubspotURL, companyID, apiKey)

	if httpReq, err = http.NewRequest("PUT", url, bytes.NewBuffer(reqBytes)); err != nil {
		return err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient = http.Client{}
	if respRaw, err = httpClient.Do(httpReq); err != nil {
		return err
	}
	defer respRaw.Body.Close()

	if respRaw.StatusCode != http.StatusOK {
		if bodyBytes, err = ioutil.ReadAll(respRaw.Body); err != nil {
			return err
		}

		return fmt.Errorf("Error: %s \n%s", respRaw.Status, string(bodyBytes))
	}

	return nil
}

// AddCompanyContact will add a contact to a company
func AddCompanyContact(apiKey string, companyID int, contactID int) error {
	const (
		hubspotURL = "https://api.hubapi.com/companies/v2/companies/%d/contacts/%d?hapikey=%s"
	)

	var (
		err        error
		url        string
		httpClient http.Client
		httpReq    *http.Request
		respRaw    *http.Response
		bodyBytes  []byte
	)

	url = fmt.Sprintf(hubspotURL, companyID, contactID, apiKey)

	if httpReq, err = http.NewRequest("PUT", url, bytes.NewBuffer([]byte{})); err != nil {
		return err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient = http.Client{}
	if respRaw, err = httpClient.Do(httpReq); err != nil {
		return err
	}
	defer respRaw.Body.Close()

	if respRaw.StatusCode != http.StatusOK {
		if bodyBytes, err = ioutil.ReadAll(respRaw.Body); err != nil {
			return err
		}

		return fmt.Errorf("Error: %s \n%s", respRaw.Status, string(bodyBytes))
	}

	return nil
}
