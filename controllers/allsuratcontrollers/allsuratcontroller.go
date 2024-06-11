package allsuratcontrollers

import (
	"Backend_TA/config"
	"Backend_TA/models"
	"Backend_TA/utils"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
		updateTanggal := dataSurat.UpdatedAt.String()[0:10] + " " + dataSurat.UpdatedAt.String()[11:19]
		data[i] = fiber.Map{
			"id_surat":      dataSurat.ID,
			"id_masyarakat": dataSurat.Id_masyarakat,  //Id surat yg digenerate
			"nik":           dataSurat.Masyarakat.NIK, //NIK yg unik
			"nama":          dataSurat.Masyarakat.Nama,
			"jns_surat":     dataSurat.Jns_surat,
			"status":        dataSurat.Status,
			"updated_at":    updateTanggal,
			"keterangan":    dataSurat.Keterangan,
		}
	}
	return c.JSON(data)
}

func ShowSuratById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id surat kosong"})
	}

	var surat models.Surat
	if err := models.DB.First(&surat, "id =?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"msg": err.Error()})
	}
	return c.JSON(surat)
}

func ShowSuratByNik(c *fiber.Ctx) error {
	id := c.Params("id")
	tx := models.DB
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id user kosong"})
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
		updateTanggal := dataDoc.UpdatedAt.String()[0:10] + " " + dataDoc.UpdatedAt.String()[11:19]
		data[i] = fiber.Map{
			"id_surat":      dataDoc.ID,
			"id_masyarakat": dataDoc.Id_masyarakat,
			"nik":           dataDoc.Masyarakat.NIK,
			"nama":          dataDoc.Masyarakat.Nama,
			"jns_surat":     dataDoc.Jns_surat,
			"status":        dataDoc.Status,
			"updated_at":    updateTanggal,
			"keterangan":    dataDoc.Keterangan,
		}
	}

	return c.JSON(data)
}

func ShowDataDoc(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id surat kosong"})
	}

	tx := models.DB
	var doc_syarat []models.Dokumen_Syarat
	if err := tx.Preload("Surat").Joins("JOIN surat ON surat.id = dokumen_syarat.id_surat").Where("surat.id =?", id).Find(&doc_syarat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "Data tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
	}

	data := make([]fiber.Map, len(doc_syarat))
	for i, dataDoc := range doc_syarat {
		data[i] = fiber.Map{
			"id":        dataDoc.ID,              //ID dokumen syarat
			"filename":  dataDoc.Filename,        //Nama dokumen syarat
			"jns_surat": dataDoc.Surat.Jns_surat, //jenis surat yang diajukan
		}
	}

	return c.JSON(data)
}

// func ShowDocSyarat(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	if id == "" {
// 		return c.Status(400).JSON(fiber.Map{"msg": "Id kosong!"})
// 	}

// 	var doc_syarat models.Dokumen_Syarat
// 	if err := models.DB.Where("id =?", id).First(&doc_syarat).Error; err != nil {
// 		return c.Status(400).JSON(fiber.Map{"msg": err})
// 	}

// 	return c.Download(fmt.Sprintf("./public/ktp_baru/%s", doc_syarat.Filename))
// }

// Masih terdapat kesalahan, akan dibenahi saat tersambung ke frontend
func UpdateDocSyarat(c *fiber.Ctx) error {
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

func Update(c *fiber.Ctx) error {

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

	if err := models.DB.Where("id =?", id).Updates(&inputSurat).Error; err != nil {
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

func ShowDocSyarat(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id kosong!"})
	}

	var doc_syarat models.Dokumen_Syarat
	if err := models.DB.Where("id =?", id).First(&doc_syarat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	// var dataSurat models.S

	// Membuat objek ChaCha20
	dataKey := config.RenderEnv("KEY_CHACHA20")
	dataNonce := config.RenderEnv("NONCE")
	key, _ := hex.DecodeString(dataKey)
	nonce, _ := hex.DecodeString(dataNonce)

	var keyArray [32]byte
	var nonceArray [12]byte

	copy(keyArray[:], key)
	copy(nonceArray[:], nonce)

	chacha := &utils.ChaCha20{
		Key:     keyArray,
		Nonce:   nonceArray,
		Counter: 1,
	}

	// Ambil nama file input dari parameter
	inputFileName := doc_syarat.Filename
	if inputFileName == "" {
		return c.Status(http.StatusBadRequest).SendString("Nama file input tidak boleh kosong")
	}

	// Persiapkan jalur file input dan output
	inputFilePath := filepath.Join("./public/ktp_baru", inputFileName)
	outputFileName := strings.TrimSuffix(inputFileName, filepath.Ext(inputFileName)) + ".pdf"
	outputFilePath := filepath.Join(os.TempDir(), outputFileName)

	// Dekripsi file
	err := utils.DecryptFile(chacha, inputFilePath, outputFilePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal melakukan dekripsi: " + err.Error())
	}

	return c.Download(outputFilePath)
}
