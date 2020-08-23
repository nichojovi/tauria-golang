package auth

import (
	"context"
	"net/http"

	"github.com/nichojovi/tauria-test/internal/entity"
	"github.com/nichojovi/tauria-test/internal/service"
	"github.com/nichojovi/tauria-test/internal/utils/encrypt"
	"github.com/nichojovi/tauria-test/internal/utils/response"
	"github.com/nichojovi/tauria-test/internal/utils/router"
	"github.com/opentracing/opentracing-go"
)

var (
	authPrefix = "Bearer "
)

type Opts struct {
	UserService service.UserService
}

type Module struct {
	UserService service.UserService
}

type AuthService interface {
	Authorize(h router.Handle) router.Handle
}

func New(o *Opts) *Module {
	return &Module{
		UserService: o.UserService,
	}
}

func getRequestHeader(r *http.Request) map[string]string {
	h := make(map[string]string, 0)
	for name := range r.Header {
		h[name] = r.Header.Get(name)
	}
	return h
}

func (m *Module) getUserInfoByRequest(ctx context.Context, r *http.Request) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "getUserInfoByRequest")
	defer span.Finish()

	username := r.Header.Get("user_name")
	password := encrypt.SHA1(r.Header.Get("password"))

	if len(username) < entity.MinCharacter || len(password) < entity.MinCharacter {
		return nil, response.ErrBadRequest
	}

	resp, err := m.UserService.GetUserAuth(ctx, username, password)
	if err != nil {
		return nil, response.ErrInternalServerError
	}

	if resp == nil {
		return nil, response.ErrForbidden
	}

	return resp, nil
}

func (m *Module) Authorize(h router.Handle) router.Handle {
	return func(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "Authorize")
		defer span.Finish()

		u, err := m.getUserInfoByRequest(ctx, r)
		if err != nil {
			return response.NewJSONResponse().SetError(err).SetLog("error", err)
		}

		header := getRequestHeader(r)
		ctx = context.WithValue(ctx, "RequestHeader", header)
		ctx = context.WithValue(ctx, "AuthDetail", u)
		r = r.WithContext(ctx)

		return h(w, r)
	}
}

func GetAuthDetailFromContext(ctx context.Context) *entity.User {
	v := ctx.Value("AuthDetail")
	if v == nil {
		return nil
	}
	return v.(*entity.User)
}
