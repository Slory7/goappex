package datamodels

import (
	"time"
)

type UserLogins struct {
	Id        int64
	UserID    int64     `xorm:"'userid' notnull index"`
	LoginTime time.Time `xorm:"'logintime' notnull"`
}

func (u *UserLogins) TableName() string {
	return "user_logins"
}
