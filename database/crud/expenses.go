package crud

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/LuisFlahan4051/basic-counter-app/database"
)

type Expense struct {
	Id        uint       `json:"id"`
	Value     float64    `json:"value"`
	Type      string     `json:"type"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func NewExpense(expense *Expense) error {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	tableName := "expenses"
	query, data, err := GetQuery(tableName, *expense, "INSERT", true)
	if err != nil {
		return fmt.Errorf("can't get the query %s ERROR: %s", tableName, err.Error())
	}

	err = db.QueryRow(query, data...).Scan(
		&expense.Id,
		&expense.Value,
		&expense.Type,
		&expense.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("can't execute the %s query ERROR: %s", tableName, err.Error())
	}

	if expense.Id == 0 {
		return fmt.Errorf("can't create %s", tableName)
	}

	return nil
}

func GetExpenses(root bool) ([]Expense, error) {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	tableName := "expenses"
	query, _, err := GetQuery(tableName, Expense{}, "SELECT", false)
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

	var expenses []Expense
	for rows.Next() {
		var expense Expense
		err = rows.Scan(
			&expense.Id,
			&expense.Value,
			&expense.Type,
			&expense.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("can't execute the %s query ERROR: %s", tableName, err.Error())
		}

		expenses = append(expenses, expense)
	}

	if len(expenses) == 0 {
		return nil, fmt.Errorf("%s not found", tableName)
	}

	return expenses, nil
}

func GetExpense(id uint) (Expense, error) {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	tableName := "expenses"
	expense := Expense{}
	query, _, err := GetQuery(tableName, expense, "SELECT", false)
	if err != nil {
		return expense, fmt.Errorf("can't get the query %s ERROR: %s", tableName, err.Error())
	}
	query += " WHERE id = $1"

	err = db.QueryRow(query, id).Scan(
		&expense.Id,
		&expense.Value,
		&expense.Type,
		&expense.CreatedAt,
	)
	if err != nil {
		return expense, fmt.Errorf("can't get the query %s ERROR: %s", tableName, err.Error())
	}

	if expense.Id == 0 {
		return expense, fmt.Errorf("%s not found", tableName)
	}

	return expense, nil
}

func DeleteExpense(id uint) error {
	return DeleteFromTableById("expenses", id)
}

func UpdateExpense(updatingExpense *Expense) error {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	tableName := "expenses"

	if IsDeleted(tableName, updatingExpense.Id) {
		return errors.New("regist is deleted")
	}

	query, data, err := GetQuery(tableName, *updatingExpense, "UPDATE", true)
	querySplit := strings.Split(query, "RETURNING") // Separate "UPDATE () SET () WHERE id = ()" + <stringToIntroduce> + "()"
	query = fmt.Sprintf("%s AND deleted_at IS NULL RETURNING %s", querySplit[0], querySplit[1])
	if err != nil {
		return fmt.Errorf("can't get the query %s ERROR: %s", tableName, err.Error())
	}

	err = db.QueryRow(query, data...).Scan(
		&updatingExpense.Id,
		&updatingExpense.Value,
		&updatingExpense.Type,
		&updatingExpense.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("can't execute the %s query ERROR: %s", tableName, err.Error())
	}

	if updatingExpense.Id == 0 {
		return fmt.Errorf("%s not found", tableName)
	}

	return nil
}
