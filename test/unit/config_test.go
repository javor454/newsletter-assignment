package unit

import (
	"testing"

	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	envSetFn       func()
	expectedErrMsg string
}

func Test_AppConfig_Success(t *testing.T) {
	initAppEnvVars()

	cf, err := config.NewAppConfig()
	assert.Nil(t, err)

	assert.Equal(t, 123, cf.HttpPort)
	assert.Equal(t, "http://localhost", cf.Host)
	assert.Equal(t, "jwt-secret", cf.JwtSecret)
	assert.Equal(t, []string{"http://localhost"}, cf.CorsAllowedOrigins)
	assert.Equal(t, []string{"authorization", "content-type"}, cf.CorsAllowedHeaders)
	assert.Equal(t, "Europe/Prague", cf.Timezone)
	assert.Equal(t, "sendgrid-api-key", cf.SendGridApiKey)
	assert.Equal(t, true, cf.SendMail)
}

func Test_FirebaseConfig_Success(t *testing.T) {
	initFirebaseEnvVars()

	cf, err := config.NewFirebaseConfig()
	assert.Nil(t, err)

	assert.Equal(t, true, cf.UseEmulator)
	assert.Equal(t, "firebase-emulator-host", cf.EmulatorHost)
	assert.Equal(t, "firebase-database-url", cf.DatabaseURL)
	assert.Equal(t, "firebase/service/account/file/path", cf.ServiceAccountFilePath)
}

func Test_PostgresConfig_Success(t *testing.T) {
	initPostgresEnvVars()

	cf, err := config.NewPostgresConfig()
	assert.Nil(t, err)

	assert.Equal(t, "pg-user", cf.User)
	assert.Equal(t, "pg-pass", cf.Password)
	assert.Equal(t, "pg-db/postgres", cf.Db)
	assert.Equal(t, "http://localhost", cf.Host)
	assert.Equal(t, 123, cf.Port)
	assert.Equal(t, "mig-dir", cf.MigrationsDir)
}

