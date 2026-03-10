package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cmellojr/modo-locadora/internal/almanac"
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

// getSessionMember returns the authenticated member's name and login status.
func (h *Handler) getSessionMember(r *http.Request) (name string, loggedIn bool) {
	rawID := auth.GetSessionMemberID(r, h.cookieSecret)
	if rawID == "" || h.store == nil {
		return "", false
	}
	id, err := uuid.Parse(rawID)
	if err != nil {
		return "", false
	}
	member, err := h.store.GetMemberByID(r.Context(), id)
	if err != nil || member == nil {
		return "", false
	}
	return member.ProfileName, true
}

// Logout handles POST /logout by clearing the session cookie.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_member",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ActivityView represents a formatted activity event for template display.
type ActivityView struct {
	EventType  string
	MemberName string
	GameTitle  string
	Message    string
	TimeAgo    string
}

// formatTimeAgo returns a Portuguese relative-time string.
func formatTimeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "agora"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 min atras"
		}
		return fmt.Sprintf("%d min atras", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hora atras"
		}
		return fmt.Sprintf("%d horas atras", h)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "ontem"
		}
		return fmt.Sprintf("%d dias atras", days)
	}
}

// MemberMiniView holds compact member info for sidebar display.
type MemberMiniView struct {
	ProfileName      string
	MembershipNumber string
	StatusLabel      string
	StatusBadge      string
	ActiveRentals    int
}

