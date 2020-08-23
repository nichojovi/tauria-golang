package repository

import (
	"context"

	"github.com/nichojovi/tauria-test/cmd/config"
	"github.com/nichojovi/tauria-test/internal/entity"
	"github.com/nichojovi/tauria-test/internal/utils/database"
)

type (
	userRepo struct {
		db  *database.Store
		cfg *config.MainConfig
	}
	roomRepo struct {
		db  *database.Store
		cfg *config.MainConfig
	}
)

type (
	UserRepository interface {
		GetAllUserInfo(ctx context.Context) ([]entity.User, error)
		GetUserInfo(ctx context.Context, username string) (*entity.User, error)
		RegisterUser(ctx context.Context, request entity.User) (int64, error)
		GetUserAuth(ctx context.Context, username, password string) (*entity.User, error)
		UpdatePassword(ctx context.Context, username, password string) (int64, error)
		DeleteUser(ctx context.Context, username string) (int64, error)
	}
	RoomRepository interface {
		RegisterRoom(ctx context.Context, request entity.RoomDB) (int64, error)
		UpdateHost(ctx context.Context, oldHost, newHost string, roomID int64) (int64, error)
		GetRoomInfoBasedOnName(ctx context.Context, roomName string) (*entity.RoomDB, error)
		UpdateParticipants(ctx context.Context, participants []byte, roomName string) (int64, error)
		GetRoomInfoBasedOnID(ctx context.Context, roomID int64) (*entity.RoomDB, error)
		GetAllRoomsBasedOnUsername(ctx context.Context, username string) ([]entity.RoomDB, error)
	}
)
