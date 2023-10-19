package models

type Pengantar_Kelahiran struct {
	ID            string `gorm:"primaryKey;autoIncrement" json:"id"`
	Id_surat      string `json:"id_surat"`
	Ket_kelahiran string `json:"ket_kelahiran"`
	Dokumen_lain  string `json:"dokumen_lain"`
	Surat         *Surat `gorm:"foreignKey:Id_surat;references:ID" json:"surat"`
}

func (Pengantar_Kelahiran) TableName() string {
	return "ket_kelahiran"
}
