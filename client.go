package client

import(
	"log"
	"net/http"
	"os"
	"time"
)

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

func New() *Client {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatal("No base url defined.")
		// baseURL = "https://api.com"
	}

	client := Client{
		BaseURL:    baseURL,
		Key:        os.Getenv("API_KEY"),
		Passphrase: os.Getenv("API_PASSPHRASE"),
		Secret:     os.Getenv("API_SECRET"),
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		RetryCount: 0,
	}

	if os.Getenv("API_SANDBOX") == "1" {
		client.UpdateConfig(&ClientConfig{
			BaseURL: "https://api.sandbox.com",
		})
	}

	return &client
}

func (c *Client) UpdateConfig(config *ClientConfig) {
	baseURL := config.BaseURL
	key := config.Key
	passphrase := config.Passphrase
	secret := config.Secret

	if baseURL != "" {
		c.BaseURL = baseURL
	}
	if key != "" {
		c.Key = key
	}
	if passphrase != "" {
		c.Passphrase = passphrase
	}
	if secret != "" {
		c.Secret = secret
	}
}
