package routes

import "github.com/gorilla/mux"

func SetIncomesHandleActions(router *mux.Router) {

	route := "/incomes"
	router.HandleFunc(route, seeIncomes).Methods("GET")

	route = "/income/{id}"
	router.HandleFunc(route, seeOneIncome).Methods("GET")

	route = "/income"
	router.HandleFunc(route, createIncome).Methods("POST")

	route = "/income/{id}"
	router.HandleFunc(route, updateIncome).Methods("PATCH")

	route = "/income/{id}"
	router.HandleFunc(route, deleteIncome).Methods("DELETE")

}
