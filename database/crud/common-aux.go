package crud

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LuisFlahan4051/basic-counter-app/database"
)

func valueIsNull(fieldType reflect.Kind, value reflect.Value) bool {
	switch fieldType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.String:
		return value.String() == "" || strings.Contains(value.String(), "0001-01-01")
	case reflect.Slice:
		return value.Len() == 0
	case reflect.Struct:
		return value.IsZero()
	case reflect.Interface:
		return value.IsZero()
	case reflect.Ptr:
		return value.IsNil()
	}
	return false
}

func GetValueField(fieldType reflect.Kind, value reflect.Value) interface{} {
	switch fieldType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint()
	case reflect.Float32, reflect.Float64:
		return value.Float()
	case reflect.Bool:
		return value.Bool()
	case reflect.String:
		return value.String()
	case reflect.Slice:
		return value.Interface()
	case reflect.Struct:
		return value.Interface()
	case reflect.Interface:
		return value.Interface()
	case reflect.ValueOf(time.Time{}).Kind():
		return value.Interface()
	case reflect.Ptr:
		return value.Elem()
	}
	return nil
}

// Can add `flahan:"ignore"` into structs to ignore a field
func GetStructFields(anyStruct interface{}, slices bool, nulls bool) ([]string, map[string]interface{}) {
	var fields []string
	valueFields := make(map[string]interface{})

	_struct := reflect.ValueOf(anyStruct)

	for i := 0; i < _struct.NumField(); i++ {
		_field := _struct.Type().Field(i)
		value := _struct.Field(i)

		fieldType := _field.Type.Kind()
		flahanTag := _struct.Type().Field(i).Tag.Get("flahan")

		if flahanTag == "ignore" {
			continue
		}
		switch fieldType {
		case reflect.Struct, reflect.Interface:
			subStruct := _struct.Field(i).Interface()
			subFields, subValueFields := GetStructFields(subStruct, slices, nulls)
			fields = append(fields, subFields...)
			for subKey, subValue := range subValueFields {
				valueFields[subKey] = subValue
			}
		}

		valueFieldIsNull := valueIsNull(fieldType, value)

		jsonTag := _struct.Type().Field(i).Tag.Get("json")
		if jsonTag != "" {
			if (!slices && fieldType == reflect.Slice) || (!nulls && valueFieldIsNull) {
				continue
			}

			jsonTagOptions := strings.Split(jsonTag, ",")
			fieldName := jsonTagOptions[0]
			fields = append(fields, fieldName)

			valueField := GetValueField(fieldType, value)
			if valueField != nil {
				valueFields[fieldName] = valueField
			}
		}
	}

	return fields, valueFields
}

func DecodeRowsToJson(rows *sql.Rows) ([]string, error) {
	columns, _ := rows.Columns()
	countColumns := len(columns)
	values := make([]interface{}, countColumns)
	valuePointers := make([]interface{}, countColumns)

	var jsonEncodes []string
	for rows.Next() {
		for i := range columns {
			valuePointers[i] = &values[i]
		}

		rows.Scan(valuePointers...)
		mapToDecode := make(map[string]interface{})

		for i, column := range columns {
			value := values[i]

			bytes, ok := value.([]byte)
			var finalValue interface{}
			if ok {
				finalValue = string(bytes)
			} else {
				finalValue = value
			}

			mapToDecode[column] = finalValue
		}

		jsonEncode, err := json.Marshal(mapToDecode)
		if err != nil {
			return nil, fmt.Errorf("error while encoding map to json ERROR: %v", err)
		}

		jsonEncodes = append(jsonEncodes, string(jsonEncode))
	}

	return jsonEncodes, nil
}

// ------------------------ Query Builder ------------------------ //

// RETURNS: "$1, $2, $3..."
func getIndexFormated(slice []string) string {
	var index string
	for i := 0; i < len(slice)-1; i++ {
		index += "$" + strconv.Itoa(i+1) + ", "
	}
	index += "$" + strconv.Itoa(len(slice))
	return index
}

