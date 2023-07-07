package crud

import (
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

func GetExpenses(sinceFilter *time.Time, toFilter *time.Time, typeFilter *string) ([]Expense, error) {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	tableName := "expenses"
	query, _, err := GetQuery(tableName, Expense{}, "SELECT", false)
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

	query, data, err := GetQuery(tableName, *updatingExpense, "UPDATE", true)
	querySplit := strings.Split(query, "RETURNING") // Separate "UPDATE () SET () WHERE id = ()" + <stringToIntroduce> + "()"
	query = fmt.Sprintf("%s RETURNING %s", querySplit[0], querySplit[1])
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
