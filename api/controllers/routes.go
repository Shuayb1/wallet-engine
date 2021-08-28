package controllers

import middlewares "opay/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Wallets routes
	s.Router.HandleFunc("/createwallet", middlewares.SetMiddlewareJSON(s.CreateWallet)).Methods("POST")
	s.Router.HandleFunc("/wallets", middlewares.SetMiddlewareJSON(s.GetWallets)).Methods("GET")
	s.Router.HandleFunc("/wallet/{id}", middlewares.SetMiddlewareJSON(s.GetWallet)).Methods("GET")
	s.Router.HandleFunc("/activate_a_wallet/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.ActivateAWallet))).Methods("PUT")
	s.Router.HandleFunc("/deactivate_a_wallet/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DeactivateAWallet))).Methods("PUT")
	s.Router.HandleFunc("/credit_a_wallet/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreditWallet))).Methods("POST")
	s.Router.HandleFunc("/debit_a_wallet/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DebitWallet))).Methods("POST")

	s.Router.HandleFunc("/deletewallet/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteWallet)).Methods("DELETE")

}
