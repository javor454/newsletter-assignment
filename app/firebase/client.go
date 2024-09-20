package firebase

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/javor454/newsletter-assignment/app/config"
	"google.golang.org/api/option"
)

type Client struct {
	*db.Client
}

func NewClient(ctx context.Context, conf *config.FirebaseConfig) (*Client, error) {
	if conf.UseEmulator {
		if err := os.Setenv("FIREBASE_DATABASE_EMULATOR_HOST", conf.EmulatorHost); err != nil {
			return nil, fmt.Errorf("error setting FIREBASE_DATABASE_EMULATOR_HOST: %v", err)
		}
	}

	opt := option.WithCredentialsFile(conf.ServiceAccountFilePath)
	cf := &firebase.Config{DatabaseURL: conf.DatabaseURL}

	app, err := firebase.NewApp(ctx, cf, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase database: %w", err)
	}

	return &Client{client}, nil
}
