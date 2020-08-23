package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nichojovi/tauria-test/internal/entity"
	"github.com/nichojovi/tauria-test/internal/utils/auth"
	"github.com/nichojovi/tauria-test/internal/utils/response"
	"github.com/nichojovi/tauria-test/internal/utils/router"
	"github.com/opentracing/opentracing-go"
)

func (a *API) RegisterRoom(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "api.RegisterRoom")
	defer span.Finish()

	var request entity.Room
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage(err.Error())
	}

	user := auth.GetAuthDetailFromContext(ctx)
	roomRequest := entity.RoomDB{
		Name:         request.Name,
		HostUser:     user.UserName,
		Participants: []byte(request.Participants),
		Capacity:     request.Capacity,
	}

	err = a.roomService.RegisterRoom(ctx, roomRequest)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(err.Error())
	}

	return response.NewJSONResponse()
}

func (a *API) UpdateHost(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "api.UpdateHost")
	defer span.Finish()

	newHost := r.FormValue("new_host")
	roomID, err := strconv.ParseInt(r.FormValue("room_id"), 10, 64)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage(err.Error())
	}
	user := auth.GetAuthDetailFromContext(ctx)

	err = a.roomService.UpdateHost(ctx, user.UserName, newHost, roomID)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(err.Error())
	}

	return response.NewJSONResponse()
}

func (a *API) UpdateJoinStatus(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "api.UpdateJoinStatus")
	defer span.Finish()

	// join = true, leave = false
	req := &struct {
		Join     bool   `json:"join"`
		RoomName string `json:"room_name"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage(err.Error())
	}

	user := auth.GetAuthDetailFromContext(ctx)

	err = a.roomService.UpdateJoinStatus(ctx, user.UserName, req.RoomName, req.Join)
	if err != nil {
		if err == response.ErrInvalidRequest {
			return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage(err.Error())
		}
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(err.Error())
	}

	return response.NewJSONResponse()
}

func (a *API) GetRoomInfo(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "api.GetRoomInfo")
	defer span.Finish()

	roomID, err := strconv.ParseInt(router.GetHttpParam(ctx, "id"), 10, 64)
	if err != nil || roomID < 1 {
		return response.NewJSONResponse().SetError(response.ErrBadRequest)
	}

	room, err := a.roomService.GetRoomInfo(ctx, roomID)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(err.Error())
	}

	return response.NewJSONResponse().SetData(room)
}

func (a *API) GetAllRooms(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "api.GetAllRooms")
	defer span.Finish()

	username := router.GetHttpParam(ctx, "username")
	if len(username) < entity.MinCharacter {
		return response.NewJSONResponse().SetError(response.ErrBadRequest)
	}

	rooms, err := a.roomService.GetAllRoomsBasedOnUsername(ctx, username)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(err.Error())
	}

	return response.NewJSONResponse().SetData(rooms)
}
