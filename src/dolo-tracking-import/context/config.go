package context

// Configuration defines the configuration of the application
type Configuration struct {
	Hubspot HubspotConfig `json:"hubspot"`
}

// HubspotConfig defines the configuration of the Hubspot API
type HubspotConfig struct {
	APIKey string `json:"api_key"`
}
