package main

import (
	"fmt"
	"net/http"

	h "github.com/richardsric/teloauth/helper"
	p "github.com/richardsric/teloauth/public"
)

func main() {

	//r := mux.NewRouter()
	http.HandleFunc("/GoogleSignUp", p.HandleGoogleSignUp)
	http.HandleFunc("/tpn", p.RouteCreateTradeProfile)
	http.HandleFunc("/tpe", p.RouteEditTradeProfile)
	http.HandleFunc("/tpi", p.RouteEditTradeProfileInstance)
	http.HandleFunc("/signup", p.RouteSignUp)
	http.HandleFunc("/blockchaincallback", p.BlockChainUpdateCallBack)
	http.HandleFunc("/callback", p.BlockChainUpdateCallBack)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	servicemsg := "<b>WebDB Service Alert!</b>\nWebDB web proxy service has just started."
	p.SendServiceStatusIM(servicemsg)
	http.ListenAndServe(":"+h.ServicePort+"", nil)

}

func init() {
	h.GetOAuthDefaults()
	var name = "iTradeCoin Telegram Auth Service"
	var version = "0.001 DEVEL"
	var developer = "iYochu Nig LTD"

	fmt.Println("App Name: ", name)
	fmt.Println("App Version: ", version)
	fmt.Println("Developer Name: ", developer)
	fmt.Println("Service Port: ", h.ServicePort)
	fmt.Println("iTradeCoin  Telegram Auth Service: Visit http://localhost to use")

}
