package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hackclub/hackatime/models"
)

func ParseHeartbeats(r *http.Request) ([]*models.Heartbeat, error) {
	// Read body once
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Try bulk first
	var heartbeats []*models.Heartbeat
	if err := json.Unmarshal(body, &heartbeats); err == nil {
		return heartbeats, nil
	}

	// Try single if bulk fails
	var heartbeat models.Heartbeat
	if err := json.Unmarshal(body, &heartbeat); err == nil {
		return []*models.Heartbeat{&heartbeat}, nil
	}

	return nil, fmt.Errorf("failed to parse heartbeat data: %v", err)
}
