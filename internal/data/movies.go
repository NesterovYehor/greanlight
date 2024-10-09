package data

import (
	"time"

	"greenlight.nesty.net/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Year      int32     `json:"year"`
	Runtime   Runtime   `json:"runtime"`
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, input *Movie) {
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(input.Year != 0, "year", "year must be provided")
	v.Check(input.Year >= 1888, "year", "year must be greater than 1888")
	v.Check(input.Year >= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(input.Runtime > 0, "runtime", "must be greater than 0")
	v.Check(input.Runtime != 0, "runtime", "must be provided")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must be greater than 1")
	v.Check(len(input.Genres) < 5, "genres", "must not contain more than 5")

	v.Check(validator.Unique(input.Genres), "genres", "must not contain dublicate values")
}
