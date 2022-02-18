package validator

import (
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3,30}$`)

func ValidateLoginData(username, email, password1, password2 string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	if err := validateUsername(username); err != nil {
		return err
	}

	if err := validatePassword(password1, password2); err != nil {
		return err
	}

	return nil
}

func validateUsername(u string) error {
	if !usernameRegex.MatchString(u) {
		return errors.New("username not valid")
	}

	return nil
}

func validateEmail(e string) error {
	if !emailRegex.MatchString(e) {
		return errors.New("email not valid")
	}

	return nil
}

func validatePassword(password1, password2 string) error {
	if password1 != password2 {
		return errors.New("passwords do not match")
	}

	if len(password1) < 6 {
		return errors.New("password too short")
	}
	num := `[0-9]{1}`
	az := `[a-z]{1}`
	AZ := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, password1); !b || err != nil {
		return errors.New("password needs at least one number")
	}
	if b, err := regexp.MatchString(az, password1); !b || err != nil {
		return errors.New("password needs at least one small character")
	}
	if b, err := regexp.MatchString(AZ, password1); !b || err != nil {
		return errors.New("password needs at leat one uppercase character")
	}
	if b, err := regexp.MatchString(symbol, password1); !b || err != nil {
		return errors.New("password needs at least on special symbol")
	}
	return nil
}
