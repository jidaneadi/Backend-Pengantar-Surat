package main

import (
	"Backend_TA/controllers/allsuratcontrollers"
	"Backend_TA/controllers/authcontrollers"
	"Backend_TA/controllers/masyarakatcontrollers"
	"Backend_TA/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	models.ConnectDB()
	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/")
	auth := api.Group("auth")
	profile := api.Group("/profile")
	allsurat := api.Group("/surat")

	auth.Post("/refresh", authcontrollers.RefreshToken)
	auth.Post("/register", authcontrollers.Register)
	auth.Post("/login", authcontrollers.Login)

	profile.Get("/", masyarakatcontrollers.Show)
	profile.Post("/", authcontrollers.Register)
	profile.Get("/:nik", masyarakatcontrollers.ShowId)
	profile.Put("/:nik", masyarakatcontrollers.UpdateProfile)
	profile.Put("/password/:nik", masyarakatcontrollers.UpdatePassword)
	profile.Delete("/:nik", masyarakatcontrollers.DeleteProfile)

	//API utk semua surat tanpat terkecuali
	allsurat.Post("/", allsuratcontrollers.Create)
	allsurat.Get("/", allsuratcontrollers.Show)
	allsurat.Get("/:id", allsuratcontrollers.ShowId)
	allsurat.Put("/:id", allsuratcontrollers.Update)
	allsurat.Put("/:id", allsuratcontrollers.UpdateStatus)
	allsurat.Delete(":/id", allsuratcontrollers.Delete)
	app.Listen(":4001")
}
