package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"`                       // Unique integer ID for movies
	CreatedAt time.Time `json:"-"`                        // Timestamp for creation of a movie
	Title     string    `json:"title"`                    // String title for movie
	Year      int32     `json:"year,omitempty"`           // Movie release year
	Runtime   Runtime   `json:"runtime,omitempty,string"` // Movie time in minutes
	Genres    []string  `json:"genres,omitempty"`         // Slice of genres for the movie
	Version   int32     `json:"version"`                  // The version number starts at 1 and will be incremented each time the movie information is updated
}

var (
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	INSERT_QUERY = `INSERT INTO movies (title, year, runtime, genres) 
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, version
	`

	// Define the SQL query for retrieving the movie data.
	GET_QUERY = `SELECT id, created_at, title, year, runtime, genres, version
	FROM movies WHERE id = $1
	`

	// Declare the SQL query for updating the record and returning the new version
	// number.
	UPDATE_QUERY = `UPDATE movies SET title = $1, year = $2, runtime = $3,
	genres = $4, version = version + 1 WHERE version = $6 AND id = $5 RETURNING version
	`

	// Delete sql query for delete a record.
	DELETE_QUERY = `DELETE FROM movies WHERE id = $1`

	// Define a query for get multiple movie based on query parameters and sort.
	GET_ALL_QUERY = `SELECT id, created_at, title, year, runtime, genres, version
	FROM movies WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') AND 
	(genres @> $2 OR $2 = '{}')
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4
	`
)

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (movieModel MovieModel) Insert(movie *Movie) error {
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct.
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
	}

	// Create a context function for 3 second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Set returned values from database into movie instance or return any error if exists.
	return movieModel.DB.QueryRowContext(ctx, INSERT_QUERY, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (movieModel MovieModel) Get(id int64) (*Movie, error) {
	// The PostgreSQL bigserial type that we're using for the movie ID starts
	// auto-incrementing at 1 by default, so we know that no movies will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Declare an empty movie struct to hold the data returned by the query.
	var movie Movie

	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 3-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct.
	err := movieModel.DB.QueryRowContext(ctx, GET_QUERY, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}

	// At this point every things ok. so return movie struct.
	return &movie, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (movieModel MovieModel) Update(movie *Movie) error {

	// Create an slice of interfaces include parameters we want to pass QueryRow() function in follow.
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	// Define context.Context and cancel for set 3-second timeout for UPDATE_QUERY.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute Update query and scan new version that returned from database to movie.Version
	err := movieModel.DB.QueryRowContext(ctx, UPDATE_QUERY, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	// Ok.
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (movieModel MovieModel) Delete(id int64) error {
	// Return an error if the movie id less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	// Declare a ctx (context) for define a 3-second timeout for Delete query.
	// Also get a cancel function for cancel query when time riched 3 second.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := movieModel.DB.ExecContext(ctx, DELETE_QUERY, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it.
	if rows == 0 {
		return ErrRecordNotFound
	}

	// Ok.
	return nil
}

// Define GetAll() nethod on MovieModel for get all movies based on query parameters.
func (movieModel MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {
	// Create formatted query by placing order by parameters.
	query := fmt.Sprintf(GET_ALL_QUERY, filters.sortColumn(), filters.sortDirection())

	// Create a context with 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// As our SQL query now has quite a few placeholder parameters, let's collect the
	// values for the placeholders in a slice. Notice here how we call the limit() and
	// offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}

	// Using QueryContext() for execut query.
	// Result is sql.Rows.
	rows, err := movieModel.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Initialize an empty slice to hold the movie data.
	movies := []*Movie{}

	// Using rows.Next() to iterate through the rows in the result of query
	for rows.Next() {
		// Initialize an empty movie to cunstruct a new movie
		var movie Movie

		// Scan values from rows to movie instance.
		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)

		if err != nil {
			return nil, err
		}

		// Add new movie to movies slice.
		movies = append(movies, &movie)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Ok.
	return movies, nil
}
