package product

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound     = errors.New("no such product")
	ErrProductAlreadyExist = errors.New("such product already exists")
)

type Repository interface {
	GetAll() ([]Product, error)
	GetById(id uuid.UUID) (Product, error)
	Add(product Product) error
	Update(product Product) error
	Delete(id uuid.UUID) error
}
