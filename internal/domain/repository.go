package domain

import "stocker/internal/model"

type Repository interface {
	DBRepository
}

type DBRepository interface {
	GetLastRaw() (model.Raw, error)
	InsertRaw(raw model.Raw) error
}
