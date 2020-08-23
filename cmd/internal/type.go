package internal

import (
	"github.com/nichojovi/tauria-test/internal/service"
)

type (
	Service struct {
		User service.UserService
		Room service.RoomService
	}
)
