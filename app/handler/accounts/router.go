package accounts

import (
	"net/http"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
)

// Implementation of handler
// http handlerで引数以外の情報を共有するには,
// struct・method形式にしてレシーバーを利用するとscopeはっきりしてわかりやすいし簡単
type handler struct {
	app *app.App
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(app))
		r.Post("/{username}/follow", h.Follow)
		r.Post("/{username}/unfollow", h.UnFollow)
		r.Post("/update_credentials", h.Update)
	})
	r.Post("/", h.Create)
	r.Get("/{username}", h.Fetch)
	r.Get("/{username}/following", h.Following)
	r.Get("/{username}/followers", h.Followers)

	return r
}
