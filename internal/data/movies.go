package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"greenlight.nesty.net/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Year      int32     `json:"year"`
	Runtime   Runtime   `json:"runtime"`
	Genres    []string  `json:"genres"`
	Version   int64     `json:"version"`
}

type MovieModel struct {
	db *sql.DB
}

func (model *MovieModel) Insert(movie *Movie) error {
	query := `
        INSERT INTO movies (title, year, runtime, genres)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	return model.db.QueryRowContext(cntx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (model *MovieModel) Get(id int64) (*Movie, error) {
	query := `
        SELECT  id, created_at, title, year, runtime, genres, version
        FROM movies
        WHERE id = $1`

	var movie Movie

	cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := model.db.QueryRowContext(cntx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (model *MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	query := fmt.Sprintf(`
        SELECT id, created_at, title, year, runtime, genres, version
        FROM movies
        WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
        AND (genres @> $2 OR $2 = '{}')
        ORDER BY %s %s, id ASC
        LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}

	rows, err := model.db.QueryContext(cntx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	movies := []*Movie{}
	totalRecords := 1

	for rows.Next() {
		var movie Movie

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
			return nil, Metadata{}, err
		}

		movies = append(movies, &movie)
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
}

func (model *MovieModel) Update(movie *Movie) error {
	query := `
        UPDATE movies
        SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1 
        WHERE id = $5 AND version = $6
        RETURNING version
    `

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := model.db.QueryRowContext(cntx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (model *MovieModel) Delete(id int64) error {
	query := `
        DELETE FROM movies 
        WHERE id = $1
    `

	cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	result, err := model.db.ExecContext(cntx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func ValidateMovie(v *validator.Validator, input *Movie) {
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(input.Year != 0, "year", "year must be provided")
	v.Check(input.Year >= 1888, "year", "year must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(input.Runtime > 0, "runtime", "must be greater than 0")
	v.Check(input.Runtime != 0, "runtime", "must be provided")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must be greater than 1")
	v.Check(len(input.Genres) < 5, "genres", "must not contain more than 5")

	v.Check(validator.Unique(input.Genres), "genres", "must not contain dublicate values")
}
