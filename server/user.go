package server

import (
	"ToDo/handler"
	"github.com/go-chi/chi/v5"
)

func userRoutes(r chi.Router) {
	r.Group(func(user chi.Router) {
		user.Get("/info", handler.GetUserInfo)
		user.Delete("/logout", handler.Logout)
	})
}
