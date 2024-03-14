package models

import "encoding/json"

type SFPayload struct {
	Path  string          `json:"path"`
	Input json.RawMessage `json:"input"`
}
