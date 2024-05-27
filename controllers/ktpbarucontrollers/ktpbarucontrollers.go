package ktpbarucontrollers

import (
	"Backend_TA/models"
	"Backend_TA/utils"
	"fmt"
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

	errCheckContentTypeDP := utils.CheckContentType(dokumenPribadi, "application/pdf", "application/PDF")
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

	errCheckContentType := utils.CheckContentType(suratKeterangan, "application/pdf", "application/PDF")
	if errCheckContentType != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": errCheckContentType.Error(),
		})
	}

	filename := &suratKeterangan.Filename
	extenstionFile := filepath.Ext(*filename)

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

// func checkContentType(file *multipart.FileHeader, contentTypes ...string) error {
// 	if len(contentTypes) > 0 {
// 		for _, contentType := range contentTypes {
// 			contentTypeFile := file.Header.Get("Content-Type")
// 			if contentTypeFile == contentType {
// 				return nil
// 			}
// 		}

// 		return errors.New("not allowed file type")
// 	} else {
// 		return errors.New("not found content type to be checking")
// 	}
// }

// <<<<<<<<====================SIMPAN DOKUMEN KE TXT================================>>>>>>>>
// func UjiCobaFile(c *fiber.Ctx) error {
// 	key := [32]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00,
// 		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
// 	nonce := [12]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb}
// 	counter := uint32(123)

// 	chacha := &utils.ChaCha20{
// 		Key:     key,
// 		Nonce:   nonce,
// 		Counter: counter,
// 	}
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).SendString("Gagal mendapatkan file: " + err.Error())
// 	}

// 	input, err := file.Open()
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString("Gagal membuka file: " + err.Error())
// 	}
// 	defer input.Close()

// 	// Mendapatkan nama file baru dengan ekstensi .txt
// 	newName := strings.Split(file.Filename, ".")
// 	name := newName[0] + ".txt"

// 	// Membuat jalur file sementara
// 	tempFilePath := filepath.Join(os.TempDir(), name)

// 	// Melakukan enkripsi file dan menyimpannya ke file sementara
// 	err = utils.EncryptDecryptFile(chacha, input, tempFilePath)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString("Gagal melakukan enkripsi file: " + err.Error())
// 	}

// 	// Membuat direktori output jika tidak ada
// 	outputDir := "./public/hasil/"
// 	err = os.MkdirAll(outputDir, os.ModePerm)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString("Gagal membuat direktori output: " + err.Error())
// 	}

// 	// Membaca isi file yang telah dienkripsi dari file sementara
// 	encryptedData, err := os.ReadFile(tempFilePath)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString("Gagal membaca file enkripsi: " + err.Error())
// 	}

// 	// Menyiapkan jalur file .txt untuk menyimpan hasil enkripsi
// 	encryptedFilePath := filepath.Join(outputDir, name)

// 	// Menulis hasil enkripsi ke file .txt di direktori output
// 	err = os.WriteFile(encryptedFilePath, encryptedData, 0644)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString("Gagal menyimpan hasil enkripsi: " + err.Error())
// 	}

// 	return c.SendString("File berhasil dienkripsi")
// }

// func UjiCobaFile(c *fiber.Ctx) error {
// 	key := [32]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00,
// 		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
// 	nonce := [12]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb}
// 	counter := uint32(123)

// 	chacha := &utils.ChaCha20{
// 		Key:     key,
// 		Nonce:   nonce,
// 		Counter: counter,
// 	}

// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).SendString("Gagal mendapatkan file: " + err.Error())
// 	}

// 	input, err := file.Open()
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString("Gagal membuka file: " + err.Error())
// 	}
// 	defer input.Close()

// 	// Mendapatkan nama file baru dengan ekstensi .txt
// 	newName := strings.Split(file.Filename, ".")
// 	name := newName[0] + ".txt"

// 	// Membuat jalur file sementara
// 	tempFilePath := filepath.Join(os.TempDir(), name)

// 	// Melakukan enkripsi file dan menyimpannya ke file sementara
// 	err = utils.EncryptDecryptFile(chacha, input, tempFilePath)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString("Gagal melakukan enkripsi file: " + err.Error())
// 	}

// 	// Membuat direktori output jika tidak ada
// 	outputDir := "./public/hasil/"
// 	err = os.MkdirAll(outputDir, os.ModePerm)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString("Gagal membuat direktori output: " + err.Error())
// 	}

// 	// Menyiapkan jalur file enkripsi
// 	encryptedFile := filepath.Join(outputDir, file.Filename)

// 	// Memindahkan file sementara ke lokasi yang ditentukan
// 	err = os.Rename(tempFilePath, encryptedFile)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString("Gagal menyimpan file: " + err.Error())
// 	}

// Membaca isi file yang telah dienkripsi
// encryptedData, err := os.ReadFile(encryptedFile)
// if err != nil {
// 	return c.Status(http.StatusInternalServerError).SendString("Gagal membaca file enkripsi: " + err.Error())
// }

// Mengirimkan isi file yang telah dienkripsi sebagai respons
// return c.SendString("File berhasil dienkripsi: " + string(encryptedData))
// dataEnkrip := string(encryptedData)
// 	return c.JSON(fiber.Map{"data": encryptedData})
// }
