package main

import (
	"forum/database"
	"forum/internal/controller"
	"forum/internal/repository"
	"forum/internal/service"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	httpServer *http.Server
}

func main() {
	postgres := database.NewDB()

	repos := repository.NewRepository(postgres)
	services := service.NewService(repos)
	handler := controller.NewHandler(services)

	router := handler.InitRoutes()

	srv := new(Server)

	log.Println("Starting the server")
	if err := srv.Start("8000", router); err != nil {
		log.Fatalf("error server: %v", err)
	}
}

func (s *Server) Start(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}
