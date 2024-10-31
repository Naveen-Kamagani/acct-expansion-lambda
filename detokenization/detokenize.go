package detokenization

import (
	"acct-expansion-lambda/slog"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// MakeDetokenizeRequest sends a POST request to the specified URL with the provided payload.
func MakeDetokenizeRequest(url string, payload RequestPayload) (Response, error) {
	var response Response

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		slog.JSONLogger.With("payload", payload).Error("Error marshalling JSON payload", "error", err)
		return response, fmt.Errorf("error marshalling JSON payload: %w", err)
	}
	slog.JSONLogger.With("url", url).Info("Payload marshalled successfully")

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.JSONLogger.Error("Error creating HTTP request", "url", url, "error", err)
		return response, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("AUTH_TOKEN"))
	req.Header.Set("api-key", os.Getenv("API_KEY"))
	req.Header.Set("id-claim", os.Getenv("ID-CLAIM"))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.JSONLogger.With("url", url).Error("Error sending HTTP request", "error", err)
		return response, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	slog.JSONLogger.With("status_code", resp.StatusCode).Info("HTTP request sent successfully")

	// Read and decode the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.JSONLogger.Error("Error reading response body", "error", err)
		return response, fmt.Errorf("error reading response body: %w", err)
	}

	// Check for specific status codes
	switch resp.StatusCode {
	case 200:
		// Handle the 200 OK case if you need any specific logic
		slog.JSONLogger.With("status_code", resp.StatusCode).Info("Received 200 OK response")

	case 404:
		// Handle the 404 Not Found case
		slog.JSONLogger.With("status_code", resp.StatusCode).Warn("Resource not found (404)")
		return response, fmt.Errorf("resource not found (404) for URL: %s", url)

	case 500:
		// Handle the 500 Internal Server Error case
		slog.JSONLogger.With("status_code", resp.StatusCode).Error("Internal server error (500)")
		return response, fmt.Errorf("internal server error (500) for URL: %s", url)

	default:
		// Handle non-2xx response codes
		errMsg := fmt.Sprintf("Received non-2xx response code: %d", resp.StatusCode)
		logContext := slog.JSONLogger.With("status_code", resp.StatusCode)
		if len(body) > 0 { // Ensure there's a body to log
			logContext = logContext.With("response_body", string(body))
		}
		logContext.Warn(errMsg)
		return response, fmt.Errorf("%s, body: %s", errMsg, string(body))
	}

	// Unmarshal the response body into the Response struct
	if err := json.Unmarshal(body, &response); err != nil {
		slog.JSONLogger.With("response_body", string(body)).Error("Error unmarshalling response", "error", err)
		return response, fmt.Errorf("error unmarshalling response: %w", err)
	}

	slog.JSONLogger.With("response_body", string(body)).Info("Received API response successfully")
	return response, nil
}
