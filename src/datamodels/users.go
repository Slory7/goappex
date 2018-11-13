package datamodels

import (
	"time"
)

type User struct {
	Id          int64     //`xorm:"pk notnull autoincr"`
	UserName    string    `xorm:"'username' varchar(50) notnull unique index"`
	Email       string    `xorm:"'email' varchar(100) notnull unique index"`
	IsSuper     bool      `xorm:"'issuper' notnull"`
	Age         uint8     `xorm:"'age' notnull default(0)"`
	Memo        string    `xorm:"'memo' varchar(200)"`
	DeletedAt   time.Time `xorm:"'deleted_at' deleted"`
	CreatedDate time.Time `xorm:"'created_date' created"`
	UpdatedDate time.Time `xorm:"'updated_date' updated"`
}

func (u *User) TableName() string {
	return "users"
}
