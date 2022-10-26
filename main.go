package main

import (
	"database/sql"
	"github.com/JustinWG/simpleBank/util"
	"log"

	"github.com/JustinWG/simpleBank/api"
	db "github.com/JustinWG/simpleBank/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	server.Start(config.ServerAddress)
}
