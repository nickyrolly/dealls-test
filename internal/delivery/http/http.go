package http

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/healthcheck"
	middlewares "github.com/nickyrolly/dealls-test/internal/delivery/http/middleware"
)

type RouteConfig struct {
	Router *chi.Mux
}

func (c *RouteConfig) Setup() {
	c.SetupAPI()
}

func (c *RouteConfig) SetupAPI() {
	fmt.Println("========== Dealls Setup API ==========")

	c.Router.Use(middlewares.CorsMiddleware)
	c.Router.Get("/", healthcheck.HandleHealthCheck)
	c.Router.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middlewares.BasicMiddleware)
			// r.Post("/signup", handler.ParticipantHandler.HandlerUpsertParticipant)
		})

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthenticationMiddleware)
			// r.Post("/login", handler.ParticipantHandler.HandleGetParticipantLogin)
		})

		r.Route("/v1", func(r chi.Router) {
			// r.Use(middlewares.SessionMiddleware)
			// r.Method(http.MethodOptions, "/*", http.HandlerFunc(preflight.HandlePreflight))
			// r.Get("/participant", handler.ParticipantHandler.HandleGetParticipantData)
		})
	})
}
