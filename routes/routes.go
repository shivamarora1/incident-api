package routes

import (
	"net/http"

	"example.com/incident-api/controllers"
	"example.com/incident-api/middlewares"
	"example.com/incident-api/utils"
	"github.com/gorilla/mux"
)

func SetUpRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(middlewares.LoggingMiddleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res := utils.ResponseFailure("Currently only GET/POST methods supported for /incident")
		w.Write(res)
	})

	r.HandleFunc("/incidents", controllers.AllIncidentGetter).Methods("GET")
	r.HandleFunc("/incidents", controllers.IncidentSaver).Methods("POST")
	return r
}
