package domain

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

type User struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Age       uint8  `validate:"required,gte=0,lte=130"`
	Email     string `validate:"required,email"`
}

func ValidateS(user *User) error {
	if err := validate.Struct(user); err != nil {
		return err
	}
	return nil
}

func (u *User) Validate() error {
	validate = validator.New()
	return validate.Struct(u)
}
