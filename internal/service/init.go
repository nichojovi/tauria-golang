package service

import (
	"github.com/nichojovi/tauria-test/cmd/config"
	"github.com/nichojovi/tauria-test/internal/repository"
)

func NewUserService(user repository.UserRepository, cfg *config.MainConfig) UserService {
	return &userService{
		cfg:      cfg,
		userRepo: user,
	}
}

func NewRoomService(room repository.RoomRepository, cfg *config.MainConfig) RoomService {
	return &roomService{
		cfg:      cfg,
		roomRepo: room,
	}
}
