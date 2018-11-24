package helper

type DefaultSettings struct {
	exchangeID int
	baseURL    string
	txnFee     float64
}
type ModifyDb struct {
	AffectedRows int64
	ErrorMsg     string
}

type RowSelect struct {
	Columns  map[string]interface{}
	ErrorMsg string
}
type dbInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}
