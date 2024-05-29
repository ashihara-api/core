package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
)

// Server provides an http.Server
type Server struct {
	l *slog.Logger
	*http.Server
}

// New creates and configures a server serving all application routes.
//
// The server implements a graceful shutdown and utilizes zap.Logger for logging purposes.
func New(listenAddr string, logger *slog.Logger, mux *chi.Mux) (*Server, error) {
	errorLog := slog.NewLogLogger(logger.Handler(), slog.LevelError)
	return NewWithHTTPServer(&http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ErrorLog:     errorLog,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}, logger)
}

// NewWithHTTPServer creates and configures a server serving all application routes.
func NewWithHTTPServer(srv *http.Server, logger *slog.Logger) (*Server, error) {
	return &Server{logger, srv}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown
func (srv *Server) Start() {
	srv.l.Info("Starting server...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srv.l.Error("Could not listen on", slog.String("addr", srv.Addr), slog.Any("error", err))
		}
	}()
	srv.l.Info("Server is ready to handle requests", slog.String("addr", srv.Addr))
	srv.gracefulShutdown()
}

func (srv *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	srv.l.Info("Server is shutting down", slog.String("reason", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		srv.l.Error("Could not gracefully shutdown the server", slog.Any("error", err))
	}
	srv.l.Info("Server stopped")
}
