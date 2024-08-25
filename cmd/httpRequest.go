package cmd

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
)

func createRequest(method string, token string, url string) (*http.Request, error) {
	// Create a new HTTP request for GET request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Adding Bearer token to the request
	req.Header.Set("Authorization", "Bearer "+token)

	// Set the Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func createHttpClient() *http.Client {
	// Create a new HTTP client and make the request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return client
}

func getResponse(config Config, url string) (string, error) {
	logger := log.New(os.Stdout, "Log: ", log.Ldate|log.Ltime)

	req, err := createRequest("GET", config.OauthToken, url)
	if err != nil {
		logger.Printf("Error making HTTP request: %v\n", err)
		return "", err
	}

	resp, err := createHttpClient().Do(req)
	if err != nil {
		logger.Printf("Error making HTTP request: %v\n", err)
		return "", err
	}

	defer resp.Body.Close()

	// Check if the response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		logger.Printf("Error: received non-OK HTTP status: %s\n", resp.Status)
		return "", err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("Error reading response body: %v\n", err)
		return "", err
	}

	return string(body), nil
}

func downloadFile(config Config, url string, filePath string) (bool, error) {
	logger := log.New(os.Stdout, "Log:", log.Ldate|log.Ltime)

	req, err := createRequest("GET", config.OauthToken, url)
	if err != nil {
		logger.Printf("Error making HTTP request: %v\n", err)
	}

	resp, err := createHttpClient().Do(req)
	if err != nil {
		logger.Printf("Error making HTTP request: %v\n", err)
	}

	defer resp.Body.Close()

	// Check if the response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		logger.Printf("Error: received non-OK HTTP status: %s\n", resp.Status)
		return false, err
	}

	// Create the output file
	out, err := os.Create(filePath)
	if err != nil {
		logger.Printf("Error creating output file: %v\n", err)
		return false, err
	}
	defer out.Close()

	// Copy the response body to the output file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logger.Printf("Error saving file: %v\n", err)
		return false, err
	}

	return true, nil
}
