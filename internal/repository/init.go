package repository

import (
	"github.com/nichojovi/tauria-test/cmd/config"
	"github.com/nichojovi/tauria-test/internal/utils/database"
)

func NewUserRepository(db *database.Store, config *config.MainConfig) UserRepository {
	return &userRepo{
		db:  db,
		cfg: config,
	}
}

func NewRoomRepository(db *database.Store, config *config.MainConfig) RoomRepository {
	return &roomRepo{
		db:  db,
		cfg: config,
	}
}
