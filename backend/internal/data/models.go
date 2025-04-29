package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict raised")
)

type Models struct {
	Book            BookModel
	Tokens          TokenModel
	Permissions     PermissionModel
	Users           UserModel
	WritingVariants WritingVariantModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Book:            BookModel{DB: db},
		Tokens:          TokenModel{DB: db},
		Permissions:     PermissionModel{DB: db},
		Users:           UserModel{DB: db},
		WritingVariants: WritingVariantModel{DB: db},
	}
}
