package models

type Pengantar_KTP struct {
	ID             string `gorm:"primaryKey;autoIncrement" json:"id"`
	Id_surat       string `json:"id_surat"`
	Dokumen_syarat string `json:"dokumen_syarat"`
	Surat          *Surat `gorm:"foreignKey:Id_surat;references:ID" json:"surat"`
}

func (Pengantar_KTP) TableName() string {
	return "pengantar_ktp"
}
