package row

import "time"

type Newsletter struct {
	ID          string
	PublicID    string
	Name        string
	Description *string
	CreatedAt   time.Time
}
