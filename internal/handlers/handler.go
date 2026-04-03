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
	adminEmail   string
}

// NewHandler creates a new Handler with the provided store, cookie secret and admin email.
func NewHandler(store database.Store, cookieSecret, adminEmail string) *Handler {
	return &Handler{store: store, cookieSecret: cookieSecret, adminEmail: adminEmail}
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
	case "verdict_completed":
		return fmt.Sprintf("%s detonou %s! Zerou com estilo!", a.MemberName, a.GameTitle)
	case "verdict_enjoyed":
		return fmt.Sprintf("%s aproveitou bastante %s!", a.MemberName, a.GameTitle)
	case "verdict_quick_play":
		return fmt.Sprintf("%s deu uma partidinha em %s.", a.MemberName, a.GameTitle)
	case "verdict_not_for_me":
		return fmt.Sprintf("%s tentou %s mas nao deu.", a.MemberName, a.GameTitle)
	case "verdict_gave_up":
		return fmt.Sprintf("%s desistiu de %s. O Tio esta decepcionado!", a.MemberName, a.GameTitle)
	case "relic":
		return fmt.Sprintf("%s virou Reliquia da Casa! 10 socios ja detonaram!", a.GameTitle)
	case "club_created":
		return fmt.Sprintf("%s formou a turma %s! Quem vai entrar?", a.MemberName, a.GameTitle)
	case "club_joined":
		return fmt.Sprintf("%s entrou na turma %s!", a.MemberName, a.GameTitle)
	default:
		return ""
	}
}

// LayoutData holds shared data for the layout template (top nav + sidebars).
// Page-specific handler structs embed this to make all fields available at the top level.
type LayoutData struct {
	PageTitle    string
	IsLoggedIn   bool
	IsAdmin      bool
	MemberName   string
	MemberMini   *MemberMiniView
	ShameEntries []database.ShameEntry
	Activities   []ActivityView
	AlmanacEntry string
}

