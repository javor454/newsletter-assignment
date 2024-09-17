package domain

type Subscription struct {
	id                 *ID
	newsletterPublicID *ID
	email              *Email
}

func NewSubscription(newsletterPublicID *ID, email *Email) *Subscription {
	return &Subscription{
		id:                 NewID(),
		newsletterPublicID: newsletterPublicID,
		email:              email,
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
