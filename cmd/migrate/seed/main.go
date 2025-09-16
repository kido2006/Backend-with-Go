package main

import (
	"backendwithgo/internal/db"
	"backendwithgo/internal/env"
	"backendwithgo/internal/store"
	"log"
)

func main() {
	addr := env.GetString("DB_ADDR", "niga:123456789@tcp(db:3306)/myapp?parseTime=true")

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	store := store.NewSQL(conn)

	db.Seed(store, conn)
}
