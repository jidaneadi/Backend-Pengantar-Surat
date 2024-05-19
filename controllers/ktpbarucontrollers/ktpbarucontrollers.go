package ktpbarucontrollers

import (
	"Backend_TA/models"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateKTPBaru(c *fiber.Ctx) error {

	id := c.Params("id")
	tx := models.DB
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id user kosong"})
	}

	//Cek data masyarakat apakah NIK sesuai atau tidak
	var masyarakat models.Masyarakat
	if err := models.DB.Where("nik =?", id).First(&masyarakat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "User not found"})
		}
		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
	}

	//Melakukan input surat
	var surat models.Surat
	surat.Id_masyarakat = masyarakat.Idm
	if err := c.BodyParser(&surat); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}
	surat.Jns_surat = "ktp_baru"
	surat.Status = "diproses"
	if surat.Id_masyarakat == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Id masyarakat kosong"})
	}

	//Menginputkan data dokumen syarat
	var doc_syarat models.Dokumen_Syarat

	//Input file dokumen probadi
	dokumenPribadi, err := c.FormFile("dokumen_pribadi")

	if dokumenPribadi == nil {
		return c.Status(400).JSON(fiber.Map{"msg": "File dokumen pribadi tidak boleh kosong!"})
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	errCheckContentTypeDP := checkContentType(dokumenPribadi, "application/pdf", "application/PDF")
	if errCheckContentTypeDP != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": errCheckContentTypeDP.Error(),
		})
	}

	nameDokumenPribadi := &dokumenPribadi.Filename
	extenstionFileDP := filepath.Ext(*nameDokumenPribadi)

	namaFile1 := masyarakat.NIK + surat.ID + "1" + extenstionFileDP

	//Menyimpan file dokumen pribadi pada database
	if err := c.SaveFile(dokumenPribadi, fmt.Sprintf("./public/%s/%s", surat.Jns_surat, namaFile1)); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Input file keterangan RT
	suratKeterangan, err := c.FormFile("keterangan_rt")

	if suratKeterangan == nil {
		return c.Status(400).JSON(fiber.Map{"msg": "File dokumen pribadi tidak boleh kosong!"})
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	errCheckContentType := checkContentType(suratKeterangan, "application/pdf", "application/PDF")
	if errCheckContentType != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": errCheckContentType.Error(),
		})
	}

	filename := &suratKeterangan.Filename
	extenstionFile := filepath.Ext(*filename)
	// newFilename := fmt.Sprintf("gambar-satu%s", extenstionFile)

	//Nama file ke 2
	namaFile2 := masyarakat.NIK + surat.ID + "2" + extenstionFile

	//Menyimpan file dokumen pribadi pada folder project
	if err := c.SaveFile(suratKeterangan, fmt.Sprintf("./public/%s/%s", surat.Jns_surat, namaFile2)); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Create data surat
	if err := tx.Create(&surat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}
	doc_syarat.Id_surat = surat.ID
	if err := c.BodyParser(&doc_syarat); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Create data dokumen surat
	doc_syarat.Filename = namaFile1
	if err := tx.Create(&doc_syarat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Create data dokumen surat
	doc_syarat.Filename = namaFile2
	doc_syarat.ID = ""
	if err := tx.Create(&doc_syarat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"msg": "Surat berhasil diajukan"})
}

func checkContentType(file *multipart.FileHeader, contentTypes ...string) error {
	if len(contentTypes) > 0 {
		for _, contentType := range contentTypes {
			contentTypeFile := file.Header.Get("Content-Type")
			if contentTypeFile == contentType {
				return nil
			}
		}

		return errors.New("not allowed file type")
	} else {
		return errors.New("not found content type to be checking")
	}
}
