package domain

type Subscription struct {
	id                 *ID
	newsletterPublicID *ID
	email              *Email
	token              string
}

func NewSubscription(newsletterPublicID *ID, email *Email, token string) *Subscription {
	return &Subscription{
		id:                 NewID(),
		newsletterPublicID: newsletterPublicID,
		email:              email,
		token:              token,
	}
}

func (s *Subscription) ID() *ID {
	return s.id
}

func (s *Subscription) NewsletterPublicID() *ID {
	return s.newsletterPublicID
}

func (s *Subscription) Email() *Email {
	return s.email
}

func (s *Subscription) Token() string {
	return s.token
}
