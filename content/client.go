package content

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Placeholder client configuration struct to pass client info into meta
type ClientConfig struct {
	ClientID     string
	ClientSecret string
	URL          string
	AuthUrl      string
	HTTPClient   *http.Client
}

type Client struct {
	url        string
	httpClient *http.Client
}

func NewClient(config *ClientConfig) (*Client, error) {

	if config.AuthUrl == "" {
		config.AuthUrl = "https://auth.adis.ws/oauth/token"
	}

	if config.URL == "" {
		config.URL = "https://api.amplience.net/v2/content"
	}

	auth := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     config.AuthUrl,
	}

	// If a custom httpClient is passed use that
	var httpClient *http.Client
	if config.HTTPClient != nil {
		httpClient = auth.Client(
			context.WithValue(oauth2.NoContext, oauth2.HTTPClient, config.HTTPClient))
	} else {
		httpClient = auth.Client(context.TODO())
	}

	client := &Client{
		url:        config.URL,
		httpClient: httpClient,
	}
	return client, nil
}

func (client *Client) request(method string, path string, body []byte, output interface{}) error {
	url := fmt.Sprintf("%s%s", client.url, path)
	log.Printf("%s: %s\n", method, url)

	buf := bytes.NewBuffer(body)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, method, url, buf)

	req.Header.Set("content-type", "application/json")

	logRequest(req)

	resp, err := client.httpClient.Do(req)

	if resp != nil {
		logResponse(resp)
	}

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(output)
	if err != nil {
		return err
	}

	return nil

}
