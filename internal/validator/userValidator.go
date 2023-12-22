package validator

import "Meow/internal/data"

// Email validator function will check for empty and
// be in email format.
func ValidateEmail(v *Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(Matches(email, EmailRX), "email", "must be a valid email address")
}

// Check some limitations for user password
// length and empty.
func ValidatePasswordPlainText(v *Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

// Main user validator function
func ValidateUser(v *Validator, user *data.User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) < 500, "name", "must not be more than 500 character")

	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	ValidateEmail(v, user.Email)
	ValidatePasswordPlainText(v, *user.Password.Plaintext)

	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user).
	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
