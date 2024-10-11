package data

import (
	"database/sql"
	"errors"
)

var (
	ErrEditConflict   = errors.New("edit conflict")
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Movie MovieModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Movie: MovieModel{
			db: db,
		},
	}
}
