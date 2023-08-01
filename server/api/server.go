package api

import "github.com/gofiber/fiber/v2"

type Server struct {
	listenAddr string
}

func New(listenAddr string) *Server {
	server := Server{
		listenAddr: listenAddr,
	}
	return &server
}

func (s *Server) Start() error {
	app := fiber.New()
	app.Get("/state", GetState)

	app.Post("/rename", Rename)
	app.Post("/move", Move)
	app.Post("/copy", Copy)
	app.Post("/delete", Delete)

	app.Post("/mkdir", MkDir)

	app.Put("/upload", StartUploading)
	app.Put("/upload/chunk", UploadChunk)
	app.Put("/upload/completed", UploadComplete)

	app.Get("/download", StartDownloading)
	app.Get("/download/chunk", DownloadChunk)

	app.Post("/authenticate", Authenticate)
	app.Post("/refresh", Refresh)
	app.Post("/logout", Logout)

	err := app.Listen(s.listenAddr)
	return err
}