// formatActivityMessage returns a human-readable message for an activity event.
func formatActivityMessage(a database.ActivityEntry) string {
	switch a.EventType {
	case "penalty":
		return fmt.Sprintf("%s foi penalizado(a) por atrasar %s!", a.MemberName, a.GameTitle)
	case "redemption":
		return fmt.Sprintf("%s soprou o cartucho e foi redimido(a)!", a.MemberName)
	case "new_game":
		return fmt.Sprintf("Nova fita no acervo: %s!", a.GameTitle)
	case "prestige":
		return fmt.Sprintf("%s atingiu prestigio! Socio(a) exemplar!", a.MemberName)
	case "verdict_complete":
		return fmt.Sprintf("%s zerou %s! Respeito total!", a.MemberName, a.GameTitle)
	case "verdict_partial":
		return fmt.Sprintf("%s jogou %s mas nao terminou.", a.MemberName, a.GameTitle)
	case "verdict_quit":
		return fmt.Sprintf("%s desistiu de %s. O Tio esta decepcionado!", a.MemberName, a.GameTitle)
	default:
		return ""
	}
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

	membershipNumber, err := h.store.NextMembershipNumber(r.Context())
	if err != nil {
		http.Error(w, "Failed to generate membership number", http.StatusInternalServerError)
		return
	}

	member := &models.Member{
		ID:               uuid.New(),
		ProfileName:      req.ProfileName,
		Email:            req.Email,
		PasswordHash:     string(hashedPassword),
		FavoriteConsole:   req.FavoriteConsole,
		MembershipNumber: membershipNumber,
		JoinedAt:         time.Now(),
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

// PlatformView represents a platform for display in the console selection grid.
type PlatformView struct {
	Platform  string
	GameCount int
	LogoURL   string
}

// platformLogoFile returns the logo filename for a platform name.
func platformLogoFile(platform string) string {
	name := strings.ToLower(strings.ReplaceAll(platform, " ", "-"))
	return "/static/img/logos/" + name + ".svg"
}

// GameView represents a game for display in the shelf.
type GameView struct {
	ID              string
	Title           string
	Platform        string
	CoverURL        string
	Summary         string
	SourceMagazine  string
	TotalCopies     int
	AvailableCopies int
	RenterName      string
}

// ListGames handles GET /games. Without ?platform= it shows the platform selection page.
// With ?platform=X it shows the filtered games shelf for that platform.
func (h *Handler) ListGames(w http.ResponseWriter, r *http.Request, platformsTmpl, gamesTmpl *template.Template) {
	memberName := "Visitante"
	var memberID string
	var isLoggedIn bool

	rawMemberID := auth.GetSessionMemberID(r, h.cookieSecret)
	if rawMemberID != "" && h.store != nil {
		id, err := uuid.Parse(rawMemberID)
		if err == nil {
			member, err := h.store.GetMemberByID(r.Context(), id)
			if err == nil && member != nil {
				memberName = member.ProfileName
				memberID = rawMemberID
				isLoggedIn = true
			}
		}
	}

	platform := r.URL.Query().Get("platform")

	// No platform filter → show platform selection page.
	if platform == "" {
		var platformViews []PlatformView
		if h.store != nil {
			platforms, _ := h.store.ListPlatforms(r.Context())
			for _, ps := range platforms {
				platformViews = append(platformViews, PlatformView{
					Platform:  ps.Platform,
					GameCount: ps.GameCount,
					LogoURL:   platformLogoFile(ps.Platform),
				})
			}
		}

		var activityViews []ActivityView
		if h.store != nil {
			activities, _ := h.store.ListRecentActivities(r.Context(), 5)
			for _, a := range activities {
				av := ActivityView{
					EventType:  a.EventType,
					MemberName: a.MemberName,
					GameTitle:  a.GameTitle,
					Message:    formatActivityMessage(a),
					TimeAgo:    formatTimeAgo(a.CreatedAt),
				}
				activityViews = append(activityViews, av)
			}
		}

		// Build member mini-card for left sidebar.
		var memberMini *MemberMiniView
		if isLoggedIn && h.store != nil {
			id, err := uuid.Parse(rawMemberID)
			if err == nil {
				member, err := h.store.GetMemberByID(r.Context(), id)
				if err == nil && member != nil {
					activeCount, overdueCount, _ := h.store.GetMemberRentalStats(r.Context(), id)
					statusLabel := "Ativo"
					statusBadge := "is-success"
					if member.Status == models.MemberStatusEmDebito || overdueCount > 0 {
						statusLabel = "Em Debito"
						statusBadge = "is-error"
					} else if activeCount > 0 {
						statusLabel = "Jogando"
						statusBadge = "is-primary"
					}
					memberMini = &MemberMiniView{
						ProfileName:      member.ProfileName,
						MembershipNumber: member.MembershipNumber,
						StatusLabel:      statusLabel,
						StatusBadge:      statusBadge,
						ActiveRentals:    activeCount,
					}
				}
			}
		}

		// Fetch Wall of Shame for left sidebar.
		var shameEntries []database.ShameEntry
		if h.store != nil {
			shameEntries, _ = h.store.GetTopShameEntries(r.Context(), 5)
		}

		data := struct {
			MemberName   string
			IsLoggedIn   bool
			Platforms    []PlatformView
			Activities   []ActivityView
			Ephemeride   string
			MemberMini   *MemberMiniView
			ShameEntries []database.ShameEntry
		}{
			MemberName:   memberName,
			IsLoggedIn:   isLoggedIn,
			Platforms:    platformViews,
			Activities:   activityViews,
			Ephemeride:   almanac.TodaysEphemeride(),
			MemberMini:   memberMini,
			ShameEntries: shameEntries,
		}

		if err := platformsTmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Platform filter present → show games for that platform.
	var games []GameView

	if h.store != nil {
		gamesAvail, err := h.store.ListGamesWithAvailability(r.Context(), platform)
		if err == nil {
			for _, ga := range gamesAvail {
				games = append(games, GameView{
					ID:              ga.Game.ID.String(),
					Title:           ga.Game.Title,
					Platform:        ga.Game.Platform,
					CoverURL:        ga.Game.CoverURL,
					TotalCopies:     ga.TotalCopies,
					AvailableCopies: ga.AvailableCopies,
					RenterName:      ga.RenterName,
				})
			}
		}
	}

	debtError := r.URL.Query().Get("error") == "em_debito"

	// Build completed games map for golden star display.
	completedGames := make(map[string]bool)
	if isLoggedIn && h.store != nil {
		memberUUID, err := uuid.Parse(memberID)
		if err == nil {
			completedIDs, _ := h.store.ListCompletedGameIDs(r.Context(), memberUUID)
			for _, cid := range completedIDs {
				completedGames[cid.String()] = true
			}
		}
	}

	data := struct {
		MemberName     string
		MemberID       string
		IsLoggedIn     bool
		Games          []GameView
		Platform       string
		DebtError      bool
		CompletedGames map[string]bool
	}{
		MemberName:     memberName,
		MemberID:       memberID,
		IsLoggedIn:     isLoggedIn,
		Games:          games,
		Platform:       platform,
		DebtError:      debtError,
		CompletedGames: completedGames,
	}

	if err := gamesTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GameDetailPage handles GET /games/{id} and renders the game detail page.
func (h *Handler) GameDetailPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de jogo inválido", http.StatusBadRequest)
		return
	}

	detail, err := h.store.GetGameDetail(r.Context(), id)
	if err != nil {
		http.Error(w, "Falha ao carregar jogo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if detail == nil {
		http.Error(w, "Jogo não encontrado", http.StatusNotFound)
		return
	}

	memberName, isLoggedIn := h.getSessionMember(r)

	debtError := r.URL.Query().Get("error") == "em_debito"

	data := struct {
		Detail     *database.GameDetail
		MemberName string
		IsLoggedIn bool
		DebtError  bool
	}{
		Detail:     detail,
		MemberName: memberName,
		IsLoggedIn: isLoggedIn,
		DebtError:  debtError,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Carteirinha handles GET /carteirinha and renders the member's profile card.
func (h *Handler) Carteirinha(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	rawMemberID := auth.GetSessionMemberID(r, h.cookieSecret)
	if rawMemberID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, err := uuid.Parse(rawMemberID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	member, err := h.store.GetMemberByID(r.Context(), id)
	if err != nil || member == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	activeCount, overdueCount, _ := h.store.GetMemberRentalStats(r.Context(), id)

	isEmDebito := member.Status == models.MemberStatusEmDebito

	statusLabel := "Jogador Honesto"
	statusBadge := "is-success"
	if isEmDebito || overdueCount > 0 {
		statusLabel = "Socio em Debito com o Tio"
		statusBadge = "is-error"
	} else if activeCount > 0 {
		statusLabel = "Jogador Ativo"
		statusBadge = "is-primary"
	}

	successMsg := r.URL.Query().Get("success")

	var memberRentals []database.MemberRental
	memberRentals, _ = h.store.ListMemberActiveRentals(r.Context(), id)

	data := struct {
		Member        *models.Member
		ActiveRentals int
		OverdueCount  int
		StatusLabel   string
		StatusBadge   string
		Success       string
		IsEmDebito    bool
		LateCount     int
		Rentals       []database.MemberRental
	}{
		Member:        member,
		ActiveRentals: activeCount,
		OverdueCount:  overdueCount,
		StatusLabel:   statusLabel,
		StatusBadge:   statusBadge,
		Success:       successMsg,
		IsEmDebito:    isEmDebito,
		LateCount:     member.LateCount,
		Rentals:       memberRentals,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SavePasswordNotes handles POST /carteirinha/notes.
func (h *Handler) SavePasswordNotes(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	rawMemberID := auth.GetSessionMemberID(r, h.cookieSecret)
	if rawMemberID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	memberID, err := uuid.Parse(rawMemberID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	notes := r.FormValue("notes")
	if err := h.store.UpdateMemberNotes(r.Context(), memberID, notes); err != nil {
		http.Error(w, "Falha ao salvar notas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/carteirinha?success=1", http.StatusSeeOther)
}

// RentGame handles POST /rent.
func (h *Handler) RentGame(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	rawMemberID := auth.GetSessionMemberID(r, h.cookieSecret)
	if rawMemberID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	memberID, err := uuid.Parse(rawMemberID)
	if err != nil {
		http.Error(w, "Sessão inválida", http.StatusBadRequest)
		return
	}

	// Block rental if member is in debt.
	status, err := h.store.GetMemberStatus(r.Context(), memberID)
	if err != nil {
		http.Error(w, "Falha ao verificar status: "+err.Error(), http.StatusInternalServerError)
		return
	}
	gameIDStr := r.FormValue("game_id")

	if status == models.MemberStatusEmDebito {
		http.Redirect(w, r, "/games/"+gameIDStr+"?error=em_debito", http.StatusSeeOther)
		return
	}

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "ID de jogo inválido", http.StatusBadRequest)
		return
	}

	if err := h.store.RentGame(r.Context(), gameID, memberID); err != nil {
		http.Error(w, "Falha ao alugar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/games/"+gameID.String(), http.StatusSeeOther)
}

// AdminReturns handles GET /admin/returns and renders the active rentals for check-in.
func (h *Handler) AdminReturns(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	memberName, _ := h.getSessionMember(r)
	successMsg := r.URL.Query().Get("success")

	rentals, err := h.store.ListActiveRentals(r.Context())
	if err != nil {
		http.Error(w, "Failed to list rentals: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		MemberName string
		Rentals    []database.ActiveRental
		Success    string
	}{
		MemberName: memberName,
		Rentals:    rentals,
		Success:    successMsg,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ReturnGame handles POST /admin/return-game.
func (h *Handler) ReturnGame(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	rentalID, err := uuid.Parse(r.FormValue("rental_id"))
	if err != nil {
		http.Error(w, "ID de aluguel inválido", http.StatusBadRequest)
		return
	}

	if err := h.store.ReturnGame(r.Context(), rentalID); err != nil {
		http.Error(w, "Falha ao devolver: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/returns?success=Fita+devolvida", http.StatusSeeOther)
}

// AdminStock handles GET /admin/stock and renders the IGDB search page.
func (h *Handler) AdminStock(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	query := r.URL.Query().Get("q")
	magazine := r.URL.Query().Get("magazine")
	selectedIDStr := r.URL.Query().Get("selected")
	successMsg := r.URL.Query().Get("success")

	var results []igdb.GameData
	var selected *igdb.GameData

	if query != "" {
		clientID := os.Getenv("TWITCH_CLIENT_ID")
		clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

		if clientID != "" && clientSecret != "" {
			token, err := igdb.GetAccessToken(clientID, clientSecret)
			if err == nil {
				results, _ = igdb.SearchGame(clientID, token.AccessToken, query)
			}
		}

		if selectedIDStr != "" {
			selectedID := 0
			fmt.Sscanf(selectedIDStr, "%d", &selectedID)
			for i := range results {
				if results[i].ID == selectedID {
					selected = &results[i]
					break
				}
			}
		}
	}

	memberName, _ := h.getSessionMember(r)

	data := struct {
		MemberName string
		Query      string
		Magazine   string
		Results    []igdb.GameData
		Selected   *igdb.GameData
		Success    string
	}{
		MemberName: memberName,
		Query:      query,
		Magazine:   magazine,
		Results:    results,
		Selected:   selected,
		Success:    successMsg,
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

	platform := r.FormValue("platform")
	if platform == "" {
		platform = "N/A"
	}

	coverURL := r.FormValue("cover_url")
	if strings.Contains(coverURL, "t_thumb") {
		coverURL = strings.Replace(coverURL, "t_thumb", "t_cover_big", 1)
	}

	game := &models.Game{
		ID:             uuid.New(),
		Title:          r.FormValue("title"),
		IgdbID:         r.FormValue("igdb_id"),
		Platform:       platform,
		Summary:        r.FormValue("summary"),
		CoverURL:       coverURL,
		SourceMagazine: r.FormValue("magazine"),
		AcquiredAt:     time.Now(),
	}

	if err := h.store.AddGame(r.Context(), game); err != nil {
		http.Error(w, "Failed to purchase game: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_ = h.store.InsertActivity(r.Context(), "new_game", "", game.Title)

	http.Redirect(w, r, "/admin/edit/"+game.ID.String(), http.StatusSeeOther)
}

// AdminInventory handles GET /admin/inventory and renders the game list table.
func (h *Handler) AdminInventory(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	memberName, _ := h.getSessionMember(r)
	successMsg := r.URL.Query().Get("success")

	games, err := h.store.ListGames(r.Context())
	if err != nil {
		http.Error(w, "Failed to list games: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		MemberName string
		Games      []models.Game
		Success    string
	}{
		MemberName: memberName,
		Games:      games,
		Success:    successMsg,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// EditGame handles GET /admin/edit/{id} and renders the edit form for a game.
func (h *Handler) EditGame(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de jogo inválido", http.StatusBadRequest)
		return
	}

	game, err := h.store.GetGameByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve game", http.StatusInternalServerError)
		return
	}
	if game == nil {
		http.Error(w, "Jogo não encontrado", http.StatusNotFound)
		return
	}

	memberName, _ := h.getSessionMember(r)

	data := struct {
		MemberName string
		Game       *models.Game
	}{
		MemberName: memberName,
		Game:       game,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateGame handles POST /admin/update-game and updates game fields in the database.
func (h *Handler) UpdateGame(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Falha ao processar formulário", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID de jogo inválido", http.StatusBadRequest)
		return
	}

	game, err := h.store.GetGameByID(r.Context(), id)
	if err != nil || game == nil {
		http.Error(w, "Jogo não encontrado", http.StatusNotFound)
		return
	}

	game.Title = r.FormValue("title")
	game.Platform = r.FormValue("platform")
	game.Summary = r.FormValue("summary")
	game.SourceMagazine = r.FormValue("magazine")

	// Handle cover file upload.
	file, header, err := r.FormFile("cover_file")
	if err == nil {
		defer file.Close()

		ext := filepath.Ext(header.Filename)
		if ext == "" {
			ext = ".jpg"
		}
		filename := id.String() + ext
		savePath := filepath.Join("web", "static", "covers", filename)

		dst, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "Falha ao salvar capa: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Falha ao gravar capa: "+err.Error(), http.StatusInternalServerError)
			return
		}

		game.CoverURL = "/static/covers/" + filename
	} else {
		// No upload — preserve existing cover_url from hidden field.
		coverURL := r.FormValue("cover_url")
		if coverURL != "" {
			if strings.Contains(coverURL, "t_thumb") {
				coverURL = strings.Replace(coverURL, "t_thumb", "t_cover_big", 1)
			}
			game.CoverURL = coverURL
		}
	}

	if err := h.store.UpdateGame(r.Context(), game); err != nil {
		http.Error(w, "Failed to update game: "+err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("/admin/inventory?success=%s",
		template.URLQueryEscaper(game.Title))
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
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

// HandleIndex handles GET / and renders the landing page with the Wall of Shame.
// Authenticated members are redirected straight to the shelf.
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	memberName, isLoggedIn := h.getSessionMember(r)

	var shameEntries []database.ShameEntry
	if h.store != nil {
		entries, err := h.store.GetTopShameEntries(r.Context(), 5)
		if err == nil {
			shameEntries = entries
		}
	}

	data := struct {
		MemberName   string
		IsLoggedIn   bool
		ShameEntries []database.ShameEntry
	}{
		MemberName:   memberName,
		IsLoggedIn:   isLoggedIn,
		ShameEntries: shameEntries,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleRedeem handles POST /carteirinha/redeem, clearing the member's debt status.
func (h *Handler) HandleRedeem(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	rawMemberID := auth.GetSessionMemberID(r, h.cookieSecret)
	if rawMemberID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	memberID, err := uuid.Parse(rawMemberID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := h.store.RedeemMember(r.Context(), memberID); err != nil {
		http.Error(w, "Falha na redenção: "+err.Error(), http.StatusInternalServerError)
		return
	}

	member, err := h.store.GetMemberByID(r.Context(), memberID)
	if err == nil && member != nil {
		_ = h.store.InsertActivity(r.Context(), "redemption", member.ProfileName, "")
	}

	http.Redirect(w, r, "/carteirinha?success=redencao", http.StatusSeeOther)
}

// HandleMemberReturn handles POST /carteirinha/return, allowing a member to self-return a game.
func (h *Handler) HandleMemberReturn(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	rawMemberID := auth.GetSessionMemberID(r, h.cookieSecret)
	if rawMemberID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	memberID, err := uuid.Parse(rawMemberID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rentalID, err := uuid.Parse(r.FormValue("rental_id"))
	if err != nil {
		http.Error(w, "ID de aluguel invalido", http.StatusBadRequest)
		return
	}

	verdict := r.FormValue("verdict")
	if verdict != "zerei" && verdict != "joguei_um_pouco" && verdict != "desisti" {
		verdict = ""
	}

	// Get game title before the return (for activity logging).
	gameTitle, _ := h.store.GetRentalGameTitle(r.Context(), rentalID)

	if err := h.store.ReturnGameByMember(r.Context(), rentalID, memberID, verdict); err != nil {
		http.Error(w, "Falha ao devolver: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fire verdict activity event.
	if verdict != "" && gameTitle != "" {
		member, _ := h.store.GetMemberByID(r.Context(), memberID)
		if member != nil {
			var eventType string
			switch verdict {
			case "zerei":
				eventType = "verdict_complete"
			case "joguei_um_pouco":
				eventType = "verdict_partial"
			case "desisti":
				eventType = "verdict_quit"
			}
			_ = h.store.InsertActivity(r.Context(), eventType, member.ProfileName, gameTitle)
		}
	}

	// Check for prestige milestone (every 10th on-time return).
	onTimeCount, err := h.store.CountOnTimeReturns(r.Context(), memberID)
	if err == nil && onTimeCount > 0 && onTimeCount%10 == 0 {
		member, err := h.store.GetMemberByID(r.Context(), memberID)
		if err == nil && member != nil {
			_ = h.store.InsertActivity(r.Context(), "prestige", member.ProfileName, "")
		}
	}

	http.Redirect(w, r, "/carteirinha?success=devolucao", http.StatusSeeOther)
}

// GetGame handles GET /games/{id}.
func (h *Handler) GetGame(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

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
