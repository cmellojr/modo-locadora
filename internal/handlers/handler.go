package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/cmellojr/modo-locadora/internal/auth"
	"github.com/cmellojr/modo-locadora/internal/database"
	"github.com/cmellojr/modo-locadora/internal/igdb"
	"github.com/cmellojr/modo-locadora/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Handler handles HTTP requests for the system.
type Handler struct {
	store        database.Store
	cookieSecret string
}

// NewHandler creates a new Handler with the provided store and cookie secret.
func NewHandler(store database.Store, cookieSecret string) *Handler {
	return &Handler{store: store, cookieSecret: cookieSecret}
}

// CreateMemberRequest defines the input for member registration.
type CreateMemberRequest struct {
	ProfileName     string `json:"profile_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	FavoriteConsole string `json:"favorite_console"`
}

// CreateMember handles POST /members.
func (h *Handler) CreateMember(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	member := &models.Member{
		ID:              uuid.New(),
		ProfileName:     req.ProfileName,
		Email:           req.Email,
		PasswordHash:    string(hashedPassword),
		FavoriteConsole: req.FavoriteConsole,
		JoinedAt:        time.Now(),
	}

	if err := h.store.CreateMember(r.Context(), member); err != nil {
		http.Error(w, "Failed to create member", http.StatusInternalServerError)
		return
	}

	// Do not expose the password hash in the response.
	member.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)
}

// Login handles POST /login, validating credentials and setting a signed cookie.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	profileName := r.FormValue("profile_name")
	password := r.FormValue("password")

	if profileName == "" || password == "" {
		http.Error(w, "Nome e senha são obrigatórios", http.StatusBadRequest)
		return
	}

	// If store is not available, deny login (no way to validate).
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	member, err := h.store.GetMemberByProfileName(r.Context(), profileName)
	if err != nil {
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	if member == nil {
		http.Error(w, "Nome ou senha inválidos", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(member.PasswordHash), []byte(password)); err != nil {
		http.Error(w, "Nome ou senha inválidos", http.StatusUnauthorized)
		return
	}

	auth.SetSessionCookie(w, member.ID.String(), h.cookieSecret)
	http.Redirect(w, r, "/games", http.StatusSeeOther)
}

// GameView represents a game for display in the shelf.
type GameView struct {
	ID              string
	Title           string
	Platform        string
	CoverURL        string
	CopiesAvailable int
}

// ListGames handles GET /games and renders the games shelf.
func (h *Handler) ListGames(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	memberName := "Visitante"

	memberID := auth.GetSessionMemberID(r, h.cookieSecret)
	if memberID != "" && h.store != nil {
		id, err := uuid.Parse(memberID)
		if err == nil {
			member, err := h.store.GetMemberByID(r.Context(), id)
			if err == nil && member != nil {
				memberName = member.ProfileName
			}
		}
	}

	var dbGames []models.Game
	if h.store != nil {
		dbGames, _ = h.store.ListGames(r.Context())
	}

	// For now, let's mix mock data with DB data if DB is empty
	var releases []GameView
	var catalog []GameView

	if len(dbGames) > 0 {
		for _, g := range dbGames {
			view := GameView{
				ID:              g.ID.String(),
				Title:           g.Title,
				Platform:        g.Platform,
				CoverURL:        g.CoverURL,
				CopiesAvailable: 1, // Defaulting to 1 for now until copies are implemented
			}
			// Logic to separate: let's say last 2 are releases
			if len(releases) < 2 {
				releases = append(releases, view)
			} else {
				catalog = append(catalog, view)
			}
		}
	} else {
		// Mock data for testing
		releases = []GameView{
			{
				ID:              "00000000-0000-0000-0000-000000000001",
				Title:           "Chrono Trigger",
				Platform:        "SNES",
				CoverURL:        "https://images.igdb.com/igdb/image/upload/t_cover_big/co1v9x.jpg",
				CopiesAvailable: 1,
			},
		}

		catalog = []GameView{
			{
				ID:              "00000000-0000-0000-0000-000000000002",
				Title:           "Top Gear",
				Platform:        "SNES",
				CoverURL:        "https://images.igdb.com/igdb/image/upload/t_cover_big/co2607.jpg",
				CopiesAvailable: 0,
			},
			{
				ID:              "00000000-0000-0000-0000-000000000003",
				Title:           "Super Metroid",
				Platform:        "SNES",
				CoverURL:        "https://images.igdb.com/igdb/image/upload/t_cover_big/co1tpz.jpg",
				CopiesAvailable: 2,
			},
		}
	}

	data := struct {
		MemberName string
		Releases   []GameView
		Catalog    []GameView
	}{
		MemberName: memberName,
		Releases:   releases,
		Catalog:    catalog,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AdminStock handles GET /admin/stock and renders the IGDB search page.
func (h *Handler) AdminStock(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	query := r.URL.Query().Get("q")
	magazine := r.URL.Query().Get("magazine")

	var results []igdb.GameData
	if query != "" {
		clientID := os.Getenv("TWITCH_CLIENT_ID")
		clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

		if clientID != "" && clientSecret != "" {
			token, err := igdb.GetAccessToken(clientID, clientSecret)
			if err == nil {
				results, _ = igdb.SearchGame(clientID, token.AccessToken, query)
			}
		}
	}

	data := struct {
		Query    string
		Magazine string
		Results  []igdb.GameData
	}{
		Query:    query,
		Magazine: magazine,
		Results:  results,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PurchaseGame handles POST /admin/purchase.
func (h *Handler) PurchaseGame(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	game := &models.Game{
		ID:             uuid.New(),
		Title:          r.FormValue("title"),
		IgdbID:         r.FormValue("igdb_id"),
		Platform:       "SNES", // Default for now
		Summary:        r.FormValue("summary"),
		CoverURL:       r.FormValue("cover_url"),
		SourceMagazine: r.FormValue("magazine"),
		AcquiredAt:     time.Now(),
	}

	if err := h.store.AddGame(r.Context(), game); err != nil {
		http.Error(w, "Failed to purchase game: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/games", http.StatusSeeOther)
}

// SearchGame handles GET /search?q=... and returns raw JSON from IGDB.
func (h *Handler) SearchGame(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		http.Error(w, "IGDB credentials not configured", http.StatusServiceUnavailable)
		return
	}

	token, err := igdb.GetAccessToken(clientID, clientSecret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get IGDB token: %v", err), http.StatusInternalServerError)
		return
	}

	games, err := igdb.SearchGame(clientID, token.AccessToken, query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to search IGDB: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

// GetGame handles GET /games/{id}.
func (h *Handler) GetGame(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Assuming the ID is passed as a query parameter "id" for now,
	// as stdlib's default ServeMux doesn't support path parameters easily before Go 1.22.
	// Since we are on Go 1.24.3, we can use the new path parameter syntax if we want.
	idStr := r.PathValue("id")
	if idStr == "" {
		idStr = r.URL.Query().Get("id")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	game, err := h.store.GetGameByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve game", http.StatusInternalServerError)
		return
	}

	if game == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}
