package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/rugwirobaker/helmes"
	"github.com/rugwirobaker/helmes/api"
)

func main() {

	port := os.Getenv("PORT")
	id := os.Getenv("HELMES_SMS_APP_ID")
	secret := os.Getenv("HELMES_SMS_APP_SECRET")
	sender := os.Getenv("HELMES_SENDER_IDENTITY")
	callback := os.Getenv("HELMES_CALLBACK_URL")

	cli := provideClient()

	service, err := helmes.NewSendService(cli, id, secret, sender, callback)
	if err != nil {
		log.Fatalf("could not initialize sms service: %v", err)
	}

	events := helmes.NewPubsub()
	defer events.Close()

	log.Println("initialized helmes api")
	api := api.New(service, events)
	mux := chi.NewMux()
	mux.Mount("/api", api.Handler())

	if len(port) == 0 {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("starting application at port %v", port)

	err = srv.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-idleConnsClosed
}
