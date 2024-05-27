package ktpbarucontrollers

import (
	"Backend_TA/config"
	"Backend_TA/models"
	"Backend_TA/utils"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateKTPBaru(c *fiber.Ctx) error {

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

	id := c.Params("id") //Mengambil parameter id
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

	//Melakukan input data surat
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

	//Input file dokumen probadi
	dokumenPribadi, err := c.FormFile("dokumen_pribadi")

	//Mengecek apakah data surat kosong atau tidak
	if dokumenPribadi == nil {
		return c.Status(400).JSON(fiber.Map{"msg": "File dokumen pribadi tidak boleh kosong!"})
	}

	//Melihat apakah ada error lain saat input file
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Memvalidasi bahwa file harus bertipe PDF
	errCheckContentTypeDP := utils.CheckContentType(dokumenPribadi, "application/pdf", "application/PDF")
	if errCheckContentTypeDP != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": errCheckContentTypeDP.Error(),
		})
	}

	//Membuka file input
	inputDokPribadi, err := dokumenPribadi.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal membuka file: " + err.Error())
	}
	defer inputDokPribadi.Close()

	namaFile1 := masyarakat.NIK + surat.ID + "1" + ".txt"

	// Membuat jalur file sementara
	pathDokPribadi := filepath.Join(os.TempDir(), namaFile1)

	// Melakukan enkripsi file dan menyimpannya ke file sementara
	err = utils.EncryptDecryptFile(chacha, inputDokPribadi, pathDokPribadi)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal melakukan enkripsi file: " + err.Error())
	}

	// Membuat direktori output jika tidak ada
	dirOutputDokPribadi := "./public/ktp_baru/"
	err = os.MkdirAll(dirOutputDokPribadi, os.ModePerm)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal membuat direktori output: " + err.Error())
	}

	// Membaca isi file yang telah dienkripsi dari file sementara
	enkripDokPribadi, err := os.ReadFile(pathDokPribadi)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal membaca file enkripsi: " + err.Error())
	}

	// Menyiapkan jalur file .txt untuk menyimpan hasil enkripsi
	pathFileDokPribadiNew := filepath.Join(dirOutputDokPribadi, namaFile1)

	// Menulis hasil enkripsi ke file .txt di direktori output
	err = os.WriteFile(pathFileDokPribadiNew, enkripDokPribadi, 0644)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal menyimpan hasil enkripsi: " + err.Error())
	}

	//======================================================
	//Input file keterangan RT
	suratKeterangan, err := c.FormFile("keterangan_rt")

	//Pengecekan apakah kosong atau tidak
	if suratKeterangan == nil {
		return c.Status(400).JSON(fiber.Map{"msg": "File dokumen pribadi tidak boleh kosong!"})
	}

	//Mengecek apakah ada error yang lain atau tidak
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	//Pengecekan bahwa file harus berformat pdf
	errCheckContentType := utils.CheckContentType(suratKeterangan, "application/pdf", "application/PDF")
	if errCheckContentType != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": errCheckContentType.Error(),
		})
	}

	inputSuratKet, err := suratKeterangan.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal membuka file: " + err.Error())
	}
	defer inputSuratKet.Close()

	//Nama file ke 2
	namaFile2 := masyarakat.NIK + surat.ID + "2" + ".txt"

	// Membuat instance ChaCha20 untuk keterangan RT
	chachaSuratKet := &utils.ChaCha20{
		Key:     keyArray,
		Nonce:   nonceArray,
		Counter: 1,
	}

	// Membuat jalur file sementara
	tempFilePath := filepath.Join(os.TempDir(), namaFile2)

	// Melakukan enkripsi file dan menyimpannya ke file sementara
	err = utils.EncryptDecryptFile(chachaSuratKet, inputSuratKet, tempFilePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal melakukan enkripsi file: " + err.Error())
	}

	// Membuat direktori output jika tidak ada
	outputDir := "./public/ktp_baru/"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal membuat direktori output: " + err.Error())
	}

	// Membaca isi file yang telah dienkripsi dari file sementara
	encryptedData, err := os.ReadFile(tempFilePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal membaca file enkripsi: " + err.Error())
	}

	// Menyiapkan jalur file .txt untuk menyimpan hasil enkripsi
	encryptedFilePath := filepath.Join(outputDir, namaFile2)

	// Menulis hasil enkripsi ke file .txt di direktori output
	err = os.WriteFile(encryptedFilePath, encryptedData, 0644)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal menyimpan hasil enkripsi: " + err.Error())
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

	fmt.Println(dataKey)
	fmt.Println(dataNonce)
	return c.JSON(fiber.Map{"msg": "Surat berhasil diajukan"})
}

//================================================================
// func UjiCobaFile(c *fiber.Ctx) error {
// 	dataKey := config.RenderEnv("KEY_CHACHA20")
// 	dataNonce := config.RenderEnv("NONCE")
// 	key, _ := hex.DecodeString(dataKey)
// 	nonce, _ := hex.DecodeString(dataNonce)

// 	var keyArray [32]byte
// 	var nonceArray [12]byte
// 	// nameDokumenPribadi := dokumenPribadi.Filename
// 	// chacha := &utils.ChaCha20{} // Create instance of ChaCha20

// 	copy(keyArray[:], key)
// 	copy(nonceArray[:], nonce)

// 	chacha := &utils.ChaCha20{
// 		Key:     keyArray,
// 		Nonce:   nonceArray,
// 		Counter: 1,
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

// func UjiCobaDekripsiFile(c *fiber.Ctx) error {
// 	dataKey := config.RenderEnv("KEY_CHACHA20")
// 	dataNonce := config.RenderEnv("NONCE")
// 	key, _ := hex.DecodeString(dataKey)
// 	nonce, _ := hex.DecodeString(dataNonce)

// 	var keyArray [32]byte
// 	var nonceArray [12]byte

// 	copy(keyArray[:], key)
// 	copy(nonceArray[:], nonce)

// 	chacha := &utils.ChaCha20{
// 		Key:     keyArray,
// 		Nonce:   nonceArray,
// 		Counter: 1,
// 	}

// 	// Ambil nama file input dari parameter
// 	inputFileName := "330900611060100012031.txt"
// 	if inputFileName == "" {
// 		return c.Status(http.StatusBadRequest).SendString("Nama file input tidak boleh kosong")
// 	}

// 	// Persiapkan jalur file input dan output
// 	inputFilePath := filepath.Join("./public/hasil", inputFileName)
// 	outputFileName := strings.TrimSuffix(inputFileName, filepath.Ext(inputFileName)) + ".pdf"
// 	outputFilePath := filepath.Join("./public/hasil", outputFileName)

// 	// Dekripsi file
// 	err := DecryptFile(chacha, inputFilePath, outputFilePath)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString("Gagal melakukan dekripsi: " + err.Error())
// 	}

// 	return c.SendFile(outputFilePath)
// }

// func DecryptFile(chacha *utils.ChaCha20, inputFilePath, outputFilePath string) error {
// 	// Baca data terenkripsi dari file
// 	encryptedData, err := os.ReadFile(inputFilePath)
// 	if err != nil {
// 		return err
// 	}

// 	// Dekripsi data
// 	decryptedData := make([]byte, len(encryptedData))
// 	chacha.Counter = 1 // Reset counter untuk dekripsi
// 	chacha.XORKeyStream(decryptedData, encryptedData)

// 	// Simpan hasil dekripsi ke file output
// 	err = os.WriteFile(outputFilePath, decryptedData, 0644)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

//==================================
