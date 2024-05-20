package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Masyarakat struct {
	Idm          string    `gorm:"primaryKey;autoIncrement;column:idm" json:"idm"`
	NIK          string    `json:"nik" validate:"numeric,min=15"`
	No_Kk        string    `json:"no_kk" validate:"numeric,min=15"`
	Nama         string    `json:"nama"`
	No_hp        string    `json:"no_hp" validate:"numeric"`
	Gender       Gender    `gorm:"default:laki-laki" json:"gender"`
	Tempat_lahir string    `json:"tempat_lahir"`
	Birthday     string    `json:"birthday" validate:"datetime=2006-02-02"`
	Agama        Agama     `grom:"default:islam" json:"agama"`
	Pekerjaan    string    `json:"pekerjaan"`
	Status       Status    `gorm:"default:belum kawin" json:"status"`
	Alamat       string    `json:"alamat"`
	Warga_Negara string    `json:"warga_negara"`
	CreatedAt    time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
	User         *User     `gorm:"foreignKey:NIK;references:ID" json:"user"`
}
type Gender string

type Agama string

type Status string

const (
	Laki      Gender = "laki-laki"
	Perempuan Gender = "perempuan"
	Islam     Agama  = "islam"
	Kristen   Agama  = "kristen"
	Katholik  Agama  = "katholik"
	Hindu     Agama  = "hindu"
	Budha     Agama  = "budha"
	Khonghucu Agama  = "khonghucu"
	Kawin     Status = "kawin"
	Belum     Status = "belum kawin"
	Chidup    Status = "cerai hidup"
	Cmati     Status = "cerai mati"
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
