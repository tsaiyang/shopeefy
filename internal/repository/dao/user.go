package dao

import "context"

type UserDAO interface {
	Insert(ctx context.Context, user User) error
}

type User struct{}
