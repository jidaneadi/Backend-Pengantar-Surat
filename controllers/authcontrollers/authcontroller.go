package authcontrollers

import (
	"Backend_TA/models"
	"Backend_TA/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	if user.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Email required"})
	}

	if user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Password required"})
	}

	tx := models.DB
	var email models.User
	cek := tx.Where("email =?", user.Email).First(&email)
	if cek.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "Email not found"})
	}

	user.Password = utils.EncryptHash(user.Password)
	if user.Password != email.Password {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"msg": "Password Failed"})
	}

	refreshClaim := jwt.MapClaims{}
	refreshClaim["id"] = email.ID
	refreshClaim["role"] = email.Role
	refreshClaim["exp"] = time.Now().Add(time.Hour * 5).Unix()

	accesClaims := jwt.MapClaims{}
	accesClaims["id"] = email.ID
	accesClaims["role"] = email.Role
	accesClaims["exp"] = time.Now().Add(time.Minute * 3).Unix()

	accesToken, err := utils.GenerateAccesTokens(&accesClaims)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"msg": "Wrong credential"})
	}

	refreshToken, err := utils.GenerateRefreshTokens(&refreshClaim)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"msg": "Wrong credential in Refresh"})
	}

	var masyarakat models.Masyarakat

	cekId := models.DB.Where("nik = ?", email.ID).First(&masyarakat)
	if cekId.RowsAffected > 0 {
		return c.JSON(fiber.Map{
			"role":          email.Role,
			"nama":          masyarakat.Nama,
			"acces_token":   accesToken,
			"refresh_token": refreshToken,
			"msg":           "Login Berhasil",
		})
	}
	return c.JSON(fiber.Map{
		"msg":           "Login Berhasil",
		"role":          email.Role,
		"acces_token":   accesToken,
		"refresh_token": refreshToken,
	})

}

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	if user.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Nik kosong"})
	}

	if user.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "E-mail kosong"})
	}

	if user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Password kosong"})
	}

	if user.Konf_pass == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Konfirmasi password kosong"})
	}

	if err := models.ValidateUser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg_validate": err.Error()})
	}

	if user.Password != user.Konf_pass {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Password dan konfirmasi password tidak sama"})
	}

	tx := models.DB
	var cekData models.User
	cekNik := tx.Where("id = ?", user.ID).First(&cekData)
	if cekNik.RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{"msg": "NIK atau E-mail sudah digunakan"})
	}
	cekEmail := tx.Where("email = ?", user.Email).First(&cekData)
	if cekEmail.RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{"msg": "NIK atau E-mail sudah digunakan"})
	}

	user.Password = utils.EncryptHash(user.Password)
	user.Konf_pass = utils.EncryptHash(user.Konf_pass)

	var masyarakat models.Masyarakat
	masyarakat.NIK = user.ID
	if err := c.BodyParser(&masyarakat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	if masyarakat.Nama == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Nama kosong"})
	}

	if masyarakat.No_hp == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Nomor hp kosong"})
	}

	if masyarakat.Tempat_lahir == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Tempat lahir kosong"})
	}

	if masyarakat.Birthday == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Tanggal lahir kosong"})
	}

	if masyarakat.Alamat == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Alamat kosong"})
	}

	if err := models.ValidateMasyarakat(&masyarakat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	if err := tx.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg_required": err.Error()})
	}

	if err := tx.Create(&masyarakat).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg_required": err.Error()})
	}

	return c.JSON(fiber.Map{
		"msg": "Registrasi Berhasil",
	})
}

func RefreshToken(c *fiber.Ctx) error {

	var refreshToken models.Token
	if err := c.BodyParser(&refreshToken); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	if refreshToken.Refresh_token == "" {
		return c.Status(404).JSON(fiber.Map{"msg": "Refresh Token Kosong"})
	}

	claims, err := utils.DecodeRefreshTokens(refreshToken.Refresh_token)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"msg": "Invalid refresh token",
		})
	}

	newClaims := jwt.MapClaims{}
	newClaims["id"] = claims["id"].(string)
	newClaims["role"] = claims["role"].(string)
	newClaims["exp"] = time.Now().Add(time.Hour * 5).Unix()

	newAccesClaims := jwt.MapClaims{}
	newAccesClaims["id"] = claims["id"].(string)
	newAccesClaims["role"] = claims["role"].(string)
	newAccesClaims["exp"] = time.Now().Add(time.Minute * 3).Unix()

	newAccesTokens, err := utils.GenerateAccesTokens(&newAccesClaims)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"msg": "Failed to generate new access token",
		})
	}

	newRefreshToken, err := utils.GenerateRefreshTokens(&newClaims)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"msg": "Failed to generate new refresh token",
		})
	}

	return c.JSON(fiber.Map{
		"acces_token":   newAccesTokens,
		"refresh_token": newRefreshToken,
		"msg":           "Refresh token generated successfully",
	})
}
