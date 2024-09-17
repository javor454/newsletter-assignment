package domain

import "time"

type Newsletter struct {
	id          *ID
	publicID    *ID
	name        string
	description *string
	createdAt   time.Time
}

func NewNewsletter(name string, description *string) *Newsletter {
	return &Newsletter{
		id:          NewID(),
		publicID:    NewID(),
		name:        name,
		description: description,
		createdAt:   time.Now(),
	}
}

func CreateNewsletterFromExisting(id, publicID *ID, name string, description *string, createdAt time.Time) *Newsletter {
	return &Newsletter{
		id:          id,
		publicID:    publicID,
		name:        name,
		description: description,
		createdAt:   createdAt,
	}
}

func (u *Newsletter) ID() *ID {
	return u.id
}

func (u *Newsletter) PublicID() *ID {
	return u.publicID
}

func (u *Newsletter) Name() string {
	return u.name
}

func (u *Newsletter) Description() *string {
	return u.description
}

func (u *Newsletter) CreatedAt() time.Time {
	return u.createdAt
}
