package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cmellojr/modo-locadora/internal/config"
	"github.com/cmellojr/modo-locadora/internal/database"
	"github.com/cmellojr/modo-locadora/internal/handlers"
	"github.com/cmellojr/modo-locadora/internal/jobs"
	"github.com/cmellojr/modo-locadora/internal/middleware"
)

func main() {
	seedFlag := flag.Bool("seed", false, "Populate database with sample data and exit")
	flag.Parse()

	config.LoadConfig()

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

	if *seedFlag {
		if store == nil {
			log.Fatal("Cannot seed: no database connection. Set DATABASE_URL.")
		}
		pgStore, ok := store.(interface {
			ExecRaw(ctx context.Context, sql string) error
		})
		if !ok {
			log.Fatal("Store does not support raw SQL execution.")
		}
		// Apply all migrations in order, then seed.
		// Try "migrations/" first (Docker), fall back to "internal/database/migrations/" (local dev).
		migrationsDir := "internal/database/migrations/"
		if _, err := os.Stat("migrations"); err == nil {
			migrationsDir = "migrations/"
		}
		sqlFiles := []string{
			migrationsDir + "001_initial_schema.sql",
			migrationsDir + "002_update_games_table.sql",
			migrationsDir + "003_membership_and_rental_support.sql",
			migrationsDir + "004_password_notes.sql",
			migrationsDir + "005_auto_return_reputation.sql",
			migrationsDir + "006_activities_feed.sql",
			migrationsDir + "007_seed_initial_data.sql",
			migrationsDir + "008_cover_display.sql",
		}
		for _, f := range sqlFiles {
			data, err := os.ReadFile(f)
			if err != nil {
				log.Fatalf("Failed to read %s: %v", f, err)
			}
			if err := pgStore.ExecRaw(ctx, string(data)); err != nil {
				log.Fatalf("Failed to execute %s: %v", f, err)
			}
			log.Printf("Applied: %s", f)
		}
		log.Println("Migrations + seed concluidos com sucesso!")
		return
	}

	cookieSecret := os.Getenv("COOKIE_SECRET")
	if cookieSecret == "" {
		log.Println("Warning: COOKIE_SECRET not set. Using insecure default. Set it in production!")
		cookieSecret = "modo-locadora-dev-secret-change-me"
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		log.Println("Warning: ADMIN_EMAIL not set. Admin routes will be inaccessible.")
	}

	h := handlers.NewHandler(store, cookieSecret, adminEmail)

	// Start the overdue rental checker background job.
	if store != nil {
		jobs.StartOverdueChecker(ctx, store, 5*time.Minute)
	}

	indexTmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Fatalf("failed to parse index template: %v", err)
	}

	platformsTmpl, err := template.ParseFiles("web/templates/platforms.html")
	if err != nil {
		log.Fatalf("failed to parse platforms template: %v", err)
	}

	gamesTmpl, err := template.ParseFiles("web/templates/games.html")
	if err != nil {
		log.Fatalf("failed to parse games template: %v", err)
	}

	gameDetailTmpl, err := template.ParseFiles("web/templates/game_detail.html")
	if err != nil {
		log.Fatalf("failed to parse game detail template: %v", err)
	}

	adminStockTmpl, err := template.ParseFiles("web/templates/admin_stock.html")
	if err != nil {
		log.Fatalf("failed to parse admin stock template: %v", err)
	}

	adminInventoryTmpl, err := template.ParseFiles("web/templates/admin_inventory.html")
	if err != nil {
		log.Fatalf("failed to parse admin inventory template: %v", err)
	}

	adminEditTmpl, err := template.ParseFiles("web/templates/admin_edit.html")
	if err != nil {
		log.Fatalf("failed to parse admin edit template: %v", err)
	}

	carteirinhaTmpl, err := template.ParseFiles("web/templates/carteirinha.html")
	if err != nil {
		log.Fatalf("failed to parse carteirinha template: %v", err)
	}

	adminReturnsTmpl, err := template.ParseFiles("web/templates/admin_returns.html")
	if err != nil {
		log.Fatalf("failed to parse admin returns template: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		h.HandleIndex(w, r, indexTmpl)
	})
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("POST /logout", h.Logout)
	mux.HandleFunc("GET /games", func(w http.ResponseWriter, r *http.Request) {
		h.ListGames(w, r, platformsTmpl, gamesTmpl)
	})
	mux.HandleFunc("GET /search", h.SearchGame)

	// Admin routes — protected by RequireAdmin middleware.
	mux.HandleFunc("GET /admin/stock", middleware.RequireAdmin(cookieSecret, adminEmail, store, func(w http.ResponseWriter, r *http.Request) {
		h.AdminStock(w, r, adminStockTmpl)
	}))
	mux.HandleFunc("POST /admin/purchase", middleware.RequireAdmin(cookieSecret, adminEmail, store, h.PurchaseGame))
	mux.HandleFunc("GET /admin/inventory", middleware.RequireAdmin(cookieSecret, adminEmail, store, func(w http.ResponseWriter, r *http.Request) {
		h.AdminInventory(w, r, adminInventoryTmpl)
	}))
	mux.HandleFunc("GET /admin/edit/{id}", middleware.RequireAdmin(cookieSecret, adminEmail, store, func(w http.ResponseWriter, r *http.Request) {
		h.EditGame(w, r, adminEditTmpl)
	}))
	mux.HandleFunc("POST /admin/update-game", middleware.RequireAdmin(cookieSecret, adminEmail, store, h.UpdateGame))
	mux.HandleFunc("GET /admin/returns", middleware.RequireAdmin(cookieSecret, adminEmail, store, func(w http.ResponseWriter, r *http.Request) {
		h.AdminReturns(w, r, adminReturnsTmpl)
	}))
	mux.HandleFunc("POST /admin/return-game", middleware.RequireAdmin(cookieSecret, adminEmail, store, h.ReturnGame))

	// Member routes — protected by RequireAuth middleware.
	mux.HandleFunc("GET /carteirinha", middleware.RequireAuth(cookieSecret, func(w http.ResponseWriter, r *http.Request) {
		h.Carteirinha(w, r, carteirinhaTmpl)
	}))
	mux.HandleFunc("POST /rent", middleware.RequireAuth(cookieSecret, h.RentGame))
	mux.HandleFunc("POST /carteirinha/notes", middleware.RequireAuth(cookieSecret, h.SavePasswordNotes))
	mux.HandleFunc("POST /carteirinha/redeem", middleware.RequireAuth(cookieSecret, h.HandleRedeem))
	mux.HandleFunc("POST /carteirinha/return", middleware.RequireAuth(cookieSecret, h.HandleMemberReturn))

	// Serve static files from web/static
	fileServer := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("POST /members", h.CreateMember)
	mux.HandleFunc("GET /games/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.GameDetailPage(w, r, gameDetailTmpl)
	})

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
