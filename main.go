package main

import (
	"Backend_TA/controllers/allsuratcontrollers"
	"Backend_TA/controllers/authcontrollers"
	"Backend_TA/controllers/ktpbarucontrollers"
	"Backend_TA/controllers/ktplamacontrollers"
	"Backend_TA/controllers/masyarakatcontrollers"
	"Backend_TA/controllers/ujicobacontrollers"
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
	ktplama := api.Group("/ktplama")
	pengujian := api.Group("/uji")

	auth.Post("/refresh", authcontrollers.RefreshToken)
	auth.Post("/register", authcontrollers.Register)
	auth.Post("/login", authcontrollers.Login)

	profile.Get("/" /*, middlewares.Auth*/, masyarakatcontrollers.Show)
	profile.Post("/" /*, middlewares.Auth*/, authcontrollers.Register)
	profile.Get("/:nik" /*, middlewares.Auth*/, masyarakatcontrollers.ShowId)
	profile.Put("/:nik" /*, middlewares.Auth*/, masyarakatcontrollers.UpdateProfile)
	profile.Put("/password/:nik" /*, middlewares.Auth*/, masyarakatcontrollers.UpdatePassword)
	profile.Delete("/:nik" /*, middlewares.Auth*/, masyarakatcontrollers.DeleteProfile)

	//API surat KTP Baru
	ktpbaru.Post("/:id" /*, middlewares.Auth*/, ktpbarucontrollers.CreateKTPBaru)

	//API surat KTP lama
	ktplama.Post("/:id" /*, middlewares.Auth*/, ktplamacontrollers.CreateKTPLama)

	//API utk semua sura tanpat terkecuali
	allsurat.Get("/" /*, middlewares.Auth*/, allsuratcontrollers.ShowSurat)
	allsurat.Get("/show/:id" /*, middlewares.Auth*/, allsuratcontrollers.ShowSuratById)  //Melihat seluruh data surat berdasarkan id surat
	allsurat.Get("/:id" /*, middlewares.Auth*/, allsuratcontrollers.ShowSuratByNik)      //Melihat data surat berdasarkan NIK user yang mengajukan
	allsurat.Get("/data_doc/:id", allsuratcontrollers.ShowDataDoc)                       //Melihat seluruh data dokumen syarat dari pengajuan surat
	allsurat.Get("/doc/:id" /*, middlewares.Auth*/, allsuratcontrollers.ShowDocSyarat)   //Melakukan download dokumen syarat setiap surat berdasarkan id doc_syarat
	allsurat.Put("/doc/:id" /*, middlewares.Auth*/, allsuratcontrollers.UpdateDocSyarat) //Melakukan update input dokumen syarat
	allsurat.Put("/:id" /*, middlewares.Auth*/, allsuratcontrollers.Update)              //Melakukan update hanya pada data surat tanpa mengubah data dokumen syarat
	allsurat.Delete("/:id" /*, middlewares.Auth*/, allsuratcontrollers.Delete)           //Menghapus data surat

	// Progam test enkrip isi file
	pengujian.Post("/enkrip/", ujicobacontrollers.EncryptFile)
	pengujian.Post("/dekrip/", ujicobacontrollers.DecryptFile)

	//Progam test enkrip file
	pengujian.Post("/enkripfile/", ujicobacontrollers.UjiCobaFile)
	pengujian.Post("/dekripfile/", ujicobacontrollers.UjiCobaFileDek)
	// pengujian.Post("/enkripfile2/", ktpbarucontrollers.UjiCobaFile)
	// pengujian.Post("/dekripfile2/", ktpbarucontrollers.UjiCobaDekripsiFile)
	app.Listen(":4001")
}
