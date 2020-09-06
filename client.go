package client

type Client struct {
	BaseURL string
	Secret string
	Key string
	Passphrase string
	HTTPClient *http.Client
	RetryCount int
}

type ClientConfig struct {
	BaseURL string
	Secret string
	Key string
	Passphrase string
}

