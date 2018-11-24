package helper

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

var dbi = dbInfo{
	Host:     "itradecoin-db-instance.cppwjzff6pjg.us-east-2.rds.amazonaws.com",
	Port:     5433,
	User:     "postgres",
	Password: "ikechukwu1",
	Dbname:   "webdb",
}

/*
var dbi = dbInfo{
	Host:     "localhost",
	Port:     5433,
	User:     "postgres",
	Password: "1234",
	Dbname:   "webdb",
}
*/
// DbCon this hold pointer to db
type DbCon struct {
	Db *sql.DB
}

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	dbi.Host, dbi.Port, dbi.User, dbi.Password, dbi.Dbname)

// OpenConnection is used open the connection to main DB
func OpenConnection() (*DbCon, error) {

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		//panic(err)
		//return nil, err
		fmt.Println("Database connection failed due to", err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Database connection failed due to", err)
		//return nil, err
	}
	//fmt.Println("Successfully connected To DB!")

	return &DbCon{
		Db: db,
	}, nil
}

// Close  our db conncetion.
func (con *DbCon) Close() error {
	return con.Db.Close()
}

//DbSelectRow returns a value selected from the database or error
func DBSelectRow(query string, num ...interface{}) *RowSelect {
	columns := make(map[string]interface{})
	var querySlice []string
	var cols []string
	if c := strings.Contains(query, "FROM"); c {
		querySlice = strings.SplitN(query, "FROM", 2)
	}
	//if c := strings.Contains(query, "from"); c {
	//querySlice = strings.SplitN(query, "from", 2)
	//}
	w := querySlice[0]  //this has the first part of the query before FROM
	sq := querySlice[1] //this has the second part after FROM
	//fmt.Println("The Query Is:", w, sq)
	if con := strings.Contains(w, "SELECT"); con {
		cols = strings.SplitAfter(w, "SELECT")
	}
	if con := strings.Contains(w, "select"); con {
		cols = strings.SplitAfter(w, "select")
	}
	col := strings.Replace(cols[1], " ", "", -1)
	colSlice := strings.Split(col, ",")
	con, err := OpenConnection()
	if err != nil {
		return &RowSelect{nil, err.Error()}
	}
	defer con.Close()
	var ar interface{}
	for _, d := range colSlice {
		err = con.Db.QueryRow(fmt.Sprintf("SELECT %v FROM %s", d, sq), num...).Scan(&ar)
		if err != nil {
			return &RowSelect{nil, err.Error()}
		}
		columns[d] = ar
	}
	//fmt.Println("EMAIL: ", columns)
	return &RowSelect{columns, ""}
}

//DBModify Returns the number of rows affected after executing a db query
func DBModify(query string, num ...interface{}) *ModifyDb {
	con, err := OpenConnection()
	if err != nil {
		//fmt.Println("connection failed:", err)
		return &ModifyDb{-1, err.Error()}
	}
	defer con.Close()
	queryFields := strings.Fields(query)
	var queryString string
	if strings.ToUpper(queryFields[0]) == "INSERT" {
		var args []string
		count := len(num)
		for i := 1; i <= count; i++ {
			args = append(args, fmt.Sprintf("$%v", i))
		}
		arguements := strings.Join(args, ",")
		queryString = query + "VALUES(" + arguements + ")"
		fmt.Println("insert query is :", queryString)
	} else if strings.ToUpper(queryFields[0]) == "UPDATE" {
		queryString = query
	}
	//fmt.Println("Update query is :", queryString)
	result, err := con.Db.Exec(queryString, num...)
	if err != nil {
		return &ModifyDb{-1, err.Error()}
	}
	re, err := result.RowsAffected()
	if err != nil {
		return &ModifyDb{-1, err.Error()}
	}
	return &ModifyDb{re, ""}
}

//DBInsertReturn inserts into a row and Returns any specified column value
//e.g the primary key if succesfull or error if not successful
func DBInsertReturn(query string, num ...interface{}) (interface{}, error) {
	con, err := OpenConnection()
	if err != nil {
		//fmt.Println("connection failed:", err)
		return nil, err
	}
	defer con.Close()
	var retSlice []string
	var res interface{}
	if c := strings.Contains(query, "RETURNING"); c {
		retSlice = strings.SplitN(query, "RETURNING", 2)
	}
	if c := strings.Contains(query, "returning"); c {
		retSlice = strings.SplitN(query, "returning", 2)
	}
	queryFields := strings.Fields(query)
	var queryString string
	if strings.ToUpper(queryFields[0]) == "INSERT" {
		var args []string
		count := len(num)
		for i := 1; i <= count; i++ {
			args = append(args, fmt.Sprintf("$%v", i))
		}
		arguements := strings.Join(args, ",")
		queryString = retSlice[0] + "VALUES(" + arguements + ") RETURNING" + retSlice[1]
		//	fmt.Println("insert query is :", queryString)
	} /*else if strings.ToUpper(queryFields[0]) == "UPDATE" {
		queryString = query
		fmt.Println("Update query is :", queryString)
	}*/
	err = con.Db.QueryRow(queryString, num...).Scan(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//DbSelect returns a column value selected from the database or error
func (con *DbCon) DBSelect(query string, num ...interface{}) (result interface{}, err error) {
	var res interface{}
	value := con.Db.QueryRow(query, num...)
	err = value.Scan(&res)

	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	switch selected := res.(type) {
	case string:
		if selected == "" {
			return nil, errors.New("empty")
		}
		return selected, nil
	case int:
		return selected, nil
	case nil:
		return nil, errors.New("null")
	}
	return res, nil
}

//DbSelect returns a value selected from the database or error
func DBSelect(query string, num ...interface{}) (result interface{}, err error) {
	con, err := OpenConnection()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	var res interface{}
	value := con.Db.QueryRow(query, num...)
	err = value.Scan(&res)

	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	switch selected := res.(type) {
	case string:
		if selected == "" {
			return nil, errors.New("empty")
		}
		return selected, nil
	case int:
		return selected, nil
	case nil:
		return nil, errors.New("null")
	}
	return res, nil
}

//DBSelectNoErr returns a value selected from the database without checking if there is error
func DBSelectNoErr(query string, num ...interface{}) (result interface{}) {
	con, err := OpenConnection()
	if err != nil {
		//fmt.Println("connection failed:", err)
		return nil
	}
	defer con.Close()
	var res interface{}
	value := con.Db.QueryRow(query, num...)
	_ = value.Scan(&res)
	switch selected := res.(type) {
	case string:
		if selected == "" {
			return "empty"
		}
		return selected
	case int:
		return selected
	case nil:
		return "null"
	}
	return res
}

//ExecuteDBSqlSelect This returns the result of the db query
func ExecuteDBSqlSelect(query string) string {
	con, err := OpenConnection()
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()

	var result string
	row, err := con.Db.Query(query)
	if err != nil {
		fmt.Println("telegram_bot.go: Selection Failed Due To: ", err)
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&result)
		if err != nil {
			fmt.Println("Row Scan Failed Due To:\n", err)
		}

	}
	return result
}
