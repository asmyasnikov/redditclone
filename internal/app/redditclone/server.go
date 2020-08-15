package redditclone

import (
	"github.com/asmyasnikov/redditclone/api"
	"github.com/asmyasnikov/redditclone/internal/logger"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type server struct {
	log *logger.Logger
}

// Run run server with configuration cfg
func Run(cfg *api.Configuration) error {
	l, err := logger.New(cfg.Log)
	if err != nil {
		log.Fatalf("cannot init logger with config '%v': %v", cfg.Log, err)
		return err
	}
	s := server{
		log: l,
	}
	if err := http.ListenAndServe(cfg.Server.HTTPListen, router()); err != nil {
		s.log.Error(err)
	}
	return nil
}

func router() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
	return r
}