// RETURNS: "field1 = $1, field2 = $2, field3 = $3..."
func getIndexAndFieldsFormated(slice []string) string {
	var index string
	for i := 0; i < len(slice)-1; i++ {
		index += fmt.Sprintf("%s = $%d, ", slice[i], i+1)
	}
	i := len(slice) - 1
	index += fmt.Sprintf("%s = $%d", slice[i], i+1)
	return index
}

// Returns an interface needed to be inserted in the database. User before GetStructFields function.
func getDataFields(fieldsSlice []string, fieldsValuesMap map[string]interface{}) []interface{} {
	var data []interface{}
	for _, field := range fieldsSlice {
		valueString := fmt.Sprintf("%v", fieldsValuesMap[field])

		regularExpresion := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
		if regularExpresion.MatchString(valueString) {
			date, _ := time.Parse("2006-01-02 15:04:05", strings.Split(valueString, " -")[0])

			data = append(data, date)
		} else {
			data = append(data, valueString)
		}
	}
	return data
}

// OPERATION: "INSERT", "UPDATE", "SELECT"; | RETURNS query string, data interface and error
func GetQuery(table string, strc interface{}, operation string, returning bool) (string, []interface{}, error) {
	var query string
	var data []interface{}

	switch operation {
	case "INSERT":
		fieldsSlice, fieldsValuesMap := GetStructFields(strc, false, false)
		fields := strings.Join(fieldsSlice, ", ")

		index := getIndexFormated(fieldsSlice)

		query = fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", table, fields, index)

		if returning {
			allFieldsSlice, _ := GetStructFields(strc, false, true)
			allFields := strings.Join(allFieldsSlice, ", ")
			query += " RETURNING " + allFields
		}

		data = getDataFields(fieldsSlice, fieldsValuesMap)
	case "UPDATE":
		fieldsSlice, fieldsValuesMap := GetStructFields(strc, false, false)

		FieldsWithIndex := getIndexAndFieldsFormated(fieldsSlice)
		query = fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", table, FieldsWithIndex, fieldsValuesMap["id"])

		data = getDataFields(fieldsSlice, fieldsValuesMap)

		if returning {
			allFieldsSlice, _ := GetStructFields(strc, false, true)
			allFields := strings.Join(allFieldsSlice, ", ")
			query += " RETURNING " + allFields
		}

	case "SELECT":
		fieldsSlice, _ := GetStructFields(strc, false, true)
		fields := strings.Join(fieldsSlice, ", ")

		query = fmt.Sprintf("SELECT %s FROM %s", fields, table)

	default:
		return "", nil, errors.New("operation not valid")
	}

	return query, data, nil
}

// tables := make(map[string]interface{}); tables["tableName1"] = models.Table1; tables["tableName2"] = models.Table2; ... RETURNS a query string
func GetMixSelect(tables map[string]interface{}) string {

	var mixFields []string
	var mixTables []string

	for table, strct := range tables {
		log.Println(table)
		log.Println(strct)

		nameSingleTable := table + "_single"
		tableAndSingleTable := fmt.Sprintf("%s %s", table, nameSingleTable)

		mixTables = append(mixTables, tableAndSingleTable)

		fields, _ := GetStructFields(strct, false, true)
		var tableAndFields []string

		for _, field := range fields {
			tableAndFields = append(tableAndFields, fmt.Sprintf("%s.%s", nameSingleTable, field))
		}
		mixFields = append(mixFields, strings.Join(tableAndFields, ", "))
	}

	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(mixFields, ", "), strings.Join(mixTables, ", "))

	return query
}

func IsDeleted(table string, id uint) bool {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	var deletedAt *time.Time
	query := fmt.Sprintf("SELECT deleted_at FROM %s WHERE id = $1", table)
	err := db.QueryRow(query, id).Scan(&deletedAt)

	if err != nil {
		return false
	}

	return deletedAt != nil
}

func DeleteFromTableById(tableName string, id uint) error {
	db := database.GetConnection(database.DATABASE_NAME)
	defer db.Close()

	if IsDeleted(tableName, id) {
		return errors.New("already deleted")
	}

	query := fmt.Sprintf("UPDATE %s SET deleted_at = NOW() WHERE id = $1", tableName)
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("can't execute the %s query ERROR: %s", tableName, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete one %s ERROR: %s", tableName, err.Error())
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s not found", tableName)
	}

	return nil
}
