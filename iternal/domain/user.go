package domain

import "github.com/gofrs/uuid"

type User struct {
	Base
	Login    string `gorm:"column:login"`
	Password string `gorm:"column:password"`
	PermLVL  int    `gorm:"column:permLVL"`
}

func (u *User) ID() uuid.UUID {
	return u.Base.ID
}

func (u *User) Log() string {
	return u.Login
}

func (u *User) Permissions() int {
	return u.PermLVL
}

func (u *User) Pass() string {
	return u.Password
}
