package domain

type Newsletter struct {
	id          *ID
	name        string
	description *string
}

func NewNewsletter(name string, description *string) *Newsletter {
	return &Newsletter{
		id:          NewID(),
		name:        name,
		description: description,
	}
}

func CreateNewsletterFromExisting(id *ID, name string, description *string) *Newsletter {
	return &Newsletter{
		id:          id,
		name:        name,
		description: description,
	}
}

func (u *Newsletter) Id() *ID {
	return u.id
}

func (u *Newsletter) Name() string {
	return u.name
}

func (u *Newsletter) Description() *string {
	return u.description
}