// buildLayoutData loads shared layout data (auth state, sidebar content) for every page.
func (h *Handler) buildLayoutData(r *http.Request, pageTitle string) LayoutData {
	ld := LayoutData{PageTitle: pageTitle}

	if h.store == nil {
		return ld
	}

	// Authenticate
	rawID := auth.GetSessionMemberID(r, h.cookieSecret)
	if rawID != "" {
		id, err := uuid.Parse(rawID)
		if err == nil {
			member, err := h.store.GetMemberByID(r.Context(), id)
			if err == nil && member != nil {
				ld.IsLoggedIn = true
				ld.MemberName = member.ProfileName
				ld.IsAdmin = h.adminEmail != "" && member.Email == h.adminEmail

				// Left sidebar: member mini-card
				activeCount, overdueCount, _ := h.store.GetMemberRentalStats(r.Context(), id)
				statusLabel := "Ativo"
				statusBadge := "is-success"
				if member.Status == models.MemberStatusInDebt || overdueCount > 0 {
					statusLabel = "Em Debito"
					statusBadge = "is-error"
				} else if activeCount > 0 {
					statusLabel = "Jogando"
					statusBadge = "is-primary"
				}
				ld.MemberMini = &MemberMiniView{
					ProfileName:      member.ProfileName,
					MembershipNumber: member.MembershipNumber,
					StatusLabel:      statusLabel,
					StatusBadge:      statusBadge,
					ActiveRentals:    activeCount,
				}
			}
		}
	}

	// Right sidebar: activities, shame, almanac
	activities, _ := h.store.ListRecentActivities(r.Context(), 5)
	for _, a := range activities {
		ld.Activities = append(ld.Activities, ActivityView{
			EventType:  a.EventType,
			MemberName: a.MemberName,
			GameTitle:  a.GameTitle,
			Message:    formatActivityMessage(a),
			TimeAgo:    formatTimeAgo(a.CreatedAt),
		})
	}
	ld.ShameEntries, _ = h.store.GetTopShameEntries(r.Context(), 5)
	ld.AlmanacEntry = almanac.TodaysEphemeride()

	return ld
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
		http.Error(w, "Name and password are required", http.StatusBadRequest)
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
		http.Error(w, "Invalid name or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(member.PasswordHash), []byte(password)); err != nil {
		http.Error(w, "Invalid name or password", http.StatusUnauthorized)
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
	aliases := map[string]string{
		"Super Nintendo": "snes",
	}
	if alias, ok := aliases[platform]; ok {
		return "/static/img/logos/" + alias + ".svg"
	}
	name := strings.ToLower(strings.ReplaceAll(platform, " ", "-"))
	return "/static/img/logos/" + name + ".svg"
}

// GameView represents a game for display in the shelf.
type GameView struct {
	ID              string
	Title           string
	Platform        string
	CoverURL        string
	CoverDisplay    string
	Summary         string
	SourceMagazine  string
	TotalCopies     int
	AvailableCopies int
	RenterName      string
}

// ListGames handles GET /games. Without ?platform= it shows the platform selection page.
// With ?platform=X it shows the filtered games shelf for that platform.
func (h *Handler) ListGames(w http.ResponseWriter, r *http.Request, platformsTmpl, gamesTmpl *template.Template) {
	platform := r.URL.Query().Get("platform")

	// No platform filter → show platform selection page.
	if platform == "" {
		ld := h.buildLayoutData(r, "Acervo de Cartuchos")

		fixedPlatforms := []string{"Mega Drive", "Super Nintendo", "NES", "Master System", "Atari 2600"}
		platformCounts := make(map[string]int)
		if h.store != nil {
			platforms, _ := h.store.ListPlatforms(r.Context())
			for _, ps := range platforms {
				platformCounts[ps.Platform] = ps.GameCount
			}
		}
		var platformViews []PlatformView
		for _, p := range fixedPlatforms {
			platformViews = append(platformViews, PlatformView{
				Platform:  p,
				GameCount: platformCounts[p],
				LogoURL:   platformLogoFile(p),
			})
		}

		data := struct {
			LayoutData
			Platforms []PlatformView
		}{
			LayoutData: ld,
			Platforms:  platformViews,
		}

		if err := platformsTmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Platform filter present → show games for that platform.
	ld := h.buildLayoutData(r, platform)

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
					CoverDisplay:    ga.Game.CoverDisplay,
					TotalCopies:     ga.TotalCopies,
					AvailableCopies: ga.AvailableCopies,
					RenterName:      ga.RenterName,
				})
			}
		}
	}

	debtError := r.URL.Query().Get("error") == "in_debt"

	completedGames := make(map[string]bool)
	if ld.IsLoggedIn && h.store != nil {
		rawID := auth.GetSessionMemberID(r, h.cookieSecret)
		memberUUID, err := uuid.Parse(rawID)
		if err == nil {
			completedIDs, _ := h.store.ListCompletedGameIDs(r.Context(), memberUUID)
			for _, cid := range completedIDs {
				completedGames[cid.String()] = true
			}
		}
	}

	data := struct {
		LayoutData
		Games          []GameView
		Platform       string
		DebtError      bool
		CompletedGames map[string]bool
	}{
		LayoutData:     ld,
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
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	detail, err := h.store.GetGameDetail(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to load game: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if detail == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	ld := h.buildLayoutData(r, detail.Game.Title)

	data := struct {
		LayoutData
		Detail    *database.GameDetail
		DebtError bool
	}{
		LayoutData: ld,
		Detail:     detail,
		DebtError:  r.URL.Query().Get("error") == "in_debt",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Membership handles GET /membership and renders the member's profile card.
func (h *Handler) Membership(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
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

	ld := h.buildLayoutData(r, "Membership Card")

	activeCount, overdueCount, _ := h.store.GetMemberRentalStats(r.Context(), id)
	isInDebt := member.Status == models.MemberStatusInDebt

	statusLabel := "Jogador Honesto"
	statusBadge := "is-success"
	if isInDebt || overdueCount > 0 {
		statusLabel = "Socio em Debito com o Tio"
		statusBadge = "is-error"
	} else if activeCount > 0 {
		statusLabel = "Jogador Ativo"
		statusBadge = "is-primary"
	}

	memberRentals, _ := h.store.ListMemberActiveRentals(r.Context(), id)
	onTimeCount, _ := h.store.CountOnTimeReturns(r.Context(), id)
	completedGameIDs, _ := h.store.ListCompletedGameIDs(r.Context(), id)
	memberTitle := models.ComputeMemberTitle(len(completedGameIDs), onTimeCount)
	memberClubs, _ := h.store.ListMemberClubs(r.Context(), id)

	data := struct {
		LayoutData
		Member        *models.Member
		ActiveRentals int
		OverdueCount  int
		StatusLabel   string
		StatusBadge   string
		Success       string
		IsInDebt    bool
		LateCount     int
		Rentals       []database.MemberRental
		Title         models.MemberTitle
		Clubs         []database.MemberClubView
	}{
		LayoutData:    ld,
		Member:        member,
		ActiveRentals: activeCount,
		OverdueCount:  overdueCount,
		StatusLabel:   statusLabel,
		StatusBadge:   statusBadge,
		Success:       r.URL.Query().Get("success"),
		IsInDebt:    isInDebt,
		LateCount:     member.LateCount,
		Rentals:       memberRentals,
		Title:         memberTitle,
		Clubs:         memberClubs,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SavePasswordNotes handles POST /membership/notes.
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
		http.Error(w, "Failed to save notes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/membership?success=1", http.StatusSeeOther)
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
		http.Error(w, "Invalid session", http.StatusBadRequest)
		return
	}

	// Block rental if member is in debt.
	status, err := h.store.GetMemberStatus(r.Context(), memberID)
	if err != nil {
		http.Error(w, "Failed to check status: "+err.Error(), http.StatusInternalServerError)
		return
	}
	gameIDStr := r.FormValue("game_id")

	if status == models.MemberStatusInDebt {
		http.Redirect(w, r, "/games/"+gameIDStr+"?error=in_debt", http.StatusSeeOther)
		return
	}

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	if err := h.store.RentGame(r.Context(), gameID, memberID); err != nil {
		http.Error(w, "Failed to rent: "+err.Error(), http.StatusInternalServerError)
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

	ld := h.buildLayoutData(r, "Returns")

	rentals, err := h.store.ListActiveRentals(r.Context())
	if err != nil {
		http.Error(w, "Failed to list rentals: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		LayoutData
		Rentals []database.ActiveRental
		Success string
	}{
		LayoutData: ld,
		Rentals:    rentals,
		Success:    r.URL.Query().Get("success"),
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
		http.Error(w, "Invalid rental ID", http.StatusBadRequest)
		return
	}

	if err := h.store.ReturnGame(r.Context(), rentalID); err != nil {
		http.Error(w, "Failed to return: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/returns?success=Game+returned", http.StatusSeeOther)
}

// AdminStock handles GET /admin/stock and renders the IGDB search page.
func (h *Handler) AdminStock(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	ld := h.buildLayoutData(r, "Abastecer Prateleiras")

	query := r.URL.Query().Get("q")
	magazine := r.URL.Query().Get("magazine")
	selectedIDStr := r.URL.Query().Get("selected")

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

	data := struct {
		LayoutData
		Query    string
		Magazine string
		Results  []igdb.GameData
		Selected *igdb.GameData
		Success  string
	}{
		LayoutData: ld,
		Query:      query,
		Magazine:   magazine,
		Results:    results,
		Selected:   selected,
		Success:    r.URL.Query().Get("success"),
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

	ld := h.buildLayoutData(r, "Acervo")

	items, err := h.store.ListGamesWithPopularity(r.Context())
	if err != nil {
		http.Error(w, "Failed to list games: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		LayoutData
		Items   []database.GameInventoryItem
		Success string
	}{
		LayoutData: ld,
		Items:      items,
		Success:    r.URL.Query().Get("success"),
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

	ld := h.buildLayoutData(r, "Edit "+game.Title)

	rentalHistory, _ := h.store.ListGameRentalHistory(r.Context(), id, 5)

	data := struct {
		LayoutData
		Game          *models.Game
		RentalHistory []database.GameRentalHistoryEntry
	}{
		LayoutData:    ld,
		Game:          game,
		RentalHistory: rentalHistory,
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
		http.Error(w, "Failed to process form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	game, err := h.store.GetGameByID(r.Context(), id)
	if err != nil || game == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	game.Title = r.FormValue("title")
	game.Platform = r.FormValue("platform")
	game.Summary = r.FormValue("summary")
	game.SourceMagazine = r.FormValue("magazine")
	game.CoverDisplay = r.FormValue("cover_display")
	if game.CoverDisplay == "" {
		game.CoverDisplay = "cover"
	}

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
			http.Error(w, "Failed to save cover: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to write cover: "+err.Error(), http.StatusInternalServerError)
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

// HandleIndex handles GET / and renders the landing page.
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	ld := h.buildLayoutData(r, "Welcome")

	data := struct {
		LayoutData
	}{
		LayoutData: ld,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleRedeem handles POST /membership/redeem, clearing the member's debt status.
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
		http.Error(w, "Failed to redeem member: "+err.Error(), http.StatusInternalServerError)
		return
	}

	member, err := h.store.GetMemberByID(r.Context(), memberID)
	if err == nil && member != nil {
		_ = h.store.InsertActivity(r.Context(), "redemption", member.ProfileName, "")
	}

	http.Redirect(w, r, "/membership?success=redeemed", http.StatusSeeOther)
}

// HandleMemberReturn handles POST /membership/return, allowing a member to self-return a game.
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
		http.Error(w, "Invalid rental ID", http.StatusBadRequest)
		return
	}

	verdict := r.FormValue("verdict")
	validVerdicts := map[string]bool{
		"completed": true, "enjoyed": true, "quick_play": true,
		"not_for_me": true, "gave_up": true,
	}
	if !validVerdicts[verdict] {
		verdict = ""
	}

	// Get game title before the return (for activity logging).
	gameTitle, _ := h.store.GetRentalGameTitle(r.Context(), rentalID)

	if err := h.store.ReturnGameByMember(r.Context(), rentalID, memberID, verdict); err != nil {
		http.Error(w, "Failed to return: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fire verdict activity event.
	if verdict != "" && gameTitle != "" {
		member, _ := h.store.GetMemberByID(r.Context(), memberID)
		if member != nil {
			_ = h.store.InsertActivity(r.Context(), "verdict_"+verdict, member.ProfileName, gameTitle)
		}

		// Check if game just became a "Reliquia da Casa" (10+ completions).
		if verdict == "completed" {
			completionCount, err := h.store.CountGameCompletions(r.Context(), rentalID)
			if err == nil && completionCount == 10 {
				_ = h.store.InsertActivity(r.Context(), "relic", "", gameTitle)
			}
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

	http.Redirect(w, r, "/membership?success=returned", http.StatusSeeOther)
}

// ── Club handlers ───────────────────────────────────────────────────────────

// getSessionMemberID extracts and parses the member UUID from the session cookie.
func (h *Handler) getSessionMemberID(r *http.Request) (uuid.UUID, bool) {
	raw := auth.GetSessionMemberID(r, h.cookieSecret)
	if raw == "" {
		return uuid.UUID{}, false
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.UUID{}, false
	}
	return id, true
}

// requireClubAdmin verifies the session user is an admin of the club identified by {id} in the path.
// Returns the member ID, club ID, and true if authorized. Writes an error response and returns false otherwise.
func (h *Handler) requireClubAdmin(w http.ResponseWriter, r *http.Request) (uuid.UUID, uuid.UUID, bool) {
	memberID, ok := h.getSessionMemberID(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return uuid.UUID{}, uuid.UUID{}, false
	}

	clubID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid club ID", http.StatusBadRequest)
		return uuid.UUID{}, uuid.UUID{}, false
	}

	role, err := h.store.GetClubMemberRole(r.Context(), clubID, memberID)
	if err != nil || role != models.ClubRoleAdmin {
		http.Error(w, "Restricted to club admins", http.StatusForbidden)
		return uuid.UUID{}, uuid.UUID{}, false
	}

	return memberID, clubID, true
}

// ListClubs handles GET /clubs.
func (h *Handler) ListClubs(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	ld := h.buildLayoutData(r, "Clubs")

	var viewerID *uuid.UUID
	if id, ok := h.getSessionMemberID(r); ok {
		viewerID = &id
	}

	var clubs []database.ClubListItem
	if h.store != nil {
		clubs, _ = h.store.ListClubs(r.Context(), viewerID)
	}

	data := struct {
		LayoutData
		Clubs   []database.ClubListItem
		Success string
	}{
		LayoutData: ld,
		Clubs:      clubs,
		Success:    r.URL.Query().Get("success"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ClubDetail handles GET /clubs/{id}.
func (h *Handler) ClubDetail(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	clubID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid club ID", http.StatusBadRequest)
		return
	}

	detail, err := h.store.GetClubDetail(r.Context(), clubID)
	if err != nil {
		http.Error(w, "Failed to load club: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if detail == nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	ld := h.buildLayoutData(r, detail.Club.Name)

	var viewerRole string
	var viewerID uuid.UUID
	if id, ok := h.getSessionMemberID(r); ok {
		viewerID = id
		viewerRole, _ = h.store.GetClubMemberRole(r.Context(), clubID, id)
	}

	data := struct {
		LayoutData
		Detail      *database.ClubDetail
		ViewerRole  string
		IsMember    bool
		IsClubAdmin bool
		IsCreator   bool
		Success     string
	}{
		LayoutData:  ld,
		Detail:      detail,
		ViewerRole:  viewerRole,
		IsMember:    viewerRole != "",
		IsClubAdmin: viewerRole == models.ClubRoleAdmin,
		IsCreator:   viewerID == detail.Club.CreatedBy,
		Success:     r.URL.Query().Get("success"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ClubFormPage handles GET /clubs/new and GET /clubs/{id}/edit.
func (h *Handler) ClubFormPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template, isEdit bool) {
	var club *models.Club
	if isEdit {
		_, clubID, ok := h.requireClubAdmin(w, r)
		if !ok {
			return
		}
		var err error
		club, err = h.store.GetClubByID(r.Context(), clubID)
		if err != nil || club == nil {
			http.Error(w, "Club not found", http.StatusNotFound)
			return
		}
	}

	title := "Create Club"
	if isEdit && club != nil {
		title = "Edit " + club.Name
	}
	ld := h.buildLayoutData(r, title)

	data := struct {
		LayoutData
		Club   *models.Club
		IsEdit bool
	}{
		LayoutData: ld,
		Club:       club,
		IsEdit:     isEdit,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateClub handles POST /clubs.
func (h *Handler) CreateClub(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	memberID, ok := h.getSessionMemberID(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to process form", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "Club name is required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	club := &models.Club{
		ID:          uuid.New(),
		Name:        name,
		Description: r.FormValue("description"),
		WebsiteURL:  r.FormValue("website_url"),
		CreatedBy:   memberID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Handle badge file upload.
	file, header, err := r.FormFile("badge_file")
	if err == nil {
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		if ext == "" {
			ext = ".png"
		}
		filename := club.ID.String() + ext
		savePath := filepath.Join("web", "static", "clubs", filename)
		dst, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "Failed to save badge: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to write badge: "+err.Error(), http.StatusInternalServerError)
			return
		}
		club.BadgeURL = "/static/clubs/" + filename
	}

	if err := h.store.CreateClub(r.Context(), club); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			http.Error(w, "A club with this name already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create club: "+err.Error(), http.StatusInternalServerError)
		return
	}

	member, _ := h.store.GetMemberByID(r.Context(), memberID)
	if member != nil {
		_ = h.store.InsertActivity(r.Context(), "club_created", member.ProfileName, club.Name)
	}

	http.Redirect(w, r, "/clubs/"+club.ID.String()+"?success=created", http.StatusSeeOther)
}

// UpdateClub handles POST /clubs/{id}/edit.
func (h *Handler) UpdateClub(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	_, clubID, ok := h.requireClubAdmin(w, r)
	if !ok {
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to process form", http.StatusBadRequest)
		return
	}

	club, err := h.store.GetClubByID(r.Context(), clubID)
	if err != nil || club == nil {
		http.Error(w, "Club not found", http.StatusNotFound)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "Club name is required", http.StatusBadRequest)
		return
	}

	club.Name = name
	club.Description = r.FormValue("description")
	club.WebsiteURL = r.FormValue("website_url")

	// Handle badge file upload.
	file, header, err := r.FormFile("badge_file")
	if err == nil {
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		if ext == "" {
			ext = ".png"
		}
		filename := club.ID.String() + ext
		savePath := filepath.Join("web", "static", "clubs", filename)
		dst, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "Failed to save badge: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to write badge: "+err.Error(), http.StatusInternalServerError)
			return
		}
		club.BadgeURL = "/static/clubs/" + filename
	}

	if err := h.store.UpdateClub(r.Context(), club); err != nil {
		http.Error(w, "Failed to update club: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/clubs/"+clubID.String()+"?success=updated", http.StatusSeeOther)
}

// JoinClub handles POST /clubs/{id}/join.
func (h *Handler) JoinClub(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	memberID, ok := h.getSessionMemberID(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	clubID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid club ID", http.StatusBadRequest)
		return
	}

	if err := h.store.JoinClub(r.Context(), clubID, memberID); err != nil {
		http.Error(w, "Failed to join club: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log activity.
	member, _ := h.store.GetMemberByID(r.Context(), memberID)
	club, _ := h.store.GetClubByID(r.Context(), clubID)
	if member != nil && club != nil {
		_ = h.store.InsertActivity(r.Context(), "club_joined", member.ProfileName, club.Name)
	}

	http.Redirect(w, r, "/clubs/"+clubID.String()+"?success=joined", http.StatusSeeOther)
}

// LeaveClub handles POST /clubs/{id}/leave.
func (h *Handler) LeaveClub(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	memberID, ok := h.getSessionMemberID(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	clubID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid club ID", http.StatusBadRequest)
		return
	}

	// Prevent last admin from leaving.
	role, _ := h.store.GetClubMemberRole(r.Context(), clubID, memberID)
	if role == models.ClubRoleAdmin {
		detail, _ := h.store.GetClubDetail(r.Context(), clubID)
		if detail != nil {
			adminCount := 0
			for _, m := range detail.Members {
				if m.Role == models.ClubRoleAdmin {
					adminCount++
				}
			}
			if adminCount <= 1 {
				http.Error(w, "You are the only admin. Promote another member before leaving.", http.StatusForbidden)
				return
			}
		}
	}

	if err := h.store.LeaveClub(r.Context(), clubID, memberID); err != nil {
		http.Error(w, "Failed to leave club: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/clubs?success=left", http.StatusSeeOther)
}

// PromoteClubMember handles POST /clubs/{id}/promote.
func (h *Handler) PromoteClubMember(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	_, clubID, ok := h.requireClubAdmin(w, r)
	if !ok {
		return
	}

	targetID, err := uuid.Parse(r.FormValue("member_id"))
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	if err := h.store.PromoteClubMember(r.Context(), clubID, targetID); err != nil {
		http.Error(w, "Failed to promote member: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/clubs/"+clubID.String()+"?success=promoted", http.StatusSeeOther)
}

// RemoveClubMember handles POST /clubs/{id}/remove.
func (h *Handler) RemoveClubMember(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	memberID, clubID, ok := h.requireClubAdmin(w, r)
	if !ok {
		return
	}

	targetID, err := uuid.Parse(r.FormValue("member_id"))
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	if targetID == memberID {
		http.Error(w, "Use 'Leave club' to remove yourself.", http.StatusBadRequest)
		return
	}

	if err := h.store.RemoveClubMember(r.Context(), clubID, targetID); err != nil {
		http.Error(w, "Failed to remove member: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/clubs/"+clubID.String()+"?success=removed", http.StatusSeeOther)
}

// DeleteClub handles POST /clubs/{id}/delete.
func (h *Handler) DeleteClub(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		http.Error(w, "Database not configured", http.StatusServiceUnavailable)
		return
	}

	memberID, ok := h.getSessionMemberID(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	clubID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid club ID", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteClub(r.Context(), clubID, memberID); err != nil {
		http.Error(w, "Failed to delete club: "+err.Error(), http.StatusForbidden)
		return
	}

	http.Redirect(w, r, "/clubs?success=deleted", http.StatusSeeOther)
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
