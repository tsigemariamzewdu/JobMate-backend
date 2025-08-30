package interfaces

// PhoneValidator defines phone normalization and validation

type IPhoneValidator interface {
	Normalize(phone string) (string, error)
	Validate(phone string) error
}

type IEmailValidator interface{
	Validate (email string ) error
}
