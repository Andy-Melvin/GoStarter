package main

import(
	"fmt"
	"net/http"
	"github.com/go-chi/chi"
	//our own package here that is from the Internal, handlers
	log "github.com/sirupsen/logrus" 
)

func main() {
	log.SetReportCaller(true) 
	var r *chi.Mux=chi.NewRouter()
	handlers.Handlers(r)
	fmt.Println("Starting Go API service")
	
	err:=http.ListenAndServe("localhost:8080", r)

	if err != nil {
		log.Error(err)
	}

}