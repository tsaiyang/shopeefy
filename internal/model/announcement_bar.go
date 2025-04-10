package model

import "time"

type AnnouncementBar struct {
	Id         int64
	Shop       string
	ConfigAll  string
	ConfigPart string
	Status     AnnouncementBarStatus
	UpdateAt   time.Time
	CreateAt   time.Time
}

type AnnouncementBarStatus uint8

const (
	AnnouncementBarStatusUnknown AnnouncementBarStatus = iota
	AnnouncementBarStatusPublished
	AnnouncementBarStatusUnpublished
	AnnouncementBarStatusDeleted
)

func (status AnnouncementBarStatus) ToUint8() uint8 {
	return uint8(status)
}
