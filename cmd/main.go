package main

import (
	_api "github.com/Code-xCartel/noxus-api-svc/cmd/api"
	"github.com/Code-xCartel/noxus-api-svc/config"
	"github.com/Code-xCartel/noxus-api-svc/db"
	_ "github.com/lib/pq"
)

func main() {
	dbDriver, dbError := db.NewSQLStorage(config.Envs.PostgresConnectionString)
	if dbError != nil {
		panic(dbError)
	}
	if err := db.InitStorage(dbDriver); err != nil {
		panic(err)
	}
	server := _api.NewAPIServer(":8080", dbDriver)
	if err := server.RunApiServer(); err != nil {
		panic(err)
	}
}
