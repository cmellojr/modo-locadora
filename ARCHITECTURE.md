# Modo Locadora - System Architecture 🕹️

## 1. Vision & Concept
Modo Locadora is a retro-gaming session manager focused on scarcity, community, and nostalgia. 
It emulates the experience of 90s Brazilian video rental stores ("locadoras").

## 2. Tech Stack
- **Backend:** Go (Golang) 1.22+
- **Database:** PostgreSQL (using `pgx` driver)
- **Frontend:** Server-Side Rendering (SSR) with `html/template`.
- **Styling:** [NES.css](https://nostalgic-css.github.io/NES.css/) + Google Fonts ("Press Start 2P").
- **External API:** IGDB (Internet Game Database) for game metadata.

## 3. Design Principles (The "Nostalgic" Rulebook)
- **Clean & Static:** No heavy JavaScript, no ads, no trackers.
- **Visuals:** Use `image-rendering: pixelated` for game covers. 
- **Language:** Code, Routes, and Database must be in **English**. The UI (templates) must be in **Portuguese (BR)**.
- **Copyleft:** Licensed under **GPL v3**.

## 4. Key Entities
- **Member (User):** Represents a store member with a "Trust Score".
- **Game:** Metadata fetched from IGDB.
- **Tape (Instance):** Each game has limited physical-like instances (Scarcity).
- **Notebook (Caderninho):** Private notes for the player.
- **Legacy (Verso da Fita):** Public notes left for the next player.

## 5. Routing Strategy
- `/` -> Landing page (Balcão/Login)
- `/login` -> Authentication handler (POST)
- `/games` -> The Shelf (Prateleira)
- `/games/:id` -> Game details and rental action
- `/profile` -> Member's card and personal Notebook