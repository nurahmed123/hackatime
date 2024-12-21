package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	conf "github.com/hackclub/hackatime/config"
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

	fmt.Println(string(body))

	conf.Log().Debug("Parsing heartbeat array")

	// Try bulk first
	var heartbeats []*models.Heartbeat
	if err := json.Unmarshal(body, &heartbeats); err == nil {
		return heartbeats, nil
	} else {
		err = fmt.Errorf("failed to parse heartbeat array: %v", err)
	}

	// Try single if bulk fails
	var heartbeat models.Heartbeat
	if err := json.Unmarshal(body, &heartbeat); err == nil {
		return []*models.Heartbeat{&heartbeat}, nil
	} else {
		err = fmt.Errorf("failed to parse heartbeat: %v", err)
	}

	return nil, err
}
