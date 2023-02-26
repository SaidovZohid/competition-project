package storage

import (
	"github.com/SaidovZohid/competition-project/storage/postgres"
	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
	Url() repo.UrlStorageI
}

type storagePg struct {
	userRepo repo.UserStorageI
	urlRepo  repo.UrlStorageI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &storagePg{
		userRepo: postgres.NewUser(db),
		urlRepo:  postgres.NewUrl(db),
	}
}

func (s *storagePg) User() repo.UserStorageI {
	return s.userRepo
}

func (s *storagePg) Url() repo.UrlStorageI {
	return s.urlRepo
}
