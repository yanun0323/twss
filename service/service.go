package service

import "main/domain"

type Service struct {
	Repo domain.IRepository
}

func NewService(repo domain.IRepository) Service {
	return Service{Repo: repo}
}
