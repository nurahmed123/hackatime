package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendSlackMessage(airtableAPIKey string, userID string, message string, blocksJSON string) error {
	record := map[string]interface{}{
		"fields": map[string]interface{}{
			"requester_identifier": "Hackatime Reset Password",
			"target_slack_id":      userID,
			"message_text":         message,
			"message_blocks":       blocksJSON,
			"unfurl_links":         true,
			"unfurl_media":         true,
			"send_success":         false,
		},
	}

	payload := map[string]interface{}{
		"records": []interface{}{record},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://middleman.hackclub.com/airtable/v0/appTeNFYcUiYfGcR6/arrpheus_message_requests", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+airtableAPIKey)
	req.Header.Set("User-Agent", "waka.hackclub.com (reset password)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	var result struct {
		Records []struct{} `json:"records"`
		Error   struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("error parsing response: %v", err)
	}

	if result.Error.Type != "" {
		return fmt.Errorf("Airtable error: %s - %s", result.Error.Type, result.Error.Message)
	}

	return nil
}
