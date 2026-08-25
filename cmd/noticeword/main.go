package main

import (
	"flag"
	"log"
	"net/http"

	"noticeword/internal/httpapi"
	"noticeword/internal/service"
	"noticeword/internal/store"
)

func main() {
	dbPath := flag.String("db", "noticeword.db", "path to the bbolt database")
	address := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
	database, err := store.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	clock := store.NewSequenceClock("2025-01-01T00:00:00Z")
	app := service.New(database, clock)
	server := httpapi.NewServer(app)
	log.Printf("noticeword listening on %s", *address)
	if err := http.ListenAndServe(*address, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
