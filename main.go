package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"acct-expansion-lambda/detokenization"
	"acct-expansion-lambda/slog"

	ddlambda "github.com/DataDog/datadog-lambda-go"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//const region = os.Getenv("AWS_REGION")

const (
	detokenizePath = "/v1/data/sdm-protect/cloud-protegrity/unprotect"
	// Environment variable key for the API endpoint
	apiEndPointEnvVar = "API_ENDPOINT"
	dataElement       = "deACCOUNTNUM"
)

// CloudWatchEventData holds custom fields expected in the CloudWatch event detail
type CloudWatchEventData struct {
	AccountNumber string `json:"accountNumber"`
}

func init() {
	env := strings.ToLower(os.Getenv("ENV"))
	if env == "dev" || env == "test" {
		slog.SetDebugLevel()
	}

	slog.InitializeLoggers()
}

func handleRequest(_ context.Context, event events.CloudWatchEvent) error {
	slog.JSONLogger.With("event", event).Debug("Received following")

	// Get the API endpoint from environment variables
	apiEndPoint := os.Getenv(apiEndPointEnvVar)
	if apiEndPoint == "" {
		slog.JSONLogger.Error("API Endpoint is missing in the environment variables")
		return fmt.Errorf("API Endpoint is required in the environment variables")
	}

	// Parse the custom event data inside the "Detail" field
	var eventData CloudWatchEventData
	if err := json.Unmarshal(event.Detail, &eventData); err != nil {
		slog.JSONLogger.Error("Error unmarshalling event detail", "error", err)
		return fmt.Errorf("failed to unmarshal event detail: %w", err)
	}

	// Validate AccountNumber
	if eventData.AccountNumber == "" {
		slog.JSONLogger.Error("AccountNumber is missing in the event data")
		return fmt.Errorf("AccountNumber is required in the event data")
	}

	// Prepare the detokenization payload
	payload := detokenization.RequestPayload{
		DataElement: dataElement,
		Data:        []string{eventData.AccountNumber},
	}

	url := fmt.Sprintf("%s%s", apiEndPoint, detokenizePath)

	// Make the API call
	response, err := detokenization.MakeDetokenizeRequest(url, payload)
	if err != nil {
		slog.JSONLogger.Error("Error in detokenization request", "error", err)
		return fmt.Errorf("detokenization request failed: %w", err)
	}

	// Log the successful response
	slog.JSONLogger.With("response", response).Info("Detokenization request successful")
	return nil
}

func main() {
	lambda.Start(ddlambda.WrapFunction(handleRequest, nil))
}
