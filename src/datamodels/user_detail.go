package datamodels

import "time"

type UserDetail struct {
	UserID        int64     `xorm:"'userid' pk notnull"`
	LastLoginTime time.Time `xorm:"'lastlogintime'"`
}

func (u *UserDetail) TableName() string {
	return "user_detail"
}
