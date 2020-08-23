package repository

import (
	"context"
	"database/sql"

	"github.com/nichojovi/tauria-test/internal/entity"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	getAllUserInfoQuery = "select user_name, password, full_name, email, phone from user"
	insertUserQuery     = "insert into user(user_name, password, full_name, email, phone) values (?, SHA1(?), ?, ?, ?)"
	updatePasswordQuery = "update user set password = SHA1(?) where user_name = ?"
	deleteUserQuery     = "delete from user where user_name = ?"
)

func (ur *userRepo) GetAllUserInfo(ctx context.Context) ([]entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepo.GetAllUserInfo")
	defer span.Finish()

	var result []entity.User
	err := ur.db.GetSlave().SelectContext(ctx, &result, getAllUserInfoQuery)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ur *userRepo) GetUserInfo(ctx context.Context, username string) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepo.GetUserInfo")
	defer span.Finish()

	query := getAllUserInfoQuery + " where user_name = ?"

	result := new(entity.User)
	err := ur.db.GetSlave().GetContext(ctx, result, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (ur *userRepo) RegisterUser(ctx context.Context, request entity.User) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepo.RegisterUser")
	defer span.Finish()

	result, err := ur.db.GetMaster().ExecContext(ctx, insertUserQuery,
		request.UserName,
		request.Password,
		request.FullName,
		request.Email,
		request.Phone,
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

func (ur *userRepo) GetUserAuth(ctx context.Context, username, password string) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepo.GetUserAuth")
	defer span.Finish()

	query := getAllUserInfoQuery + " where user_name = ? and password = ?"

	result := new(entity.User)
	err := ur.db.GetSlave().GetContext(ctx, result, query, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (ur *userRepo) UpdatePassword(ctx context.Context, username, password string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepo.UpdatePassword")
	defer span.Finish()

	result, err := ur.db.GetMaster().ExecContext(ctx, updatePasswordQuery, password, username)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func (ur *userRepo) DeleteUser(ctx context.Context, username string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepo.DeleteUser")
	defer span.Finish()

	result, err := ur.db.GetMaster().ExecContext(ctx, deleteUserQuery, username)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}
