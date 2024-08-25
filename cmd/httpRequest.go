package cmd

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
)

func prepareHttpClient() *http.Client {
	// Create a new HTTP client and make the request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return client
}

func createRequest(config Config, url string) (*http.Request, error) {
	// Create a new HTTP request for GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Adding Bearer token to the request
	req.Header.Set("Authorization", "Bearer "+config.OauthToken)

	// Set the Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func getResponse(config Config, url string) (string, error) {

	client := prepareHttpClient()

	req, err := createRequest(config, url)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Check for response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func downloadFile(config Config, url string, filePath string) (bool, error) {
	client := prepareHttpClient()

	req, err := createRequest(config, url)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	// Check if the response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		return false, err
	}

	// Create the output file
	out, err := os.Create(filePath)
	if err != nil {
		return false, err
	}
	defer out.Close()

	// Copy the response body to the output file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return false, err
	}

	return true, nil
}
