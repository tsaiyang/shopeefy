package dao

import "time"

type AnnouncementBarActionSum struct {
	Id         int64 `gorm:"primaryKey,autoIncrement"`
	Shop       string
	BarId      int64
	ViewCount  int64
	ClickCount int64
	CloseCount int64
	Day        time.Time `gorm:"type:date;index"`
}
