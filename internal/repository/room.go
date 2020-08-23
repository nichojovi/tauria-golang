package repository

import (
	"context"
	"database/sql"

	"github.com/nichojovi/tauria-test/internal/entity"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	insertRoomQuery         = "insert into room(name, host_user, participant, capacity) values (?, ?, ?, ?)"
	updateHostUserQuery     = "update room set host_user = ? where host_user = ? and id = ?"
	getRoomInfo             = "select id, name, host_user, participant, capacity from room"
	updateParticipantsQuery = "update room set participant = ? where name = ?"
)

func (rr *roomRepo) RegisterRoom(ctx context.Context, request entity.RoomDB) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.RegisterRoom")
	defer span.Finish()

	result, err := rr.db.GetMaster().ExecContext(ctx, insertRoomQuery,
		request.Name,
		request.HostUser,
		request.Participants,
		request.Capacity,
	)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func (rr *roomRepo) UpdateHost(ctx context.Context, oldHost, newHost string, roomID int64) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.UpdateHost")
	defer span.Finish()

	result, err := rr.db.GetMaster().ExecContext(ctx, updateHostUserQuery, newHost, oldHost, roomID)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func (rr *roomRepo) GetRoomInfoBasedOnName(ctx context.Context, roomName string) (*entity.RoomDB, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.GetRoomInfoBasedOnName")
	defer span.Finish()

	query := getRoomInfo + " where name = ?"

	result := new(entity.RoomDB)
	err := rr.db.GetSlave().GetContext(ctx, result, query, roomName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (rr *roomRepo) UpdateParticipants(ctx context.Context, participants []byte, roomName string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.UpdateParticipants")
	defer span.Finish()

	result, err := rr.db.GetMaster().ExecContext(ctx, updateParticipantsQuery, participants, roomName)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func (rr *roomRepo) GetRoomInfoBasedOnID(ctx context.Context, roomID int64) (*entity.RoomDB, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.GetRoomInfoBasedOnID")
	defer span.Finish()

	query := getRoomInfo + " where id = ?"

	result := new(entity.RoomDB)
	err := rr.db.GetSlave().GetContext(ctx, result, query, roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (rr *roomRepo) GetAllRoomsBasedOnUsername(ctx context.Context, username string) ([]entity.RoomDB, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.GetAllRoomsBasedOnUsername")
	defer span.Finish()

	query := getRoomInfo + " where JSON_SEARCH(participant, 'all', '" + username + "') > 1"

	var result []entity.RoomDB
	err := rr.db.GetSlave().SelectContext(ctx, &result, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}
