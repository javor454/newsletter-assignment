package row

import (
	"time"
)

type MailType string

const (
	SubscriptionType MailType = "SUBSCRIPTION"
)

type Newsletter struct {
	ID          string
	PublicID    string
	Name        string
	Description *string
	CreatedAt   time.Time
}

type EmailJob struct {
	ID     string
	Type   MailType
	Params []byte
}
