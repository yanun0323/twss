package repository

import (
	"context"
	"stocker/internal/domain"
	"stocker/internal/repository/mysql"
)

type Repo struct {
	mysql.MysqlDao
}

func New(ctx context.Context) domain.Repository {
	return &Repo{
		MysqlDao: mysql.New(ctx),
	}
}
