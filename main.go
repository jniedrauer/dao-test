package main

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable",
		"localhost",
		5432,
		"testdb",
		"user",
		"password",
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}

	/*row := &Row{
		ID:   2,
		Data: "bar",
	}

	if err := Insert(row, db); err != nil {
		panic(err)
	}
	*/

	v, err := Read(2, db)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", v)
}

// Row is a single row in the database.
type Row struct {
	ID   int    `sql:"id"`
	Data string `sql:"data"`
}

// Columns returns a comma delimited list of all columns.
func Columns(row interface{}) string {
	t := reflect.TypeOf(row).Elem()
	columns := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		columns[i] = t.Field(i).Tag.Get("sql")
	}

	return strings.Join(columns, ",")
}

// Fields returns pointers to all fields in a struct.
func Fields(row interface{}) []interface{} {
	v := reflect.ValueOf(row)

	if v.Kind() != reflect.Ptr {
		panic("not a pointer")
	}

	fields := make([]interface{}, v.Elem().NumField())
	for i := range fields {
		fields[i] = v.Elem().Field(i).Addr().Interface()
	}

	return fields
}

// ErrNoResults returns nil database results.
var ErrNoResults = errors.New("no results")

// Read reads a row from the database by ID.
func Read(id int, db *sql.DB) (*Row, error) {
	row := &Row{}

	rows, err := db.Query("SELECT "+Columns(row)+" FROM test_table WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(Fields(row)...); err != nil {
			return nil, err
		}

		return row, nil
	}

	return nil, ErrNoResults
}

// fmtString generates parameterized fields for psql.
func fmtString(count int) string {
	result := "$1"
	for i := 0; i < count-1; i++ {
		result += ",$"
		result += strconv.Itoa(i + 2)
	}

	return result
}

// Insert writes a row to the database.
func Insert(row *Row, db *sql.DB) error {
	fields := Fields(row)

	query := fmt.Sprintf("INSERT INTO test_table (%s) VALUES (%s)",
		Columns(row),
		fmtString(len(fields)),
	)

	_, err := db.Exec(query, Fields(row)...)
	return err
}
