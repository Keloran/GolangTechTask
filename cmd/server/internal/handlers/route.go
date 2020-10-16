package handlers

import (
	"github.com/go-chi/chi"
)

func Routes(store Store) *chi.Mux {
	r := chi.NewRouter()

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
