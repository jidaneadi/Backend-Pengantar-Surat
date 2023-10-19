package models

import (
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        string `gorm:"primaryKey" json:"id" validate:"numeric,min=15"`
	Email     string `gorm:"unique" json:"email" validate:"email"`
	Password  string `json:"password" validate:"min=8"`
	Konf_pass string `json:"konf_pass" validate:"min=8"`
	Role      Role   `gorm:"default:masyarakat" json:"role"`
}

type NewPassword struct {
	Old_pass  string `json:"old_pass" validate:"min=8"`
	New_pass  string `json:"new_pass" validate:"min=8"`
	Konf_pass string `json:"konf_pass" validate:"min=8"`
}

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleMasyarakat Role = "masyarakat"
)

func ValidateUser(user *User) error {
	validate := validator.New()
	return validate.Struct(user)
}

func ValidatePass(password *NewPassword) error {
	validate := validator.New()
	return validate.Struct(password)
}

func (User) TableName() string {
	return "user"
}
