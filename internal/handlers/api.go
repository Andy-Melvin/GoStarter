package handlers

import (
	"github.com/go-chi/chi"
	chimiddle"github.com/go-chi/chi/middleware"
	//Other pakage of Middle ware that must 
	//be in our file here
	
)

func Handler(r *chi.Mux){
	r.use(chimiddle.StipSlashes)
	r.Route("/account",func(router chi.router)){
		router.Use(middleware.Authorisation)
		router.Get("/coins",GetCoinBalance)

	}
} 