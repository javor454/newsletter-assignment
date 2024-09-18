package firebase

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/logger"
)

// TODO (nice2have): prepare wrapper for readability in code

func NewClient(ctx context.Context, conf *config.FirebaseConfig, lg logger.Logger) (*db.Client, error) {

	// opt := option.WithCredentialsJSON(conf.ServiceAccountKey)
	// cf := &firebase.Config{DatabaseURL: fmt.Sprintf("%s:%d?ns=%s", conf.Host, conf.Port, conf.ProjectID)}
	cf := &firebase.Config{DatabaseURL: "http://0.0.0.0:9000/?ns=strv-go-newsletter-javor-jiri"}

	app, err := firebase.NewApp(ctx, cf)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %s", err.Error())
	}

	client, err := app.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase database: %s", err.Error())
	}

	// Test the connection using the built-in method
	lg.Info("Testing Firebase connection...")
	connRef := client.NewRef("http://127.0.0.1:9000/?ns=strv-go-newsletter-javor-jiri")
	var connected bool
	if err := connRef.Get(ctx, &connected); err != nil {
		return nil, fmt.Errorf("error checking connection: %s", err.Error())
	}
	if !connected {
		return nil, fmt.Errorf("not connected to Firebase")
	}
	lg.Info("Firebase reports connected")

	// Perform a write operation
	testRef := client.NewRef("test")
	testValue := fmt.Sprintf("test_value_%d", time.Now().UnixNano())
	lg.Info(fmt.Sprintf("Writing test value: %s", testValue))
	if err := testRef.Set(ctx, testValue); err != nil {
		return nil, fmt.Errorf("error writing test value: %s", err.Error())
	}

	// Read back the value
	var readValue string
	if err := testRef.Get(ctx, &readValue); err != nil {
		return nil, fmt.Errorf("error reading test value: %s", err.Error())
	}
	lg.Info(fmt.Sprintf("Read back test value: %s", readValue))

	if readValue != testValue {
		return nil, fmt.Errorf("value mismatch: wrote %s, read %s", testValue, readValue)
	}

	lg.Info("Successfully connected to Firebase emulator and verified read/write operations")

	return client, nil
}
