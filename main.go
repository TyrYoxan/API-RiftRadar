package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpAdapter "github.com/tyryoxan/API-RiftRadar/adapters/in/http"
	"github.com/tyryoxan/API-RiftRadar/adapters/out/composite"
	"github.com/tyryoxan/API-RiftRadar/adapters/out/memory"
	"github.com/tyryoxan/API-RiftRadar/adapters/out/riot"
	"github.com/tyryoxan/API-RiftRadar/app/usecases"
	"github.com/tyryoxan/API-RiftRadar/routes"
	"log"
	"net/http"
)

func main() {
	// Create a context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down...", sig)
		cancel()
	}()

	// Initialize repositories
	playerMemoryRepo := memory.NewPlayerRepo()
	riotRepo := riot.NewRiotPlayerRepo()
	playerRepo := composite.NewPlayerRepo(playerMemoryRepo, riotRepo)

	leaderboardRepo := memory.NewLeaderboardRepo()

	// Initialize services and actions
	playerService := &usecases.PlayerService{Repo: playerRepo}
	playerAction := &httpAdapter.PlayerAction{Service: playerService}

	leaderboardService := &usecases.LeaderboardService{Repo: leaderboardRepo, PlayerRepo: playerRepo}
	leaderboardAction := &httpAdapter.LeaderboardAction{Service: leaderboardService}

	// Setup router
	router := routes.SetupRouter(playerAction, leaderboardAction)

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for context cancellation (from signal handler)
	<-ctx.Done()

	// Create a deadline for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
