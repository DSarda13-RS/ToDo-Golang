package server

import (
	"ToDo/handler"
	"github.com/go-chi/chi/v5"
)

func taskRoutes(r chi.Router) {
	r.Group(func(task chi.Router) {
		task.Post("/create", handler.CreateTask)
		task.Post("/update", handler.UpdateTask)
		task.Get("/info", handler.GetTask)
		task.Delete("/delete", handler.DeleteTask)
	})
}
