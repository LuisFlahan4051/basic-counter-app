package crud

import (
	"fmt"
	"strings"
	"time"

	"github.com/LuisFlahan4051/basic-counter-app/database"
)

type Income struct {
	Id        uint       `json:"id"`
	Value     float64    `json:"value"`
	Type      string     `json:"type"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func NewIncome(income *Income) error {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	tableName := "incomes"
	query, data, err := GetQuery(tableName, *income, "INSERT", true)
	if err != nil {
		return fmt.Errorf("can't get the query %s ERROR: %s", tableName, err.Error())
	}

	err = db.QueryRow(query, data...).Scan(
		&income.Id,
		&income.Value,
		&income.Type,
		&income.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("can't execute the %s query ERROR: %s", tableName, err.Error())
	}

	if income.Id == 0 {
		return fmt.Errorf("can't create %s", tableName)
	}

	return nil
}

func GetIncomes(sinceFilter *time.Time, toFilter *time.Time, typeFilter *string) ([]Income, error) {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	tableName := "incomes"
	query, _, err := GetQuery(tableName, Income{}, "SELECT", false)
	if err != nil {
		return nil, fmt.Errorf("can't get the query %s ERROR: %s", tableName, err.Error())
	}

	whereIsSet := false
	if strings.Compare(*typeFilter, "") != 0 {
		query = fmt.Sprintf("%s WHERE type = '%s'", query, *typeFilter)
		whereIsSet = true
	}

	if sinceFilter != nil {
		if toFilter == nil {
			toFilter = &time.Time{}
		}

		//Verify intervals and order
		if sinceFilter.After(*toFilter) && !toFilter.IsZero() {
			sinceFilter, toFilter = toFilter, sinceFilter
		}

		//Set To to today hour 23:59:59 when Since exists and is unknown
		if toFilter.IsZero() && !sinceFilter.IsZero() {
			today, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
			todayPtr := today.Add(24 * time.Hour)
			toFilter = &todayPtr
		}
		if whereIsSet {
			query = fmt.Sprintf("%s AND created_at BETWEEN '%s' AND '%s' ORDER BY created_at DESC", query, sinceFilter, toFilter)
		} else {
			query = fmt.Sprintf("%s WHERE created_at BETWEEN '%s' AND '%s' ORDER BY created_at DESC", query, sinceFilter, toFilter)
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("can't execute the %s query ERROR: %s", tableName, err.Error())
	}
	defer rows.Close()

	var incomes []Income
	for rows.Next() {
		var income Income
		err = rows.Scan(
			&income.Id,
			&income.Value,
			&income.Type,
			&income.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("can't execute the %s query ERROR: %s", tableName, err.Error())
		}

		incomes = append(incomes, income)
	}

	if len(incomes) == 0 {
		return nil, fmt.Errorf("%s not found", tableName)
	}

	return incomes, nil
}

func GetIncome(id uint) (Income, error) {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	tableName := "incomes"
	income := Income{}
	query, _, err := GetQuery(tableName, income, "SELECT", false)
	if err != nil {
		return income, fmt.Errorf("can't get the query %s ERROR: %s", tableName, err.Error())
	}
	query += " WHERE id = $1"

	err = db.QueryRow(query, id).Scan(
		&income.Id,
		&income.Value,
		&income.Type,
		&income.CreatedAt,
	)
	if err != nil {
		return income, fmt.Errorf("can't get the query %s ERROR: %s", tableName, err.Error())
	}

	if income.Id == 0 {
		return income, fmt.Errorf("%s not found", tableName)
	}

	return income, nil
}

func DeleteIncome(id uint) error {
	return DeleteFromTableById("incomes", id)
}

func UpdateIncome(updatingIncome *Income) error {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	tableName := "incomes"

	query, data, err := GetQuery(tableName, *updatingIncome, "UPDATE", true)
	querySplit := strings.Split(query, "RETURNING") // Separate "UPDATE () SET () WHERE id = ()" + <stringToIntroduce> + "()"
	query = fmt.Sprintf("%s RETURNING %s", querySplit[0], querySplit[1])
	if err != nil {
		return fmt.Errorf("can't get the query %s ERROR: %s", tableName, err.Error())
	}

	err = db.QueryRow(query, data...).Scan(
		&updatingIncome.Id,
		&updatingIncome.Value,
		&updatingIncome.Type,
		&updatingIncome.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("can't execute the %s query ERROR: %s", tableName, err.Error())
	}

	if updatingIncome.Id == 0 {
		return fmt.Errorf("%s not found", tableName)
	}

	return nil
}
