package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/MixinNetwork/mixin-wallet/config"
	"github.com/MixinNetwork/mixin-wallet/durable"
	"github.com/MixinNetwork/mixin-wallet/services"
)

func main() {
	service := flag.String("service", "http", "run a service")
	flag.Parse()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DatebaseUser, config.DatabasePassword, config.DatabaseHost, config.DatabasePort, config.DatabaseName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panicln(err)
	}

	defer db.Close()
	database, err := durable.NewDatabase(context.Background(), db)
	if err != nil {
		log.Panicln(err)
	}

	switch *service {
	case "http":
		err := StartHTTP(database)
		if err != nil {
			log.Println(err)
		}
	default:
		hub := services.NewHub(database)
		err := hub.StartService(*service)
		if err != nil {
			log.Println(err)
		}
	}
}
