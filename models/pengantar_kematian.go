package models

type Pengantar_Kematian struct {
	ID          string `gorm:"primaryKey;autoIncrement" json:"id"`
	Id_surat    string `json:"id_surat"`
	Ket_dokter  string `json:"ket_dokter"`
	SPTJM       string `json:"sptjm"`
	KTP_pelapor string `json:"ktp_pelapor"`
	KTP_saksi   string `json:"ktp_saksi"`
	Surat       *Surat `gorm:"foreignKey:Id_surat;references:ID" json:"surat"`
}

func (Pengantar_Kematian) TableName() string {
	return "ket_kematian"
}
