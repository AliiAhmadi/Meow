// Mock types in this file will use for unit tests and access some sections
// of application without having an instance of application.

package data

import "time"

type MockMovieModel struct{}
type MockUserModel struct{}
type MockTokenModel struct{}
type MockPermissionModel struct{}

func (mock MockMovieModel) Insert(movie *Movie) error {
	return nil
}

func (mock MockMovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (mock MockMovieModel) Update(movie *Movie) error {
	return nil
}

func (mock MockMovieModel) Delete(id int64) error {
	return nil
}

func (mock MockMovieModel) GetAll(title string, genres []string, f Filters) ([]*Movie, Metadata, error) {
	return nil, Metadata{}, nil
}

// for userModel
func (mock MockUserModel) Insert(*User) error {
	return nil
}

func (mock MockUserModel) GetByEmail(string) (*User, error) {
	return nil, nil
}

func (mock MockUserModel) Update(*User) error {
	return nil
}

func (mock MockUserModel) GetForToken(scope string, plainToken string) (*User, error) {
	return nil, nil
}

// For tokenModel
func (mock MockTokenModel) DeleteAllForUser(scope string, userID int64) error {
	return nil
}

func (mock MockTokenModel) Insert(token *Token) error {
	return nil
}

func (mock MockTokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	return nil, nil
}

// For PermissionModel
func (mock MockPermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	return nil, nil
}
