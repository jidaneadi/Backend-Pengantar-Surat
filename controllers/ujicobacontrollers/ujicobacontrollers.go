package ujicobacontrollers

import (
	"Backend_TA/utils"
	"encoding/hex"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func EncryptFile(c *fiber.Ctx) error {
	key, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	nonce, _ := hex.DecodeString("000000000000000000000000")

	var keyArray [32]byte
	var nonceArray [12]byte

	copy(keyArray[:], key)
	copy(nonceArray[:], nonce)

	chacha := &utils.ChaCha20{
		Key:     keyArray,
		Nonce:   nonceArray,
		Counter: 1,
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Failed to get file: " + err.Error())
	}

	inputFilePath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveFile(file, inputFilePath); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to save file: " + err.Error())
	}

	outputFilePath := filepath.Join("./public/hasil", file.Filename)

	err = utils.EncryptDecryptFile2(chacha, inputFilePath, outputFilePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to encrypt file: " + err.Error())
	}

	return c.SendString("File encrypted successfully! Saved at " + outputFilePath)
}

func UjiCobaFile(c *fiber.Ctx) error {
	key, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	nonce, _ := hex.DecodeString("000000000000000000000000")

	var keyArray [32]byte
	var nonceArray [12]byte
	// nameDokumenPribadi := dokumenPribadi.Filename
	// chacha := &utils.ChaCha20{} // Create instance of ChaCha20

	copy(keyArray[:], key)
	copy(nonceArray[:], nonce)

	chacha := &utils.ChaCha20{
		Key:     keyArray,
		Nonce:   nonceArray,
		Counter: 1,
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Failed to get file: " + err.Error())
	}

	input, err := file.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to open file: " + err.Error())
	}
	defer input.Close()

	// newName := strings.Split(file.Filename, ".")
	// name := newName[0] + ".txt"

	tempFilePath := filepath.Join(os.TempDir(), file.Filename)
	err = utils.EncryptDecryptFile(chacha, input, tempFilePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to encrypt file: " + err.Error())
	}

	outputDir := "./public/hasil/"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to create output directory: " + err.Error())
	}

	encryptedFile := filepath.Join(outputDir, file.Filename)
	// encryptedFile := filepath.Join(outputDir, file.Filename)
	err = os.Rename(tempFilePath, encryptedFile)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to save file: " + err.Error())
	}

	encryptedData, err := os.ReadFile(encryptedFile)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal membaca file enkripsi: " + err.Error())
	}

	return c.SendString("File encrypted successfully!" + string(encryptedData))
}

func DecryptFile(c *fiber.Ctx) error {
	key, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	nonce, _ := hex.DecodeString("000000000000000000000000")

	var keyArray [32]byte
	var nonceArray [12]byte

	copy(keyArray[:], key)
	copy(nonceArray[:], nonce)

	chacha := &utils.ChaCha20{
		Key:     keyArray,
		Nonce:   nonceArray,
		Counter: 1,
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Failed to get file: " + err.Error())
	}

	input, err := file.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to open file: " + err.Error())
	}
	defer input.Close()

	tempFilePath := filepath.Join(os.TempDir(), file.Filename)
	err = utils.EncryptDecryptFile(chacha, input, tempFilePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to decrypt file: " + err.Error())
	}

	return c.Download(tempFilePath)
}

func UjiCobaFileDek(c *fiber.Ctx) error {
	key, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	nonce, _ := hex.DecodeString("000000000000000000000000")

	var keyArray [32]byte
	var nonceArray [12]byte
	// nameDokumenPribadi := dokumenPribadi.Filename
	// chacha := &utils.ChaCha20{} // Create instance of ChaCha20

	copy(keyArray[:], key)
	copy(nonceArray[:], nonce)

	chacha := &utils.ChaCha20{
		Key:     keyArray,
		Nonce:   nonceArray,
		Counter: 1,
	}
	// chacha := &utils.ChaCha20{}
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Failed to get file: " + err.Error())
	}

	input, err := file.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to open file: " + err.Error())
	}
	defer input.Close()

	tempFilePath := filepath.Join(os.TempDir(), file.Filename)
	err = utils.EncryptDecryptFile(chacha, input, tempFilePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to decrypt file: " + err.Error())
	}

	outputDir := "./public/hasildekrip/"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to create output directory: " + err.Error())
	}

	decryptedFile := filepath.Join(outputDir, file.Filename)
	err = os.Rename(tempFilePath, decryptedFile)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to save file: " + err.Error())
	}
	hasilDecrypt, err := os.ReadFile(decryptedFile)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal membaca file enkripsi: " + err.Error())
	}
	return c.SendString("File decrypted successfully!" + string(hasilDecrypt))
}
