package allsuratcontrollers

import (
	"Backend_TA/models"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ShowSurat(c *fiber.Ctx) error {
	var surat []models.Surat
	tx := models.DB
	if err := tx.Preload("Masyarakat").Find(&surat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
	}

	data := make([]fiber.Map, len(surat))
	for i, dataSurat := range surat {
		// newIdsurat := utils.EncryptHash(dataSurat.ID)
		createTanggal := dataSurat.CreatedAt.String()[0:9] + " " + dataSurat.CreatedAt.String()[12:19]
		updateTanggal := dataSurat.UpdatedAt.String()[0:9] + " " + dataSurat.UpdatedAt.String()[12:19]
		data[i] = fiber.Map{
			"id_surat":   dataSurat.ID,
			"nik":        dataSurat.Masyarakat.NIK,
			"nama":       dataSurat.Masyarakat.Nama,
			"jns_surat":  dataSurat.Jns_surat,
			"status":     dataSurat.Status,
			"created_at": createTanggal,
			"updated_at": updateTanggal,
			"keterangan": dataSurat.Keterangan,
		}
	}
	return c.JSON(data)
}

func ShowSuratByNik(c *fiber.Ctx) error {
	id := c.Params("id")
	tx := models.DB
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id surat kosong"})
	}

	var surat []models.Surat
	if err := tx.Preload("Masyarakat").Joins("JOIN masyarakat ON masyarakat.idm = surat.id_masyarakat").Where("masyarakat.nik =?", id).Find(&surat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "Data tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
	}
	data := make([]fiber.Map, len(surat))
	for i, dataDoc := range surat {
		createTanggal := dataDoc.CreatedAt.String()[0:9] + " " + dataDoc.CreatedAt.String()[12:19]
		updateTanggal := dataDoc.UpdatedAt.String()[0:9] + " " + dataDoc.UpdatedAt.String()[12:19]
		data[i] = fiber.Map{
			"id_surat":   dataDoc.ID,
			"nik":        dataDoc.Masyarakat.NIK,
			"nama":       dataDoc.Masyarakat.Nama,
			"jns_surat":  dataDoc.Jns_surat,
			"status":     dataDoc.Status,
			"created_at": createTanggal,
			"updated_at": updateTanggal,
			"keterangan": dataDoc.Keterangan,
		}
	}

	return c.JSON(data)
}

func ShowDocSyarat(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id surat kosong"})
	}

	tx := models.DB
	var doc_syarat []models.Dokumen_Syarat
	if err := tx.Where("id_surat =? ", id).Find(&doc_syarat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "Data tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
	}

	data := make([]fiber.Map, len(doc_syarat))
	for i, dataDoc := range doc_syarat {
		// newIdsurat := utils.EncryptHash(dataDoc.Id_surat)
		data[i] = fiber.Map{
			// "id_surat": newIdsurat[0:9],
			"syarat": dataDoc.Filename,
		}
	}

	return c.JSON(data)
}

// func Show(c *fiber.Ctx) error {
// 	var surat []models.Dokumen_Syarat

// 	//Join 3 tabel
// 	if err := models.DB.Preload("Surat.Masyarakat").Find(&surat).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return c.Status(404).JSON(fiber.Map{"msg": "Data null"})
// 		}
// 		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
// 	}

// 	data := make([]fiber.Map, len(surat))
// 	for i, DataSurat := range surat {
// 		newIdsurat := utils.EncryptHash(DataSurat.Id_surat)
// 		data[i] = fiber.Map{
// 			"id_surat":   newIdsurat[0:9],
// 			"nik":        DataSurat.Surat.Masyarakat.NIK,
// 			"nama":       DataSurat.Surat.Masyarakat.Nama,
// 			"syarat":     DataSurat.Filename,
// 			"jns_surat":  DataSurat.Surat.Jns_surat,
// 			"status":     DataSurat.Surat.Status,
// 			"updated_at": DataSurat.Surat.UpdatedAt.String()[0:10],
// 			"keterangan": DataSurat.Surat.Keterangan,
// 		}
// 	}
// 	return c.JSON(data)
// }

// func ShowId(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	var dataSurat []models.Dokumen_Syarat

