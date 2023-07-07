package routes

import "github.com/gorilla/mux"

func SetExpensesHandleActions(router *mux.Router) {

	route := "/expenses"
	router.HandleFunc(route, seeExpenses).Methods("GET")

	route = "/expense/{id}"
	router.HandleFunc(route, seeOneExpense).Methods("GET")

	// created_at need to be in this format: 2023-02-04T00:00:00.00Z
	route = "/expense"
	router.HandleFunc(route, createExpense).Methods("POST")

	route = "/expense/{id}"
	router.HandleFunc(route, updateExpense).Methods("PATCH")

	route = "/expense/{id}"
	router.HandleFunc(route, deleteExpense).Methods("DELETE")

}
