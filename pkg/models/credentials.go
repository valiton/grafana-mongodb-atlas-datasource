package models

import (
	"encoding/json"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Settings struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
	ApiType string `json:"apiType"`
}

func LoadSettings(settings backend.DataSourceInstanceSettings) (*Settings, error) {
	s := &Settings{}
	if err := json.Unmarshal(settings.JSONData, &s); err != nil {
		return &Settings{}, err
	}

	if val, ok := settings.DecryptedSecureJSONData["privateKey"]; ok {
		s.PrivateKey = val
	}

	return s, nil
}
