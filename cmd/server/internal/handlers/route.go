package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
  "github.com/keloran/go-probe"
)

func Routes(store Store) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/probe", probe.HTTP)

	r.Route("/buff", func(r chi.Router) {
    bh := buffHandler{
      store: store,
    }

		r.Get("/{id}", bh.GetBuff)
		r.Post("/", bh.CreateBuff)
		r.Get("/", bh.ListBuffs)
	})

	r.Route("/stream", func(r chi.Router) {
	  sh := streamHandler{
	    store: store,
    }

	  r.Get("/", sh.ListStreams)
	  r.Post("/", sh.CreateStream)
	  r.Put("/{id}", sh.AttachBuffs)
	  r.Get("/{id}", sh.GetStream)
  })

	return r
}
