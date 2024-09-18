package config

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/viper"
)

const (
	envFirebaseHost              = "CONFIG_FIREBASE_HOST"
	envFirebasePort              = "CONFIG_FIREBASE_PORT"
	envFirebaseProjectID         = "CONFIG_FIREBASE_PROJECT_ID"
	envFirebaseServiceAccountKey = "CONFIG_FIREBASE_SERVICE_ACCOUNT_KEY"
)

type FirebaseConfig struct {
	Host              string
	Port              int
	ProjectID         string
	ServiceAccountKey []byte
}

func CreateFirebaseConfig() (*FirebaseConfig, error) {
	host := viper.GetString(envFirebaseHost)
	if host == "" {
		return nil, getMissingError(envFirebaseHost)
	}
	port := viper.GetInt(envFirebasePort)
	if port == 0 {
		return nil, getMissingError(envFirebasePort)
	}
	projectID := viper.GetString(envFirebaseProjectID)
	if projectID == "" {
		return nil, getMissingError(envFirebaseProjectID)
	}
	serviceAccountKeyB64 := viper.GetString(envFirebaseServiceAccountKey)
	if serviceAccountKeyB64 == "" {
		return nil, getMissingError(envFirebaseServiceAccountKey)
	}
	serviceAccountKey, err := base64.StdEncoding.DecodeString(serviceAccountKeyB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode service account key: %s", err.Error())
	}

	return &FirebaseConfig{
		Host:              host,
		Port:              port,
		ProjectID:         projectID,
		ServiceAccountKey: serviceAccountKey,
	}, nil
}
