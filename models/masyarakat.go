package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Masyarakat struct {
	Idm          string    `gorm:"primaryKey;autoIncrement;column:idm" json:"idm"`
	NIK          string    `json:"nik" validate:"numeric,min=15"`
	Nama         string    `json:"nama"`
	No_hp        string    `json:"no_hp" validate:"numeric"`
	Gender       Gender    `gorm:"default:laki-laki" json:"gender"`
	Tempat_lahir string    `json:"tempat_lahir"`
	Birthday     string    `json:"birthday" validate:"datetime=2006-02-02"`
	Alamat       string    `json:"alamat"`
	CreatedAt    time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
	User         *User     `gorm:"foreignKey:NIK;references:ID" json:"user"`
}
type Gender string

const (
	Laki      Gender = "laki-laki"
	Perempuan Gender = "perempuan"
)

func ValidateMasyarakat(masyarakat *Masyarakat) error {
	validate := validator.New()
	return validate.Struct(masyarakat)
}

func (Masyarakat) TableName() string {
	return "masyarakat"
}

func (m *Masyarakat) BeforeCreate(tx *gorm.DB) (err error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	m.UpdateAt = time.Now()
	return nil
}
