package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cmellojr/modo-locadora/internal/database"
	"github.com/cmellojr/modo-locadora/internal/handlers"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var store database.Store
	var err error

	connString := os.Getenv("DATABASE_URL")
	if connString != "" {
		store, err = database.NewPostgresStore(ctx, connString)
		if err != nil {
			log.Printf("Warning: failed to initialize real store: %v. Proceeding without database.", err)
		}
	} else {
		log.Println("No DATABASE_URL provided. Proceeding without database.")
	}

	h := handlers.NewHandler(store)

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("POST /members", h.CreateMember)
	mux.HandleFunc("GET /games/{id}", h.GetGame)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		fmt.Printf("Rental Mode Server starting on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if store != nil {
		if pgStore, ok := store.(interface{ Close() }); ok {
			pgStore.Close()
		}
	}

	fmt.Println("Server gracefully stopped")
}
