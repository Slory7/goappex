package datamodels

type UserWithMemo struct {
	Id   int64
	Memo string
}

func (u *UserWithMemo) TableName() string {
	return "users"
}

type UserAll struct {
	User
	UserDetail
	UserLogins []UserLogins
}

type UserWithDetail struct {
	User       `xorm:"extends"`
	UserDetail `xorm:"extends"`
}

func (u *UserWithDetail) TableName() string {
	return "users"
}
