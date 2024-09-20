package config

import (
	"github.com/spf13/viper"
)

const (
	envFirebaseUseEmulator            = "CONFIG_FIREBASE_USE_EMULATOR"
	envFirebaseEmulatorHost           = "CONFIG_FIREBASE_EMULATOR_HOST"
	envFirebaseDatabaseURL            = "CONFIG_FIREBASE_DATABASE_URL"
	envFirebaseServiceAccountFilePath = "CONFIG_FIREBASE_SERVICE_ACCOUNT_FILE_PATH"
)

type FirebaseConfig struct {
	UseEmulator            bool
	EmulatorHost           string
	DatabaseURL            string
	ServiceAccountFilePath string
}

func CreateFirebaseConfig() (*FirebaseConfig, error) {
	emulatorHost := viper.GetString(envFirebaseEmulatorHost)
	if emulatorHost == "" {
		return nil, getMissingError(envFirebaseEmulatorHost)
	}
	databaseURL := viper.GetString(envFirebaseDatabaseURL)
	if databaseURL == "" {
		return nil, getMissingError(envFirebaseDatabaseURL)
	}
	useEmulator := viper.GetBool(envFirebaseUseEmulator)
	serviceAccountFilePath := viper.GetString(envFirebaseServiceAccountFilePath)
	if serviceAccountFilePath == "" {
		return nil, getMissingError(envFirebaseServiceAccountFilePath)
	}

	return &FirebaseConfig{
		EmulatorHost:           emulatorHost,
		DatabaseURL:            databaseURL,
		UseEmulator:            useEmulator,
		ServiceAccountFilePath: serviceAccountFilePath,
	}, nil
}
