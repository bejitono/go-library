package client

import(
	"log"
	"net/http"
	"os"
	"math"
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

func (c *Client) Request(method string, url string, params, result interface{}) (res *http.Response, err error) {
	for i := 0; i < c.RetryCount+1; i++ {
		retryDuration := time.Duration((math.Pow(2, float64(i))-1)/2*1000) * time.Millisecond
		time.Sleep(retryDuration)

		res, err = c.request(method, url, params, result)
		if res != nil && res.StatusCode == 429 {
			continue
		} else {
			break
		}
	}
}

func (c *Client) request(method string, url string, params, result interface{}) (res *http.Response, err error) {
	var data []byte
	body := bytes.NewReader(make([]byte, 0))

	if params != nil {
		data, err = json.Marshal(params)
		if err != nil {
			return res, err
		}

		body = bytes.NewReader(data)
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, url)
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return res, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Client 1.0")

	h, error := c.Headers(method, url, timestamp, string(data))

	if err != nil {
		return res, err
	}

	for k, v := range h {
		req.Header.Add(k, v)
	}

	res, err = c.HTTPClient.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		defer res.Body.Close()
		error := Error{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&error); error != nil {
			return res, err
		}

		return res, error(error)
	}

	if result != nil {
		decoder := json.NewDecoder(res.Body)
		if err = decoder.Decode(result); err != nil {
			return res, err
		}
	}

	return res, nil
}

// Headers generates a map that can be used as headers to authenticate a request
func (c *Client) Headers(method, url, timestamp, data string) (map[string]string, error) {
	h := make(map[string]string)
	h["CB-ACCESS-KEY"] = c.Key
	h["CB-ACCESS-PASSPHRASE"] = c.Passphrase
	h["CB-ACCESS-TIMESTAMP"] = timestamp

	message := fmt.Sprintf(
		"%s%s%s%s",
		timestamp,
		method,
		url,
		data,
	)

	sig, err := generateSig(message, c.Secret)
	if err != nil {
		return nil, err
	}
	h["CB-ACCESS-SIGN"] = sig
	return h, nil
}
