package crud

import (
	"errors"
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

func GetIncomes(root bool) ([]Income, error) {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	tableName := "incomes"
	query, _, err := GetQuery(tableName, Income{}, "SELECT", false)
	if err != nil {
		return nil, fmt.Errorf("can't get the query %s ERROR: %s", tableName, err.Error())
	}

	if !root {
		query += " WHERE deleted_at IS NULL"
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

	if IsDeleted(tableName, updatingIncome.Id) {
		return errors.New("regist is deleted")
	}

	query, data, err := GetQuery(tableName, *updatingIncome, "UPDATE", true)
	querySplit := strings.Split(query, "RETURNING") // Separate "UPDATE () SET () WHERE id = ()" + <stringToIntroduce> + "()"
	query = fmt.Sprintf("%s AND deleted_at IS NULL RETURNING %s", querySplit[0], querySplit[1])
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
