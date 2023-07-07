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

func seeOneIncome(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		commons.Logcatch(writer, http.StatusBadRequest, errors.New("invalid income/{id} value"))
		return
	}

	income, err := crud.GetIncome(uint(id))
	if err != nil {
		commons.Logcatch(writer, http.StatusNotFound, err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(income)
}

func seeIncomes(writer http.ResponseWriter, request *http.Request) {

	sinceIncome := request.URL.Query().Get("since")
	toIncome := request.URL.Query().Get("to")
	typeIncome := request.URL.Query().Get("type")

	var since, to *time.Time

	regularExpresion := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	if strings.Compare(sinceIncome, "") != 0 {
		if regularExpresion.MatchString(sinceIncome) {
			//Splitting the string to get just the date
			date, err := time.Parse("2006-01-02", strings.Split(sinceIncome, " -")[0])
			commons.Logcatch(writer, http.StatusBadRequest, err)
			since = &date
		} else {
			commons.Logcatch(writer, http.StatusBadRequest, errors.New("invalid formate of date in since parameter"))
			return
		}
	}
	if strings.Compare(toIncome, "") != 0 {
		if regularExpresion.MatchString(toIncome) {
			date, _ := time.Parse("2006-01-02", strings.Split(toIncome, " -")[0])
			to = &date
		} else {
			commons.Logcatch(writer, http.StatusBadRequest, errors.New("invalid formate of date in to parameter"))
			return
		}
	}

	incomes, err := crud.GetIncomes(since, to, &typeIncome)

	if err != nil {
		commons.Logcatch(writer, http.StatusNotFound, err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(incomes)
}

func createIncome(writer http.ResponseWriter, request *http.Request) {
	var income crud.Income
	err := json.NewDecoder(request.Body).Decode(&income)
	if err != nil {
		commons.Logcatch(writer, http.StatusBadRequest, err)
		return
	}

	if income.CreatedAt == nil {
		today := time.Now()
		income.CreatedAt = &today
	}

	err = crud.NewIncome(&income)
	if err != nil {
		commons.Logcatch(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(income)
}

func updateIncome(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		commons.Logcatch(writer, http.StatusBadRequest, errors.New("invalid income/{id} value"))
		return
	}

	var income crud.Income
	err = json.NewDecoder(request.Body).Decode(&income)
	if err != nil {
		commons.Logcatch(writer, http.StatusBadRequest, err)
		return
	}
	income.Id = uint(id)

	err = crud.UpdateIncome(&income)
	if err != nil {
		commons.Logcatch(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(income)
}

func deleteIncome(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		commons.Logcatch(writer, http.StatusBadRequest, errors.New("invalid income/{id} value"))
		return
	}

	err = crud.DeleteIncome(uint(id))
	if err != nil {
		commons.Logcatch(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode("done")
}
