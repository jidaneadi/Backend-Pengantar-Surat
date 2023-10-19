package models

type Keterangan_Tidak_Mampu struct {
	ID          string `gorm:"primaryKey;autoIncrement" json:"id"`
	Id_surat    string `json:"id_surat"`
	Ket_dokter  string `json:"ket_dokter"`
	SPTJM       string `json:"sptjm"`
	KTP_pelapor string `json:"ktp_pelapor"`
	KTP_saksi   string `json:"ktp_saksi"`
	Surat       *Surat `gorm:"foreignKey:Id_surat;references:ID" json:"surat"`
}

func (Keterangan_Tidak_Mampu) TableName() string {
	return "ket_tidak_mampu"
}