func Test_AppConfig_Fail(t *testing.T) {
	testCases := map[string]testCase{
		// App
		"http_port_zero": {
			envSetFn: func() {
				viper.Set("CONFIG_HTTP_PORT", 0)
			},
			expectedErrMsg: "missing required environment variable: CONFIG_HTTP_PORT",
		},
		"log_level_empty": {
			envSetFn: func() {
				viper.Set("CONFIG_LOG_LEVEL", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_LOG_LEVEL",
		},
		"jwt_secret_empty": {
			envSetFn: func() {
				viper.Set("CONFIG_JWT_SECRET", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_JWT_SECRET",
		},
		"cors_allowed_origins_empty": {
			envSetFn: func() {
				viper.Set("CONFIG_CORS_ALLOWED_ORIGINS", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_CORS_ALLOWED_ORIGINS",
		},
		"cors_allowed_headers_empty": {
			envSetFn: func() {
				viper.Set("CONFIG_CORS_ALLOWED_HEADERS", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_CORS_ALLOWED_HEADERS",
		},
		"timezone_empty": {
			envSetFn: func() {
				viper.Set("CONFIG_TIMEZONE", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_TIMEZONE",
		},
		"sendgrid_api_key_empty": {
			envSetFn: func() {
				viper.Set("CONFIG_SENDGRID_API_KEY", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_SENDGRID_API_KEY",
		},
		"host_empty": {
			envSetFn: func() {
				viper.Set("CONFIG_HOST", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_HOST",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			initAppEnvVars()
			initFirebaseEnvVars()

			tc.envSetFn()

			_, err := config.NewAppConfig()
			assert.Equal(t, tc.expectedErrMsg, err.Error())
		})
	}
}

func Test_FirebaseConfig_Fail(t *testing.T) {
	testCases := map[string]testCase{
		// Firebase
		"firebase_emulator_host": {
			envSetFn: func() {
				viper.Set("CONFIG_FIREBASE_EMULATOR_HOST", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_FIREBASE_EMULATOR_HOST",
		},
		"firebase_database_url": {
			envSetFn: func() {
				viper.Set("CONFIG_FIREBASE_DATABASE_URL", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_FIREBASE_DATABASE_URL",
		},
		"firebase_service_account_file_path": {
			envSetFn: func() {
				viper.Set("CONFIG_FIREBASE_SERVICE_ACCOUNT_FILE_PATH", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_FIREBASE_SERVICE_ACCOUNT_FILE_PATH",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			initFirebaseEnvVars()

			tc.envSetFn()

			_, err := config.NewFirebaseConfig()
			assert.Equal(t, tc.expectedErrMsg, err.Error())
		})
	}
}

func Test_PostgresConfig_Fail(t *testing.T) {
	testCases := map[string]testCase{
		// Firebase
		"postgres_user": {
			envSetFn: func() {
				viper.Set("CONFIG_POSTGRES_USER", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_POSTGRES_USER",
		},
		"postgres_password": {
			envSetFn: func() {
				viper.Set("CONFIG_POSTGRES_PASSWORD", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_POSTGRES_PASSWORD",
		},
		"postgres_db": {
			envSetFn: func() {
				viper.Set("CONFIG_POSTGRES_DB", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_POSTGRES_DB",
		},
		"postgres_host": {
			envSetFn: func() {
				viper.Set("CONFIG_POSTGRES_HOST", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_POSTGRES_HOST",
		},
		"postgres_port": {
			envSetFn: func() {
				viper.Set("CONFIG_POSTGRES_PORT", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_POSTGRES_PORT",
		},
		"postgres_migrations_dir": {
			envSetFn: func() {
				viper.Set("CONFIG_POSTGRES_MIGRATIONS_DIR", "")
			},
			expectedErrMsg: "missing required environment variable: CONFIG_POSTGRES_MIGRATIONS_DIR",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			initPostgresEnvVars()

			tc.envSetFn()

			_, err := config.NewPostgresConfig()
			assert.Equal(t, tc.expectedErrMsg, err.Error())
		})
	}
}

func initAppEnvVars() {
	viper.Set("CONFIG_HTTP_PORT", 123)
	viper.Set("CONFIG_LOG_LEVEL", "debug")
	viper.Set("CONFIG_JWT_SECRET", "jwt-secret")
	viper.Set("CONFIG_CORS_ALLOWED_ORIGINS", "http://localhost")
	viper.Set("CONFIG_CORS_ALLOWED_HEADERS", "authorization content-type")
	viper.Set("CONFIG_TIMEZONE", "Europe/Prague")
	viper.Set("CONFIG_SENDGRID_API_KEY", "sendgrid-api-key")
	viper.Set("CONFIG_SEND_MAIL", "true")
	viper.Set("CONFIG_HOST", "http://localhost")
}

func initFirebaseEnvVars() {
	viper.Set("CONFIG_FIREBASE_USE_EMULATOR", "true")
	viper.Set("CONFIG_FIREBASE_EMULATOR_HOST", "firebase-emulator-host")
	viper.Set("CONFIG_FIREBASE_DATABASE_URL", "firebase-database-url")
	viper.Set("CONFIG_FIREBASE_SERVICE_ACCOUNT_FILE_PATH", "firebase/service/account/file/path")
}

func initPostgresEnvVars() {
	viper.Set("CONFIG_POSTGRES_USER", "pg-user")
	viper.Set("CONFIG_POSTGRES_PASSWORD", "pg-pass")
	viper.Set("CONFIG_POSTGRES_DB", "pg-db/postgres")
	viper.Set("CONFIG_POSTGRES_HOST", "http://localhost")
	viper.Set("CONFIG_POSTGRES_PORT", 123)
	viper.Set("CONFIG_POSTGRES_MIGRATIONS_DIR", "mig-dir")
}
