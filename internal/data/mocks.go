// Mock types in this file will use for unit tests and access some sections
// of application without having an instance of application.

package data

type MockMovieModel struct{}

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
