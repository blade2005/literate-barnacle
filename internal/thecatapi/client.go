package thecatapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURLV1 = "https://api.thecatapi.com/v1"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// type successResponse struct {
// 	Code int         `json:"code"`
// 	Data interface{} `json:"data"`
// }

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: BaseURLV1,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}

}

func (c *Client) sendRequest(req *http.Request, v *[]Image) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("x-api-key", c.apiKey)
	var err error = nil

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Try to unmarshall into errorResponse
	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	// Unmarshall and populate v
	// fullResponse := new(successResponse)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Printf("body: %s\n", body)

	// fullResponse := successResponse{
	// 	Data: v,
	// }

	// if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
	// 	return err
	// }

	err = json.Unmarshal(body, &v)

	return err
}
