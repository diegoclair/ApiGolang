package controllers

import "github.com/diegoclair/ApiGolang/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareAuthentication(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	//s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Buys routes
	s.Router.HandleFunc("/buys", middlewares.SetMiddlewareJSON(s.CreateBuy)).Methods("POST")
	s.Router.HandleFunc("/buys", middlewares.SetMiddlewareAuthentication(s.GetBuys)).Methods("GET")
	s.Router.HandleFunc("/buys/{id}", middlewares.SetMiddlewareAuthentication(s.GetBuy)).Methods("GET")
	s.Router.HandleFunc("/buys/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateBuy))).Methods("PUT")
	//s.Router.HandleFunc("/buys/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteBuy)).Methods("DELETE")

	//Sales routes
	s.Router.HandleFunc("/sales", middlewares.SetMiddlewareJSON(s.CreateSale)).Methods("POST")
	s.Router.HandleFunc("/sales", middlewares.SetMiddlewareAuthentication(s.GetSales)).Methods("GET")
	s.Router.HandleFunc("/sales/{id}", middlewares.SetMiddlewareAuthentication(s.GetSale)).Methods("GET")
	s.Router.HandleFunc("/sales/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateSale))).Methods("PUT")
	//s.Router.HandleFunc("/sales/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteSale)).Methods("DELETE")

	//Reports
	s.Router.HandleFunc("/reports/id/{id}", middlewares.SetMiddlewareAuthentication(s.GetReportsByUserId)).Methods("GET")
	s.Router.HandleFunc("/reports/day/{date}", middlewares.SetMiddlewareAuthentication(s.GetReportsByDate)).Methods("GET")
}
