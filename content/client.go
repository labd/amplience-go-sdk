package content

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// ClientConfig is used to initialize a new Client
type ClientConfig struct {
	ClientID     string
	ClientSecret string
	URL          string
	AuthURL      string
	HTTPClient   *http.Client
}

type Client struct {
	url        string
	httpClient *http.Client
	logLevel   int
}

// NewClient creates a new Client object
func NewClient(config *ClientConfig) (*Client, error) {

	if config.AuthURL == "" {
		config.AuthURL = "https://auth.adis.ws/oauth/token"
	}

	if config.URL == "" {
		config.URL = "https://api.amplience.net/v2/content"
	}

	auth := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     config.AuthURL,
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

	if os.Getenv("AMPLIENCE_DEBUG") != "" {
		client.logLevel = 1
	}
	return client, nil
}

func (client *Client) request(method string, path string, body []byte, output interface{}) error {

	raw_url, err := url.Parse(path)
	if err != nil {
		return err
	}

	var url string
	if raw_url.IsAbs() {
		url = path
	} else {
		url = fmt.Sprintf("%s%s", client.url, path)
	}

	buf := bytes.NewBuffer(body)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, method, url, buf)

	req.Header.Set("content-type", "application/json")

	if client.logLevel > 0 {
		logRequest(req)
	}

	resp, err := client.httpClient.Do(req)

	if resp != nil && client.logLevel > 0 {
		logResponse(resp)
	}

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 204:
		err = json.NewDecoder(resp.Body).Decode(output)
		if err != nil {
			return err
		}
	case resp.StatusCode == 204:
		return nil
	case resp.StatusCode >= 300:
		newErr := ErrorResponse{}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		if err = json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&newErr); err != nil {
			return err
		}
		// The API sometimes returns just `{message}` instead of `{errors: [{message}]}`,
		// so we try again for those cases.
		errorObject := ErrorObject{}
		if err = json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&errorObject); err != nil {
			return err
		}
		newErr.StatusCode = resp.StatusCode
		newErr.Inner = err
		newErr.Errors = []ErrorObject{errorObject}
		return &newErr
	}

	return nil

}
