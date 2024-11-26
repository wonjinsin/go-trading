package model

import (
	"magmar/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User ...
type User struct {
	UID      string   `json:"uid" gorm:"primaryKey"`
	Email    string   `json:"email"`
	Password Password `json:"password,omitempty"`
	Nick     *string  `json:"nick,omitempty"`
}

// NewUserBySignup ...
func NewUserBySignup(su *Signup) (user *User, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(su.Password), 14)
	if err != nil {
		return nil, err
	}

	return &User{
		UID:      uuid.New().String(),
		Email:    su.Email,
		Password: Password(bytes),
		Nick:     su.Nick,
	}, nil
}

// ValidateUpdateUser ...
func (u *User) ValidateUpdateUser() bool {
	if u.UID != "" {
		return false
	}
	return true
}

// UpdateUser ...
func (u *User) UpdateUser(user *User) *User {
	u.Email = user.Email
	u.Nick = user.Nick
	return u
}

// AfterFind ...
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	login, _ := ctx.Value(util.LoginKey).(bool)
	if !login {
		u.Password = ""
	}
	return
}
