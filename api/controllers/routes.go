package controllers

import "github.com/victorsteven/fullstack/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")

	//Drugs routes
	s.Router.HandleFunc("/drugs", middlewares.SetMiddlewareJSON(s.CreateDrug)).Methods("POST")
	s.Router.HandleFunc("/drugs", middlewares.SetMiddlewareJSON(s.GetDrug)).Methods("GET")
	s.Router.HandleFunc("/drugs/{id}", middlewares.SetMiddlewareJSON(s.GetDrug)).Methods("GET")
	s.Router.HandleFunc("/drugs/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateDrug))).Methods("PUT")
	s.Router.HandleFunc("/drugs/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteDrug)).Methods("DELETE")
}
