package controllers

import "github.com/diegoclair/ApiGolang/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")

	//Buys routes
	s.Router.HandleFunc("/buys", middlewares.SetMiddlewareJSON(s.CreateBuy)).Methods("POST")
	s.Router.HandleFunc("/buys", middlewares.SetMiddlewareJSON(s.GetBuys)).Methods("GET")
	s.Router.HandleFunc("/buys/{id}", middlewares.SetMiddlewareJSON(s.GetBuy)).Methods("GET")
	s.Router.HandleFunc("/buys/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateBuy))).Methods("PUT")

	//Sales routes
	s.Router.HandleFunc("/sales", middlewares.SetMiddlewareJSON(s.CreateSale)).Methods("POST")
	s.Router.HandleFunc("/sales", middlewares.SetMiddlewareJSON(s.GetSales)).Methods("GET")
	s.Router.HandleFunc("/sales/{id}", middlewares.SetMiddlewareJSON(s.GetSale)).Methods("GET")
	s.Router.HandleFunc("/sales/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateSale))).Methods("PUT")
}