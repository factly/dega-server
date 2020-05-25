package tag

import "github.com/go-chi/chi"

// tag model
type tag struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description"`
	SpaceID     uint   `json:"space_id" validate:"required"`
}

// Router - Group of tag router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)
	r.Post("/", create)

	r.Route("/{tag_id}", func(r chi.Router) {
		r.Get("/", details)
		r.Put("/", update)
		r.Delete("/", delete)
	})

	return r

}
