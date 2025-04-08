package model

import "time"

type Shop struct {
	Id          int64
	Name        string
	AccessToken string
	IsActive    bool
	Scope       string
	ExpireAt    time.Time
	UpdateAt    time.Time
	CreateAt    time.Time
}
