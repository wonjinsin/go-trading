package model

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//Signup ...
type Signup struct {
	Email    string   `json:"email"`
	Password Password `json:"password"`
	Nick     *string  `json:"nick,omitempty"`
}

// Validate ...
func (su *Signup) Validate() bool {
	return su.Email != "" && !su.Password.IsEmpty() && su.Nick != nil
}

// Signin ...
type Signin struct {
	Email    string   `json:"email"`
	Password Password `json:"password"`
}

// CheckPassword ...
func (si *Signin) CheckPassword(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(si.Password))
	return err == nil
}

// Validate ...
func (si *Signin) Validate() bool {
	return si.Email != "" && si.Password != ""
}

func (si *Signin) String() string {
	return fmt.Sprintf("Email[%s] Password[%s]",
		si.Email,
		fmt.Sprintf("%c%s%c", si.Password[0], strings.Repeat("*", len(si.Password)-2), si.Password[len(si.Password)-1]),
	)
}
