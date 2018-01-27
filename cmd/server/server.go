package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	s "github.com/pinub/pinub/http"
	"github.com/pinub/pinub/postgres"
)

func main() {
	client, err := postgres.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer client.Close() // nolint: errcheck

	h := s.New(client, "./views")
	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
