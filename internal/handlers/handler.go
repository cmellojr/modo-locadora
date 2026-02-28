package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/cmellojr/modo-locadora/internal/database"
	"github.com/cmellojr/modo-locadora/internal/models"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for the system.
type Handler struct {
	store database.Store
}

// NewHandler creates a new Handler with the provided store.
func NewHandler(store database.Store) *Handler {
	return &Handler{store: store}
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

	member := &models.Member{
		ID:              uuid.New(),
		ProfileName:     req.ProfileName,
		Email:           req.Email,
		PasswordHash:    req.Password, // TODO: Hash password properly
		FavoriteConsole: req.FavoriteConsole,
		JoinedAt:        time.Now(),
	}

	if err := h.store.CreateMember(r.Context(), member); err != nil {
		http.Error(w, "Failed to create member", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)
}

// GameView represents a game for display in the shelf.
type GameView struct {
	ID              string
	Title           string
	Platform        string
	CoverURL        string
	CopiesAvailable int
}

// ListGames handles GET /jogos and renders the games shelf.
func (h *Handler) ListGames(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// Mock data for testing
	games := []GameView{
		{
			ID:              "1",
			Title:           "Super Mario World",
			Platform:        "SNES",
			CoverURL:        "https://upload.wikimedia.org/wikipedia/en/3/32/Super_Mario_World_Coverart.png",
			CopiesAvailable: 2,
		},
		{
			ID:              "2",
			Title:           "The Legend of Zelda: A Link to the Past",
			Platform:        "SNES",
			CoverURL:        "https://upload.wikimedia.org/wikipedia/en/2/21/The_Legend_of_Zelda_A_Link_to_the_Past_SNES_Game_Cover.jpg",
			CopiesAvailable: 0,
		},
		{
			ID:              "3",
			Title:           "Sonic the Hedgehog",
			Platform:        "Mega Drive",
			CoverURL:        "https://upload.wikimedia.org/wikipedia/en/b/ba/Sonic_the_Hedgehog_1_Genesis_box_art.jpg",
			CopiesAvailable: 1,
		},
	}

	data := struct {
		Games []GameView
	}{
		Games: games,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
