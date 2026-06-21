package repository

import (
	"chrisldo.com/todo-cli/internal/todo/db"
)

type Store interface {
	ReadDatabase() (db.DatabaseSchema, error)
	WriteDatabase(schema db.DatabaseSchema) error
}
