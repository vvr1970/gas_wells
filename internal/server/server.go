package server

import (
	"gas_wells/internal/handler"
	"gas_wells/internal/pkg/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
	logger logger.Logger
}

func New(logger logger.Logger) *Server {
	s := &Server{
		router: chi.NewRouter(),
		logger: logger,
	}

	s.setupMiddleware()
	return s
}

func (s *Server) SetupRoutes(wellHandler *handler.WellHandler) {
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/wells", http.StatusFound)
	})

	s.router.Route("/wells", func(r chi.Router) {
		r.Get("/", wellHandler.ListWells)
		r.Get("/create", wellHandler.CreateWellForm)
		r.Post("/", wellHandler.CreateWell)
		r.Get("/{id}", wellHandler.GetWell)
		r.Get("/{id}/edit", wellHandler.EditWellForm)
		r.Put("/{id}", wellHandler.UpdateWell)
		r.Delete("/{id}", wellHandler.DeleteWell)
	})

	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
	})
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(s.loggingMiddleware)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Compress(5))
	s.router.Use(middleware.Timeout(60 * time.Second))
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			s.logger.Info("Request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration", time.Since(start),
				"ip", r.RemoteAddr)
		}()

		next.ServeHTTP(ww, r)
	})
}

func (s *Server) ServeStatic(path string) {
	s.router.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir(path))))
}

func (s *Server) Handler() http.Handler {
	return s.router
}
