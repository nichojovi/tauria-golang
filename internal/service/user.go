package service

import (
	"context"

	"github.com/nichojovi/tauria-test/internal/entity"
	"github.com/opentracing/opentracing-go"
)

func (us *userService) GetAllUserInfo(ctx context.Context) ([]entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userService.GetAllUserInfo")
	defer span.Finish()

	var users []entity.User
	users, err := us.userRepo.GetAllUserInfo(ctx)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (us *userService) GetUserInfo(ctx context.Context, username string) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userService.GetUserInfo")
	defer span.Finish()

	var user *entity.User
	user, err := us.userRepo.GetUserInfo(ctx, username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (us *userService) RegisterUser(ctx context.Context, request entity.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userService.RegisterUser")
	defer span.Finish()

	_, err := us.userRepo.RegisterUser(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (us *userService) GetUserAuth(ctx context.Context, username, password string) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userService.GetUserAuth")
	defer span.Finish()

	var user *entity.User
	user, err := us.userRepo.GetUserAuth(ctx, username, password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (us *userService) UpdatePassword(ctx context.Context, username, password string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userService.UpdatePassword")
	defer span.Finish()

	_, err := us.userRepo.UpdatePassword(ctx, username, password)
	if err != nil {
		return err
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, username string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userService.DeleteUser")
	defer span.Finish()

	_, err := us.userRepo.DeleteUser(ctx, username)
	if err != nil {
		return err
	}

	return nil
}
