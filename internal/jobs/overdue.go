package jobs

import (
	"context"
	"log"
	"time"

	"github.com/cmellojr/modo-locadora/internal/database"
)

// StartOverdueChecker launches a goroutine that periodically checks for overdue
// rentals and auto-returns them, penalizing the member. Stops on ctx cancellation.
func StartOverdueChecker(ctx context.Context, store database.Store, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		log.Printf("[overdue-checker] Started. Checking every %v", interval)

		// Run once immediately on startup.
		processOverdue(ctx, store)

		for {
			select {
			case <-ctx.Done():
				log.Println("[overdue-checker] Shutting down gracefully.")
				return
			case <-ticker.C:
				processOverdue(ctx, store)
			}
		}
	}()
}

func processOverdue(ctx context.Context, store database.Store) {
	count, err := store.ProcessOverdueRentals(ctx)
	if err != nil {
		log.Printf("[overdue-checker] Error processing overdue rentals: %v", err)
		return
	}
	if count > 0 {
		log.Printf("[overdue-checker] Auto-returned %d overdue rental(s).", count)
	}
}
