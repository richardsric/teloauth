package helper

import (
	"fmt"
	"os"
	"time"
)

//TelegramURL variable holds the TelegramURL from telegram_auth_settings table
var TelegramURL, GoogleCallBackURL string

//TxnFee variable holds the TxnFee from gateway_settings table
var TxnFee float64

//Port variable holds the service_port from gateway_settings table
var ServicePort string

//TimeOut variable holds the request_time_out from gateway_settings table
var TimeOut time.Duration

// Dsets is the slice that will hold the loaded default settings from db.
var Dsets []DefaultSettings

// GetOAuthDefaults is used to load setting from db at the start up.
func GetOAuthDefaults() {
	fmt.Println("Enter To Get TelOAuth settings From DB")
	gs := DBSelectRow("SELECT service_port,telegram_url FROM telegram_auth_settings")
	if gs.ErrorMsg != "" {
		fmt.Println("GetOAuthDefaults", gs.ErrorMsg)
		os.Exit(2)
	}
	gs1 := DBSelectRow("SELECT signup_url FROM callback_url LIMIT 1")
	if gs1.ErrorMsg != "" {
		fmt.Println("GetOAuthDefaults:", gs1.ErrorMsg)
		os.Exit(2)
	}
	GoogleCallBackURL = gs1.Columns["signup_url"].(string)
	ServicePort = gs.Columns["service_port"].(string)
	TelegramURL = gs.Columns["telegram_url"].(string)

}
