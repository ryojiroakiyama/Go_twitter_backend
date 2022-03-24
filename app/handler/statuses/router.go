package statuses

import (
	"net/http"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/statuses/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(app))
		r.Post("/", h.Create)
		r.Delete("/{id}", h.Delete)
	})
	r.Get("/{id}", h.Fetch)

	return r
}
