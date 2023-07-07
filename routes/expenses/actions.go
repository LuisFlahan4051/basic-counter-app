package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LuisFlahan4051/basic-counter-app/database/crud"
	commons "github.com/LuisFlahan4051/basic-counter-app/routes/common-aux"
	"github.com/gorilla/mux"
)

func seeOneExpense(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		commons.Logcatch(writer, http.StatusBadRequest, errors.New("invalid expense/{id} value"))
		return
	}

	expense, err := crud.GetExpense(uint(id))
	if err != nil {
		commons.Logcatch(writer, http.StatusNotFound, err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(expense)
}

func seeExpenses(writer http.ResponseWriter, request *http.Request) {

	sinceExpense := request.URL.Query().Get("since")
	toExpense := request.URL.Query().Get("to")
	typeExpense := request.URL.Query().Get("type")

	var since, to *time.Time

	regularExpresion := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	if strings.Compare(sinceExpense, "") != 0 {
		if regularExpresion.MatchString(sinceExpense) {
			//Splitting the string to get just the date
			date, err := time.Parse("2006-01-02", strings.Split(sinceExpense, " -")[0])
			commons.Logcatch(writer, http.StatusBadRequest, err)
			since = &date
		} else {
			commons.Logcatch(writer, http.StatusBadRequest, errors.New("invalid formate of date in since parameter"))
			return
		}
	}
	if strings.Compare(toExpense, "") != 0 {
		if regularExpresion.MatchString(toExpense) {
			date, _ := time.Parse("2006-01-02", strings.Split(toExpense, " -")[0])
			to = &date
		} else {
			commons.Logcatch(writer, http.StatusBadRequest, errors.New("invalid formate of date in to parameter"))
			return
		}
	}

	expenses, err := crud.GetExpenses(since, to, &typeExpense)

	if err != nil {
		commons.Logcatch(writer, http.StatusNotFound, err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(expenses)
}

func createExpense(writer http.ResponseWriter, request *http.Request) {
	var expense crud.Expense
	err := json.NewDecoder(request.Body).Decode(&expense)
	if err != nil {
		commons.Logcatch(writer, http.StatusBadRequest, err)
		return
	}

	err = crud.NewExpense(&expense)
	if err != nil {
		commons.Logcatch(writer, http.StatusInternalServerError, err)
		return
	}

	if expense.CreatedAt == nil {
		today := time.Now()
		expense.CreatedAt = &today
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(expense)
}

func updateExpense(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		commons.Logcatch(writer, http.StatusBadRequest, errors.New("invalid expense/{id} value"))
		return
	}

	var expense crud.Expense
	err = json.NewDecoder(request.Body).Decode(&expense)
	if err != nil {
		commons.Logcatch(writer, http.StatusBadRequest, err)
		return
	}
	expense.Id = uint(id)

	err = crud.UpdateExpense(&expense)
	if err != nil {
		commons.Logcatch(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(expense)
}

func deleteExpense(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		commons.Logcatch(writer, http.StatusBadRequest, errors.New("invalid expense/{id} value"))
		return
	}

	err = crud.DeleteExpense(uint(id))
	if err != nil {
		commons.Logcatch(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode("done")
}
