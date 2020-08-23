package api

import (
	"encoding/json"
	"net/http"

	"github.com/nichojovi/tauria-test/internal/entity"
	"github.com/nichojovi/tauria-test/internal/utils/auth"
	"github.com/nichojovi/tauria-test/internal/utils/encrypt"
	"github.com/nichojovi/tauria-test/internal/utils/response"
	"github.com/nichojovi/tauria-test/internal/utils/router"
	"github.com/opentracing/opentracing-go"
)

func (a *API) GetAllUsers(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "api.GetAllUsers")
	defer span.Finish()

	users, err := a.userService.GetAllUserInfo(ctx)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(err.Error())
	}

	return response.NewJSONResponse().SetData(users)
}

func (a *API) GetUserInfo(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "api.GetUserInfo")
	defer span.Finish()

	username := router.GetHttpParam(ctx, "username")
	if len(username) < entity.MinCharacter {
		return response.NewJSONResponse().SetError(response.ErrBadRequest)
	}

	user, err := a.userService.GetUserInfo(ctx, username)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(err.Error())
	}

	return response.NewJSONResponse().SetData(user)
}

func (a *API) RegisterUser(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "api.RegisterUser")
	defer span.Finish()

	var request entity.User
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage(err.Error())
	}

	if len(request.UserName) < entity.MinCharacter || len(request.Password) < entity.MinCharacter {
		return response.NewJSONResponse().SetError(response.ErrBadRequest)
	}

	err = a.userService.RegisterUser(ctx, request)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(err.Error())
	}

	return response.NewJSONResponse()
}

func (a *API) UpdatePassword(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "api.UpdatePassword")
	defer span.Finish()

	user := auth.GetAuthDetailFromContext(ctx)

	oldPassword := encrypt.SHA1(r.Header.Get("password"))
	newPassword := r.Header.Get("new_password")
	if oldPassword == newPassword {
		return response.NewJSONResponse().SetError(response.ErrBadRequest)
	}

	err := a.userService.UpdatePassword(ctx, user.UserName, newPassword)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(err.Error())
	}

	return response.NewJSONResponse()
}

func (a *API) DeleteUser(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "api.DeleteUser")
	defer span.Finish()

	user := auth.GetAuthDetailFromContext(ctx)
	err := a.userService.DeleteUser(ctx, user.UserName)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(err.Error())
	}

	return response.NewJSONResponse()
}
