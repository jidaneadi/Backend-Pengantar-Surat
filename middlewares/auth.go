package middlewares

import (
	"Backend_TA/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {

	//Mengambil authorization dari header
	authorization := c.Get("Authorization")
	if authorization == "" {
		return c.Status(401).JSON(fiber.Map{"msg": "Missing authorization"})
	}

	//Memisahkan bearer dan token dari header Authorization
	tokenString := ""
	splitToken := strings.Split(authorization, "Bearer ")
	if len(splitToken) > 1 {
		tokenString = splitToken[1]
	}

	//Memeriksa token
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"msg": "Missing Token"})
	}

	//Melakukan verifikasi token
	verify, err := utils.VerifyAccesToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"msg": "Invalid acces token"})
	}

	//Proses dapat lanjut
	c.Locals("jwt", verify.Claims)
	return c.Next()
}

func PermissionCreate(c *fiber.Ctx) error {
	return c.Next()
}
