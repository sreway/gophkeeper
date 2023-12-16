package sqlite

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func errHandle(err error) error {
	var e sqlite3.Error

	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrNotFound
	}

	if errors.As(err, &e) {
		if e.Code == sqlite3.ErrConstraint {
			switch e.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				return models.ErrAlreadyExists
			default:
				return err
			}
		}
		return err
	}

	return err
}
