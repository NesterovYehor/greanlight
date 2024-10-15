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
	User  UserModel
	Token TokenModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Movie: MovieModel{db: db},
		User:  UserModel{DB: db},
		Token: TokenModel{DB: db},
	}
}
