package ktpbarucontrollers

import (
	"Backend_TA/models"
	"fmt"
	"strings"

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

	//Create data surat
	if err := tx.Create(&surat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Menginputkan data dokumen syarat
	var doc_syarat models.Dokumen_Syarat

	doc_syarat.Id_surat = surat.ID
	if err := c.BodyParser(&doc_syarat); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Input file dokumen probadi
	dokumenPribadi, err := c.FormFile("dokumen_pribadi")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Cek format
	cekDokumenPribadi := strings.Split(dokumenPribadi.Filename, ".")
	if cekDokumenPribadi[1] != "pdf" {
		return c.Status(400).JSON(fiber.Map{"msg": "File harus berekstensi PDF"})
	}

	//Input file keterangan RT
	suratKeterangan, err := c.FormFile("keterangan_rt")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Cek format
	cekSuratKeterangan := strings.Split(suratKeterangan.Filename, ".")
	if cekSuratKeterangan[1] != "pdf" {
		return c.Status(400).JSON(fiber.Map{"msg": "File harus berekstemsi PDF"})
	}

	var namaFile = masyarakat.NIK + surat.ID + "1" + "-" + surat.Jns_surat + ".pdf"

	//Menyimpan file dokumen pribadi pada database
	if err := c.SaveFile(dokumenPribadi, fmt.Sprintf("./public/%s/%s", surat.Jns_surat, namaFile)); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Create data dokumen surat
	doc_syarat.Filename = namaFile
	if err := tx.Create(&doc_syarat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//
	namaFile = masyarakat.NIK + surat.ID + "2" + "-" + surat.Jns_surat + ".pdf"

	//Menyimpan file dokumen pribadi pada folder project
	if err := c.SaveFile(suratKeterangan, fmt.Sprintf("./public/%s/%s", surat.Jns_surat, namaFile)); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Create data dokumen surat
	doc_syarat.Filename = namaFile
	doc_syarat.ID = ""
	if err := tx.Create(&doc_syarat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"msg": "Surat berhasil diajukan"})
}
