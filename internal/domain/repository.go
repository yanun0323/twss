package domain

import "stocker/internal/model"

type Repository interface {
	DBRepository
}

type DBRepository interface {
	GetLastRaw() (model.DailyRaw, error)
	InsertRaw(raw model.DailyRaw) error
}