// 	if err := models.DB.Preload("Surat.Masyarakat").
// 		Joins("JOIN surat ON surat.id = dokumen_syarat.id_surat").
// 		Joins("JOIN masyarakat ON masyarakat.idm = surat.id_masyarakat").
// 		Where("masyarakat.nik =?", id).
// 		Find(&dataSurat).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return c.Status(404).JSON(fiber.Map{"msg": "Data not found"})
// 		}
// 		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
// 	}

// 	data := make([]fiber.Map, len(dataSurat))
// 	for i, surat := range dataSurat {
// 		newIdsurat := utils.EncryptHash(surat.Id_surat)
// 		data[i] = fiber.Map{
// 			"id_surat":   newIdsurat[0:9],
// 			"nik":        surat.Surat.Masyarakat.NIK,
// 			"nama":       surat.Surat.Masyarakat.Nama,
// 			"syarat":     surat.Filename,
// 			"jns_surat":  surat.Surat.Jns_surat,
// 			"status":     surat.Surat.Status,
// 			"updated_at": surat.Surat.UpdatedAt.String()[0:10],
// 			"keterangan": surat.Surat.Keterangan,
// 		}
// 	}

// 	return c.JSON(data)
// }

func Update(c *fiber.Ctx) error {
	id := c.Params("id")
	tx := models.DB
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id surat kosong"})
	}

	var surat models.Surat
	if err := tx.Where("id = ?", id).First(&surat).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"msg": "Data surat tidak ditemukan"})
	}

	if err := c.BodyParser(&surat); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	if surat.Id_masyarakat == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id masyarakat kosong"})
	}

	if surat.Jns_surat == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Jenis surat kosong"})
	}

	if err := tx.Where("id = ? ", surat.ID).Updates(&surat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	var doc_syarat models.Dokumen_Syarat
	doc_syarat.Id_surat = surat.ID
	if err := c.BodyParser(&doc_syarat); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	if form, err := c.MultipartForm(); err == nil {
		//Pengambilan key input file
		files := form.File["filename"]

		//Pengambilan file lbh dr satu
		for i, file := range files {
			if cekFormat := strings.Split(file.Filename, "."); cekFormat[1] != "pdf" {
				return c.Status(400).JSON(fiber.Map{"msg": "File harus berekstensi PDF"})
			}

			var masyarakat models.Masyarakat
			if err := models.DB.Where("idm =? ", surat.Id_masyarakat).First(&masyarakat).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return c.Status(404).JSON(fiber.Map{"msg": "Data tidak ditemukan"})
				}
				return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
			}
			toString := strconv.Itoa(i)
			namaDokumen := masyarakat.NIK + surat.ID + toString + "-" + surat.Jns_surat + ".pdf"

			//Menyimpan surat
			if err := c.SaveFile(file, fmt.Sprintf("./public/%s/%s", surat.Jns_surat, namaDokumen)); err != nil {
				return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
			}

			doc_syarat.Filename = namaDokumen
			doc_syarat.ID = ""
			if err := tx.Where("id_surat = ? ", doc_syarat.Id_surat).Updates(&doc_syarat).Error; err != nil {
				return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
			}
		}
	}
	return c.JSON(fiber.Map{"msg": "Surat berhasil di ubah"})
}

func UpdateStatus(c *fiber.Ctx) error {

	id := c.Params("id")

	var surat models.Surat
	if err := models.DB.Where("id =?", id).First(&surat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "Data surat not found"})
		}
		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
	}

	var inputSurat models.Surat
	if err := c.BodyParser(&inputSurat); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	if surat.Id_masyarakat == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id masyarakat tidak boleh kosong"})
	}

	if surat.Jns_surat == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Jenis surat tidak boleh kosong"})
	}

	if surat.Status == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Status surat tidak boleh kosong"})
	}

	if err := models.DB.Where("id =?", id).Updates(&surat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"msg": "Status terupdate"})
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var surat models.Surat
	tx := models.DB
	if err := tx.Where("id =?", id).First(&surat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "Data not found"})
		}
		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
	}

	if err := tx.Where("id = ?", id).Delete(&surat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "Data not found"})
		}
		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
	}
	return c.JSON(fiber.Map{"msg": "Delete data sukses"})
}
