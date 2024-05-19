package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/omegaatt36/instagramrobot/logging"
)

// Server is the main controller for the bot.
type Server struct {
	port int
}

func NewServer() *Server {
	return &Server{
		port: 8080,
	}
}

func (s *Server) startHttp(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/parse/", s.addFilm)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		logging.Fatalf("listen: %s\n", err)
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logging.Fatalf("shutdown: %s\n", err)
		}
	}()
}

// Start brings bot into motion by consuming incoming updates
func (s *Server) Start(ctx context.Context) <-chan struct{} {
	logging.Info("Instagram fetcher web starting")
	closeChain := make(chan struct{})

	go s.startHttp(ctx)

	go func() {
		defer func() {
			logging.Info("Instagram fetcher web stopped")
			closeChain <- struct{}{}
			close(closeChain)
		}()

		<-ctx.Done()
		// b.Stop()
	}()

	return closeChain
}
