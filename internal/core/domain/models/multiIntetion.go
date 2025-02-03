package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type MultiIntentionInput struct {
	FromSystem string `json:"fromSystem,omitempty"`
	FromUser   string `json:"fromUser,omitempty"`
	ToSystem   string `json:"toSystem,omitempty"`
	ToUser     string `json:"toUser,omitempty"`
	Finish     bool   `json:"finish,omitempty"`
}

func (m *MultiIntentionInput) Json() string {
	var json, err = json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(json)
}

func (m *MultiIntentionInput) String() string {
	if m.ToUser != "" {
		return fmt.Sprintf("toUser(%s)", m.ToUser)
	}
	if m.ToSystem != "" {
		return fmt.Sprintf("toSystem(%s)", m.ToSystem)
	}
	return ""
}

func NewMultiIntentionInputFromString(text string) (*MultiIntentionInput, error) {
	var input = new(MultiIntentionInput)

	err := json.Unmarshal([]byte(jsonClear(text)), input)

	if err != nil {
		return nil, fmt.Errorf("error al parsear el mensaje: %w %s", err, text)
	}

	return input, err

}

func jsonClear(text string) string {
	text = strings.Replace(text, "```json", "", -1)
	text = strings.Replace(text, "```", "", -1)
	return text
}
