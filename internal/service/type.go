package service

import (
	"context"

	"github.com/nichojovi/tauria-test/cmd/config"
	"github.com/nichojovi/tauria-test/internal/entity"
	"github.com/nichojovi/tauria-test/internal/repository"
)

type (
	userService struct {
		cfg      *config.MainConfig
		userRepo repository.UserRepository
	}
	roomService struct {
		cfg      *config.MainConfig
		roomRepo repository.RoomRepository
	}
)

type (
	UserService interface {
		GetAllUserInfo(ctx context.Context) ([]entity.User, error)
		GetUserInfo(ctx context.Context, username string) (*entity.User, error)
		RegisterUser(ctx context.Context, request entity.User) error
		GetUserAuth(ctx context.Context, username, password string) (*entity.User, error)
		UpdatePassword(ctx context.Context, username, password string) error
		DeleteUser(ctx context.Context, username string) error
	}
	RoomService interface {
		RegisterRoom(ctx context.Context, request entity.RoomDB) error
		UpdateHost(ctx context.Context, oldHost, newHost string, roomID int64) error
		UpdateJoinStatus(ctx context.Context, username, roomName string, join bool) error
		GetRoomInfo(ctx context.Context, roomID int64) (*entity.Room, error)
		GetAllRoomsBasedOnUsername(ctx context.Context, username string) ([]entity.Room, error)
	}
)
