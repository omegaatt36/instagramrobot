package web

import (
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/omegaatt36/instagramrobot/appmodule/instagram"
	"github.com/omegaatt36/instagramrobot/appmodule/providers"
	"github.com/omegaatt36/instagramrobot/appmodule/threads"
	"github.com/omegaatt36/instagramrobot/logging"
)

// Server manages the web application, including routing and handling HTTP requests.
type Server struct {
	// port is the TCP port the server listens on.
	port int
	// indexPage is the parsed HTML template for the main page.
	indexPage *template.Template
	// linkProcessor is the instance of LinkProcessor used to process links.
	linkProcessor *providers.LinkProcessor
}

// NewServer creates and initializes a new Server instance.
// It parses the embedded HTML template and initializes the Instagram fetcher.
// It panics if template parsing fails.
func NewServer() *Server {
	index, err := template.ParseFS(indexHTML, "index.html")
	if err != nil {
		panic(err)
	}
	return &Server{
		port:      8080,
		indexPage: index,
		linkProcessor: providers.NewLinkProcessor(providers.NewLinkProcessorRequest{
			InstagramFetcher: instagram.NewExtractor(),
			ThreadsFetcher:   threads.NewExtractor(),
		}),
	}
}

// startHTTP configures the HTTP routes, creates an HTTP server, and starts it.
// It also sets up graceful shutdown based on the provided context.
func (s *Server) startHTTP(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/parse/", s.parse)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logging.Fatalf("shutdown: %s\n", err)
		}
	}()
}

// Start begins the web server's listening process in a separate goroutine.
// It returns a channel that closes when the server has gracefully shut down.
// It uses the provided context to trigger the shutdown process.
func (s *Server) Start(ctx context.Context) <-chan struct{} {
	logging.Info("Instagram fetcher web starting")
	closeChain := make(chan struct{})

	go s.startHTTP(ctx)

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
