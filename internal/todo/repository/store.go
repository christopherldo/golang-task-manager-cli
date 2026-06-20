package repository

import "chrisldo.com/todo-cli/internal/todo/models"

type Store interface {
	ReadDatabase() ([]byte, error)
	WriteDatabase(tasks []models.Task) error
}
