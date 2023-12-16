package models

import (
	"regexp"
)

type Email struct {
	value string
}

func (e *Email) String() string {
	return e.value
}

func NewEmail(value string) (*Email, error) {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, value)
	if !match {
		return nil, ErrInvalidEmailAddress
	}

	return &Email{
		value: value,
	}, nil
}
