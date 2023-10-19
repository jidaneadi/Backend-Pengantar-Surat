package models

import "time"

type Dokumen_Syarat struct {
	ID        string    `gorm:"primaryKey;autoIncrement" json:"id"`
	Id_surat  string    `json:"id_surat"`
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Surat     *Surat    `gorm:"foreignKey:Id_surat;references:ID" json:"surat"`
}

func (Dokumen_Syarat) TableName() string {
	return "dokumen_syarat"
}
