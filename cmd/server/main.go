package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/valve"
	"golang.org/x/net/http2"
)

var (
	fs       = flag.NewFlagSet("server", flag.ExitOnError)
	certFile = fs.String("cert.file", "certs/localhost.cert", "SSL certificate")
	keyFile  = fs.String("key.file", "certs/localhost.key", "Private key")
	httpAddr = fs.String("http.addr", "127.0.0.1:8080", "HTTP server address")
)

func main() {
	var err error
	// Our graceful valve shut-off package to manage code preemption and
	// shutdown signaling.
	vlv := valve.New()
	baseCtx := vlv.Context()
	// Parse flags.
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	err = fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("%v", err)
	}
	// Create routing mux.
	r := router()
	// Initialize http2 server.
	srv := http.Server{
		Addr:    *httpAddr,
		Handler: chi.ServerBaseContext(baseCtx, r),
	}
	err = http2.ConfigureServer(&srv, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	// Signals.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Graceful shutdown.
	go func() {
		<-c
		// sig is a ^C, handle it
		log.Println("shutting down..")

		// first valve
		vlv.Shutdown(10 * time.Second)

		// create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// start http shutdown
		srv.Shutdown(ctx)

		// verify, in worst case call cancel via defer
		select {
		case <-time.After(11 * time.Second):
			log.Println("not all connections done")
		case <-ctx.Done():

		}
	}()
	// Start server.
	log.Printf("listening on %s", *httpAddr)
	err = srv.ListenAndServeTLS(*certFile, *keyFile)
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatalf("%v", err)
		}
	}
}
