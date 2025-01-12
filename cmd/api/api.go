package _api

import (
	"database/sql"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) RunApiServer() error {
	router := http.NewServeMux()
	rootRouter := NewRootRouter(router, s.db)
	rootRouter.RegisterRoutes()
	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}
	log.Printf("api svc started %s", s.addr)
	return server.ListenAndServe()
}
