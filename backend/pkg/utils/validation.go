package utils

import (
	"regexp"
)

func IsValidEmail(email string) bool {
	// A simple regex for validating email addresses
	// This is a basic check and may not cover all edge cases
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched := regexp.MustCompile(emailRegex).MatchString(email)
	return matched
}
