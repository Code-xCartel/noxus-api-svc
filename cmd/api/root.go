package _api

import (
	"database/sql"
	"net/http"
)

type Root struct {
	router *http.ServeMux
	store  *sql.DB
}

func NewRootRouter(router *http.ServeMux, store *sql.DB) *Root {
	return &Root{router: router, store: store}
}

func (r *Root) RegisterRoutes() {
	// Add routes here
}
