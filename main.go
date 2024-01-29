package main

import (
	"Backend_TA/controllers/allsuratcontrollers"
	"Backend_TA/controllers/authcontrollers"
	"Backend_TA/controllers/ktpbarucontrollers"
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
	ktpbaru := api.Group("/ktpbaru")

	auth.Post("/refresh", authcontrollers.RefreshToken)
	auth.Post("/register", authcontrollers.Register)
	auth.Post("/login", authcontrollers.Login)

	profile.Get("/" /*middlewares.Auth,*/, masyarakatcontrollers.Show)
	profile.Post("/" /*middlewares.Auth,*/, authcontrollers.Register)
	profile.Get("/:nik" /*middlewares.Auth,*/, masyarakatcontrollers.ShowId)
	profile.Put("/:nik" /*middlewares.Auth,*/, masyarakatcontrollers.UpdateProfile)
	profile.Put("/password/:nik" /*middlewares.Auth,*/, masyarakatcontrollers.UpdatePassword)
	profile.Delete("/:nik" /*middlewares.Auth,*/, masyarakatcontrollers.DeleteProfile)

	//API surat KTP Baru
	ktpbaru.Post("/:id" /*middlewares.Auth,*/, ktpbarucontrollers.CreateKTPBaru)

	//API utk semua sura tanpat terkecuali
	allsurat.Get("/" /*middlewares.Auth,*/, allsuratcontrollers.ShowSurat)
	allsurat.Get("/:id" /*middlewares.Auth,*/, allsuratcontrollers.ShowSuratByNik)
	allsurat.Get("/doc/:id" /*middlewares.Auth,*/, allsuratcontrollers.ShowDocSyarat)
	allsurat.Put("/:id" /*middlewares.Auth,*/, allsuratcontrollers.Update)
	allsurat.Put("/:id" /*middlewares.Auth,*/, allsuratcontrollers.UpdateStatus)
	allsurat.Delete(":/id" /*middlewares.Auth,*/, allsuratcontrollers.Delete)
	app.Listen(":4001")
}
