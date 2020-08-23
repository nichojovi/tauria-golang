package api

import (
	"net/http"

	"github.com/nichojovi/tauria-test/cmd/internal"
	"github.com/nichojovi/tauria-test/internal/service"
	"github.com/nichojovi/tauria-test/internal/utils/auth"
	"github.com/nichojovi/tauria-test/internal/utils/response"
	"github.com/nichojovi/tauria-test/internal/utils/router"
)

type Options struct {
	Prefix         string
	DefaultTimeout int
	AuthService    auth.AuthService
	Service        *internal.Service
}

type API struct {
	options     *Options
	authService auth.AuthService
	userService service.UserService
	roomService service.RoomService
}

func New(o *Options) *API {
	return &API{
		options:     o,
		authService: o.AuthService,
		userService: o.Service.User,
		roomService: o.Service.Room,
	}
}

func (a *API) Register() {
	r := router.New(&router.Options{Timeout: a.options.DefaultTimeout, Prefix: a.options.Prefix})

	// Testing
	r.GET("/ping", a.Ping)

	// User Management
	r.GET("/users", a.GetAllUsers)
	r.GET("/user/:username", a.GetUserInfo)
	r.POST("/register-user", a.RegisterUser)
	r.PUT("/update-password", a.authService.Authorize(a.UpdatePassword))
	r.DELETE("/delete-user", a.authService.Authorize(a.DeleteUser))

	// Room Management
	r.POST("/register-room", a.authService.Authorize(a.RegisterRoom))
	r.PUT("/change-host", a.authService.Authorize(a.UpdateHost))
	r.PUT("/join-status", a.authService.Authorize(a.UpdateJoinStatus))
	r.GET("/room/:id", a.GetRoomInfo)
	r.GET("/find-room/:username", a.GetAllRooms)
}

func (a *API) Ping(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	return response.NewJSONResponse().SetMessage("pong")
}
