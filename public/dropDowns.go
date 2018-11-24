package public

import (
	"encoding/json"
	"fmt"

	h "github.com/richardsric/teloauth/helper"
)

func getJsonTpHistory(tpHistoryJson string) map[int]string {
	var jst interface{}
	tpHistory := make(map[int]string)

	err := json.Unmarshal([]byte(tpHistoryJson), &jst)
	if err != nil {
		fmt.Println(err)
		return nil

	}
	tpHistory[0] = "Do Not Inherit"
	for _, val := range jst.([]interface{}) {
		j1 := val.(map[string]interface{})["profile_id"].(float64)
		j2 := val.(map[string]interface{})["trade_profile"].(string)
		j3 := int(j1)

		//set tp history
		tpHistory[j3] = j2
	}

	return tpHistory
}

// GetSubscribeTradeProfile2 list all the user trade profile
func GetSubscribeTradeProfile2(tradeProfileJSON string) map[int]string {
	var jst interface{}
	tpHistory := make(map[int]string)

	err := json.Unmarshal([]byte(tradeProfileJSON), &jst)
	if err != nil {
		fmt.Println(err)
		return nil

	}
	//fmt.Printf("%T, %+v\n\n", jst, jst)
	//j := jst.([]interface{})
	tpHistory[0] = "Do Not Inherit"
	for _, val := range jst.([]interface{}) {
		//fmt.Printf("Type is: %T\n, Value is: %+v\n\n", val, val)
		j1 := val.(map[string]interface{})["profile_id"].(float64)
		j2 := val.(map[string]interface{})["trade_profile"].(string)

		j3 := int(j1)

		//set tp history
		tpHistory[j3] = j2

	}
	//fmt.Printf("\n\n\nProfile History is: %+v\n\n", tpHistory)
	//fmt.Printf("\n\n\nTphist Struct is: %+v\n\n", tphistStruct)
	return tpHistory
}

// getExchanges list all the exchanges
func getJsonExchanges(exchangesJson string) map[int]string {

	var jst interface{}
	exchangeData := make(map[int]string)

	err := json.Unmarshal([]byte(exchangesJson), &jst)
	if err != nil {
		fmt.Println(err)
		return nil

	}

	for _, val := range jst.([]interface{}) {
		j1 := val.(map[string]interface{})["exchange_id"].(float64)
		j2 := val.(map[string]interface{})["exchange_name"].(string)
		j3 := int(j1)

		//set tp history
		exchangeData[j3] = j2
	}

	return exchangeData
}

// getJsonBaseMarkets list all base markets
func getJsonBaseMarkets(baseMarketJson string) map[string]string {

	var jst interface{}
	bMarket := make(map[string]string)
	err := json.Unmarshal([]byte(baseMarketJson), &jst)
	if err != nil {

		fmt.Println(err)

	}
	for _, val := range jst.([]interface{}) {
		j2 := val.(map[string]interface{})["base_market"].(string)

		bMarket[j2] = j2

	}

	return bMarket
}

// getTradeTypes list all trades
func getJsonTradeTypes(tradeTypesJson string) map[int]string {
	var jst interface{}
	tradeTypes := make(map[int]string)
	err := json.Unmarshal([]byte(tradeTypesJson), &jst)
	if err != nil {

		fmt.Println(err)

	}
	for _, val := range jst.([]interface{}) {
		j1 := val.(map[string]interface{})["trade_type_id"].(float64)
		j2 := val.(map[string]interface{})["trade_type"].(string)
		j3 := int(j1)

		tradeTypes[j3] = j2

	}

	return tradeTypes
}

// getJsonCountryCode list all country code
func getJsonCountryCode(countryCodeJson string) map[string]string {
	var jst interface{}
	cCode := make(map[string]string)
	err := json.Unmarshal([]byte(countryCodeJson), &jst)
	if err != nil {

		fmt.Println(err)

	}
	for _, val := range jst.([]interface{}) {
		j1 := val.(map[string]interface{})["country"].(string)
		j2 := val.(map[string]interface{})["country_code"].(string)

		cCode[j2] = j1

	}

	return cCode
}

// getJsonAccountType list all trades
func getJsonAccountType(accountTypesJson string) map[int]string {
	var jst interface{}
	accountTypes := make(map[int]string)
	err := json.Unmarshal([]byte(accountTypesJson), &jst)
	if err != nil {

		fmt.Println(err)

	}
	for _, val := range jst.([]interface{}) {
		j1 := val.(map[string]interface{})["account_type_id"].(float64)
		j2 := val.(map[string]interface{})["account_type"].(string)
		j3 := int(j1)

		accountTypes[j3] = j2

	}

	return accountTypes
}

// GetExchanges list all the exchanges
func GetExchanges() allExchanges {

	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	var exchangeName string
	var exchangeID int
	exchangeData := make([]exchange, 0)
	row, err := con.Db.Query(`SELECT exchange_id,exchange_name FROM exapi_settings`)
	if err != nil {
		fmt.Println("dropDownStructs.go:GetExchanges(): Selection Failed Due To: ", err)
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&exchangeID, &exchangeName)
		if err != nil {
			fmt.Println("dropDownStructs.go:GetExchanges():Row Scan Failed Due To:\n", err)
		}
		result := exchange{
			ExchangeName: exchangeName,
			ExchnageID:   exchangeID,
		}
		exchangeData = append(exchangeData, result)
	}

	return allExchanges{
		Exchanges: exchangeData,
	}
}

// GetTradeTypes list all trades
func GetTradeTypes() allTradeTypes {

	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	var Name string
	var ID int
	tradeTypeData := make([]tradeType, 0)
	row, err := con.Db.Query(`SELECT trade_type_id,trade_type FROM trade_types`)
	if err != nil {
		fmt.Println("dropDownStructs.go:GetTradeTypes(): Selection Failed Due To: ", err)
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&ID, &Name)
		if err != nil {
			fmt.Println("dropDownStructs.go:GetTradeTypes():Row Scan Failed Due To:\n", err)
		}
		result := tradeType{
			TradeTypeName: Name,
			TradeTypeID:   ID,
		}
		tradeTypeData = append(tradeTypeData, result)
	}

	return allTradeTypes{
		TradeTypes: tradeTypeData,
	}
}

// GetBaseMarket list all base markets
func GetBaseMarket() allBaseMarket {

	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	var Name string
	//var ID int
	baseMarketData := make([]baseMarket, 0)
	row, err := con.Db.Query(`SELECT DISTINCT pcurrency FROM currency_pairs`)
	if err != nil {
		fmt.Println("dropDownStructs.go:GetBaseMarket(): Selection Failed Due To: ", err)
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&Name)
		if err != nil {
			fmt.Println("dropDownStructs.go:GetBaseMarket():Row Scan Failed Due To:\n", err)
		}
		result := baseMarket{
			baseMarketName: Name,
		}
		baseMarketData = append(baseMarketData, result)
	}

	return allBaseMarket{
		BaseMarkets: baseMarketData,
	}
}

func GetALLDropDownData() allDropDownData {
	exchanges := GetExchanges()
	tradeTypes := GetTradeTypes()
	baseMarkets := GetBaseMarket()

	return allDropDownData{
		Exchanges:   exchanges,
		TradeTypes:  tradeTypes,
		BaseMarkets: baseMarkets,
	}
}
