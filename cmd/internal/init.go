package internal

import (
	"github.com/nichojovi/tauria-test/cmd/config"
	"github.com/nichojovi/tauria-test/internal/repository"
	"github.com/nichojovi/tauria-test/internal/service"
	"github.com/nichojovi/tauria-test/internal/utils/database"
)

func GetService(db *database.Store, config *config.MainConfig) *Service {
	//REPO
	userRepository := repository.NewUserRepository(db, config)
	roomRepository := repository.NewRoomRepository(db, config)

	//SERVICE
	userService := service.NewUserService(userRepository, config)
	roomService := service.NewRoomService(roomRepository, config)

	return &Service{
		User: userService,
		Room: roomService,
	}
}
