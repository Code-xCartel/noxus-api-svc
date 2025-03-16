package _api

import (
	"database/sql"
	"github.com/Code-xCartel/noxus-api-svc/service/auth"
	"github.com/Code-xCartel/noxus-api-svc/service/friends"
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
	auth.Router(r.router, auth.NewAuthStore(r.store))
	friends.Router(r.router, friends.NewFriendsStore(r.store, auth.NewAuthStore(r.store)))
}
