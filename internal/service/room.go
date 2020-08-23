package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nichojovi/tauria-test/internal/entity"
	"github.com/nichojovi/tauria-test/internal/utils/response"
	"github.com/opentracing/opentracing-go"
)

func (rs *roomService) RegisterRoom(ctx context.Context, request entity.RoomDB) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomService.RegisterRoom")
	defer span.Finish()

	_, err := rs.roomRepo.RegisterRoom(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (rs *roomService) UpdateHost(ctx context.Context, oldHost, newHost string, roomID int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomService.UpdateHost")
	defer span.Finish()

	_, err := rs.roomRepo.UpdateHost(ctx, oldHost, newHost, roomID)
	if err != nil {
		return err
	}

	return nil
}

func (rs *roomService) UpdateJoinStatus(ctx context.Context, username, roomName string, join bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomService.UpdateHost")
	defer span.Finish()

	room, err := rs.roomRepo.GetRoomInfoBasedOnName(ctx, roomName)
	if err != nil {
		return err
	}

	if room == nil {
		log.Printf("[UpdateJoinStatus] Invalid room, roomName: %s", roomName)
		return response.ErrInvalidRequest
	}

	var listParticipant map[string]string
	var exist bool
	err = json.Unmarshal(room.Participants, &listParticipant)
	if err != nil {
		return err
	}

	// Check already join or not
	for i := 1; i <= len(listParticipant); i++ {
		if listParticipant[username] != "" {
			exist = true
		}
	}

	if join {
		// Check room capacity
		if room.Capacity < int64(len(listParticipant)+1) {
			log.Printf("[UpdateJoinStatus] Room is full, roomName: %s", roomName)
			return response.ErrInvalidRequest
		}
		if exist {
			log.Printf("[UpdateJoinStatus] User already join, roomName: %s, userName: %s", roomName, username)
			return response.ErrInvalidRequest
		}
		if username == room.HostUser || len(listParticipant) < 1 {
			listParticipant[username] = "host"
		} else {
			listParticipant[username] = "user"
		}
	} else {
		// Check user exist or not in participant
		if !exist {
			log.Printf("[UpdateJoinStatus] User already leave before, roomName: %s, userName: %s", roomName, username)
			return response.ErrInvalidRequest
		}
		delete(listParticipant, username)
	}

	// Update host if participant < 0
	if len(listParticipant) < 1 {
		_, err = rs.roomRepo.UpdateHost(ctx, room.HostUser, "", room.ID)
		if err != nil {
			return err
		}
	}

	if len(listParticipant) > 0 && room.HostUser == "" {
		_, err = rs.roomRepo.UpdateHost(ctx, "", username, room.ID)
		if err != nil {
			return err
		}
	}

	// Update participant list
	participants, err := json.Marshal(listParticipant)
	if err != nil {
		return err
	}

	_, err = rs.roomRepo.UpdateParticipants(ctx, participants, roomName)
	if err != nil {
		return err
	}

	return nil
}

func (rs *roomService) GetRoomInfo(ctx context.Context, roomID int64) (*entity.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomService.GetRoomInfo")
	defer span.Finish()

	roomDB, err := rs.roomRepo.GetRoomInfoBasedOnID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	room := &entity.Room{
		ID:           roomDB.ID,
		Name:         roomDB.Name,
		HostUser:     roomDB.HostUser,
		Participants: string(roomDB.Participants),
		Capacity:     roomDB.Capacity,
	}

	return room, nil
}

func (rs *roomService) GetAllRoomsBasedOnUsername(ctx context.Context, username string) ([]entity.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomService.GetRoomInfo")
	defer span.Finish()

	roomDB, err := rs.roomRepo.GetAllRoomsBasedOnUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	rooms := make([]entity.Room, len(roomDB))
	for i := 0; i < len(roomDB); i++ {
		rooms[i] = entity.Room{
			ID:           roomDB[i].ID,
			Name:         roomDB[i].Name,
			HostUser:     roomDB[i].HostUser,
			Participants: string(roomDB[i].Participants),
			Capacity:     roomDB[i].Capacity,
		}
	}

	return rooms, nil
}
