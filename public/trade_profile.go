package public

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	s "strings"
	"time"

	h "github.com/richardsric/teloauth/helper"
)

// RouteCreateTradeProfile is use to edit trade profile
func RouteCreateTradeProfile(w http.ResponseWriter, r *http.Request) {
	var data tradeProfileData
	data.InheritSubDesc = "Do Not Inherit"
	var subTpList, jsonExchanges, jsonBaseMarkets, jsonTradeTypes string
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("RouteCreateTradeProfile. Failed to Open DB connection:", err)
		return
	}
	defer con.Close()
	if r.Method == "GET" {

		code := r.FormValue("code")

		if code == "" { // check if code is empty
			errdata := message{
				Error:       "Sorry you do not have access to the this page.. !!!",
				Information: "Go back to telegram to Generate link to be able to access this page",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//	t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, errdata)
			return
		}

		// decode the base64 and get the actual code
		codeDecoded := deCodeBase64(code)
		if checkCode("trade_profiles_new", codeDecoded) == false { // check for exist or expire
			errdata := message{
				Error:       "Sorry the page you are trying to view is no longer available... or your link have expired !!!",
				Information: "Go back to telegram to Generate New link. Links are use only once.",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")
			t.Execute(w, errdata)
			return
		}
		selsql := `SELECT trade_profile,trade_type,trade_mode,stop_loss,profit_lock_start,wallet_exposure,exchange_id,
buy_order_timeout,profit_keep,sell_trigger,inherit_subscribers_from,min_cap,max_cap,base_market,partialbuy_timeout,
partialbuy_timeout_pl,sell_order_timeout,profile_privacy,profit_keep_readjust_pl,profit_keep_readjust,sell_trigger_readjust,
trade_commission,tp_owner_account_id,tp_history,json_base_markets,json_exchanges,json_trade_types FROM trade_profiles_new WHERE code = $1 LIMIT 1`
		err := con.Db.QueryRow(selsql, codeDecoded).Scan(&data.TradeProfile, &data.TradeType, &data.TradeMode, &data.StopLoss, &data.ProfitLockStart,
			&data.WalletExposure, &data.ExchangeID, &data.BuyOrderTimeout, &data.ProfitKeep, &data.SellTrigger, &data.InheritSubscribersFrom,
			&data.MinCap, &data.MaxCap, &data.BaseMarket, &data.PartialBuyTimeout, &data.PartialBuyTimeoutPl, &data.SellOrderTimeout,
			&data.ProfilePrivacy, &data.ProfitKeepReadjustPl, &data.ProfitKeepReadjust, &data.SellTriggerReadjust, &data.TradeCommission,
			&data.TpOwnerID, &subTpList, &jsonBaseMarkets, &jsonExchanges, &jsonTradeTypes)

		if err != nil {
			fmt.Println("RouteCreateTradeProfile profile selection with code failed:", err)
			errdata := message{
				Error:       "Sorry the page you are trying to view is no longer available... or your link have expired !!! " + fmt.Sprintf("%v", err),
				Information: "Go back to telegram to Generate New link. Links are use only once.",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")
			t.Execute(w, errdata)
			return
		}
		//set the code
		data.Code = codeDecoded
		bMarket := getJsonBaseMarkets(jsonBaseMarkets)
		types := getJsonTradeTypes(jsonTradeTypes)
		exchanges := getJsonExchanges(jsonExchanges)
		subTp := getJsonTpHistory(subTpList)
		data.InheritSubDesc = subTp[int(data.InheritSubscribersFrom)]
		if data.InheritSubDesc == "" {
			data.InheritSubDesc = "Do Not Inherit"
			data.InheritSubscribersFrom = 0
		}
		pgdata := pageMainData{
			Title:          "Create New Trade Profile",
			TradeProfile:   data,
			BaseMarketData: bMarket,
			TradeTypesData: types,
			ExchangeData:   exchanges,
			SubscriberTp:   subTp,
		}

		// template from db
		templ := getTemplate("trade_profile_new") // select the template name from db
		//fmt.Println("Gotten trade_profile_new Template From Db", templ)
		t := template.New("Trade_profile_new")
		t, _ = t.Parse(templ)

		t.Execute(w, pgdata)
		return

	}

	if r.Method == "POST" {
		codePost := r.FormValue("code")
		//fmt.Println(codeBase64)
		//code, _ := base64.StdEncoding.DecodeString(codeBase64)
		tradeProfile := r.FormValue("tradeProfile")
		//fmt.Println("tradeProfile", tradeProfile)
		tradeProfile = s.Replace(s.ToLower(tradeProfile), " ", "", -1)
		tradeProfile = s.Replace(tradeProfile, "_", "", -1)
		tradeProfile = s.Replace(tradeProfile, "-", "", -1)
		tradeProfile = s.Replace(tradeProfile, "@", "", -1)
		tradeProfile = s.Replace(tradeProfile, "*", "", -1)
		tradeProfile = s.Replace(tradeProfile, "/", "", -1)
		tradeProfile = s.Replace(tradeProfile, "\\", "", -1)
		tradeProfile = s.Replace(tradeProfile, "&", "", -1)
		tradeProfile = s.Replace(tradeProfile, "$", "", -1)
		tradeProfile = s.Replace(tradeProfile, "#", "", -1)
		tradeProfile = s.Replace(tradeProfile, "(", "", -1)
		tradeProfile = s.Replace(tradeProfile, ")", "", -1)
		tradeProfile = s.Replace(tradeProfile, "!", "", -1)
		tradeProfile = s.Replace(tradeProfile, "~", "", -1)
		tradeProfile = s.Replace(tradeProfile, "`", "", -1)
		tradeProfile = s.Replace(tradeProfile, "?", "", -1)
		tradeProfile = s.Replace(tradeProfile, ",", "", -1)
		tradeProfile = s.Replace(tradeProfile, ".", "", -1)
		tradeProfile = s.Replace(tradeProfile, "'", "", -1)
		tradeType := r.FormValue("tradeType")
		//fmt.Println("tradeType", tradeType)
		tradeMode := r.FormValue("tradeMode")
		//fmt.Println("tradeMode", tradeMode)
		tradeCommission := r.FormValue("tradeCommission")
		//fmt.Println("tradeCommission", tradeCommission)
		profilePrivacy := r.FormValue("profilePrivacy")
		//fmt.Println("profilePrivacy", profilePrivacy)
		exchangeID := r.FormValue("exchange")
		//fmt.Println("exchangeID", exchangeID)
		baseMarket := r.FormValue("baseMarket")
		//fmt.Println("baseMarket", baseMarket)
		walletExposure := r.FormValue("walletExposure")
		//fmt.Println("walletExposure", walletExposure)
		minCap := r.FormValue("minCap")
		//fmt.Println("minCap", minCap)
		maxCap := r.FormValue("maxCap")
		//fmt.Println("maxCap", maxCap)
		profitLockStart := r.FormValue("profitLockStart")
		//fmt.Println("profitLockStart", profitLockStart)
		profitKeep := r.FormValue("profitKeep")
		//fmt.Println("profitKeep", profitKeep)
		profitKeepReadJstPL := r.FormValue("profitKeepReadJstPL")
		//fmt.Println("profitKeepReadJstPL", profitKeepReadJstPL)
		profitKeepReadjust := r.FormValue("profitKeepReadjust")
		//fmt.Println("ProfitKeepReadjust", profitKeepReadjust)
		buyOrderTimeout := r.FormValue("buyOrderTimeout")
		//fmt.Println("buyOrderTimeout", buyOrderTimeout)
		partialBuyTimeout := r.FormValue("partialBuyTimeout")
		//fmt.Println("partialBuyTimeout", partialBuyTimeout)
		partialBuyTimeoutPl := r.FormValue("partialBuyTimeoutPl")
		//fmt.Println("partialBuyTimeoutPl", partialBuyTimeoutPl)
		sellTrigger := r.FormValue("sellTrigger")
		//fmt.Println("sellTrigger", sellTrigger)
		sellOrderTimeout := r.FormValue("sellOrderTimeout")
		//fmt.Println("sellOrderTimeout", sellOrderTimeout)
		sellTriggerReadjust := r.FormValue("sellTriggerReadjust")
		//fmt.Println("sellTriggerReadjust", sellTriggerReadjust)
		stopLoss := r.FormValue("stopLoss")
		//fmt.Println("stopLoss", stopLoss)
		inheritSubscribersFrom := r.FormValue("inheritSubscribersFrom")
		//fmt.Println("inheritSubscribersFrom", inheritSubscribersFrom)
		//tpOwnerID := r.FormValue("tpOwnerID")
		//fmt.Println("tpOwnerID", tpOwnerID)
		submitButtonValue := r.FormValue("action")
		//fmt.Println("submitButtonValue", submitButtonValue)
		var finalized int
		var finalizedOn time.Time
		if submitButtonValue == "SAVE" {
			finalized = 0
			finalizedOn = time.Now()

		} else {
			finalized = 1
			finalizedOn = time.Now()
		}

		//fmt.Println("finalized", finalized)
		//fmt.Println("finalizedOn", finalizedOn)
		insertq := `UPDATE trade_profiles_new SET trade_profile = $1,trade_type = $2,trade_mode=$3,stop_loss=$4,profit_lock_start=$5,
wallet_exposure=$6,exchange_id=$7,buy_order_timeout=$8,profit_keep=$9,sell_trigger=$10,inherit_subscribers_from=$11,
min_cap=$12,max_cap=$13,base_market=$14,partialbuy_timeout=$15,partialbuy_timeout_pl=$16,sell_order_timeout=$17,profile_privacy=$18,
profit_keep_readjust_pl=$19,profit_keep_readjust=$20,sell_trigger_readjust=$21,trade_commission=$22,finalized=$24,finalized_on=$25, initiated_on=now()::timestamp WHERE code = $23
 RETURNING code`
		var retcode string
		err := con.Db.QueryRow(insertq, tradeProfile, tradeType, tradeMode, stopLoss, profitLockStart, walletExposure, exchangeID, buyOrderTimeout, profitKeep, sellTrigger,
			inheritSubscribersFrom, minCap, maxCap, baseMarket, partialBuyTimeout, partialBuyTimeoutPl, sellOrderTimeout, profilePrivacy,
			profitKeepReadJstPL, profitKeepReadjust, sellTriggerReadjust, tradeCommission, codePost, finalized, finalizedOn).Scan(&retcode)
		//	fmt.Println(insert)
		//if finalized is 0 then go back to the form to keep modification
		//	fmt.Println("ReCode = ", retcode, " Posted Code =", codePost)
		//	fmt.Println("Error is: ", err)
		if err == nil && (retcode == codePost) && finalized == 1 { /// insert was successful give user success message. and send message.

			datamsg := message{
				Sucess:      "Congratulation..... you have successfully Created" + tradeProfile + "  trade profile !!!",
				Information: " you can Go back to telegram to start using your editted profile.",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("success")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, datamsg)
			// send the user message on telegram.
			//sendMessageToTelegram(tpOwnerID)
		} else if err == nil && (retcode == codePost) && finalized == 0 {
			//Redirect user to new page to continue modification.
			//Do not complete until finalized

			selsql := `SELECT trade_profile,trade_type,trade_mode,stop_loss,profit_lock_start,wallet_exposure,exchange_id,
buy_order_timeout,profit_keep,sell_trigger,inherit_subscribers_from,min_cap,max_cap,base_market,partialbuy_timeout,
partialbuy_timeout_pl,sell_order_timeout,profile_privacy,profit_keep_readjust_pl,profit_keep_readjust,sell_trigger_readjust,
trade_commission,tp_owner_account_id,tp_history,json_base_markets,json_exchanges,json_trade_types FROM trade_profiles_new WHERE code = $1 LIMIT 1`
			err := con.Db.QueryRow(selsql, codePost).Scan(&data.TradeProfile, &data.TradeType, &data.TradeMode, &data.StopLoss, &data.ProfitLockStart,
				&data.WalletExposure, &data.ExchangeID, &data.BuyOrderTimeout, &data.ProfitKeep, &data.SellTrigger, &data.InheritSubscribersFrom,
				&data.MinCap, &data.MaxCap, &data.BaseMarket, &data.PartialBuyTimeout, &data.PartialBuyTimeoutPl, &data.SellOrderTimeout,
				&data.ProfilePrivacy, &data.ProfitKeepReadjustPl, &data.ProfitKeepReadjust, &data.SellTriggerReadjust, &data.TradeCommission,
				&data.TpOwnerID, &subTpList, &jsonBaseMarkets, &jsonExchanges, &jsonTradeTypes)

			if err != nil {
				fmt.Println("RouteCreateTradeProfile profile selection with code failed:", err)
				errdata := message{
					Error:       "Sorry, your data could not be loaded for futher modifications!!! " + fmt.Sprintf("%v", err),
					Information: "Go back to telegram to Generate New link. Links are use only once.",
				}

				// template from db
				templ := getTemplate("message") // select the template name from db
				t := template.New("error")
				t, _ = t.Parse(templ)

				//template from file
				//t, _ := template.ParseFiles("static/message_page.html")
				t.Execute(w, errdata)
				return
			}

			bMarket := getJsonBaseMarkets(jsonBaseMarkets)
			types := getJsonTradeTypes(jsonTradeTypes)
			exchanges := getJsonExchanges(jsonExchanges)
			subTp := getJsonTpHistory(subTpList)
			//set the chosen value of inheritance
			data.InheritSubDesc = subTp[int(data.InheritSubscribersFrom)]
			if data.InheritSubDesc == "" {
				data.InheritSubDesc = "Do Not Inherit"
				data.InheritSubscribersFrom = 0
			}
			//assign the submitted code to the new data
			data.Code = codePost
			pgdata := pageMainData{
				Title:          "New Trade Profile. Last Update Successful",
				TradeProfile:   data,
				BaseMarketData: bMarket,
				TradeTypesData: types,
				ExchangeData:   exchanges,
				SubscriberTp:   subTp,
			}

			// template from db
			templ := getTemplate("trade_profile_new") // select the template name from db
			//fmt.Println("Gotten trade_profile_new Template From Db", templ)
			t := template.New("Trade_profile_new")
			t, _ = t.Parse(templ)

			t.Execute(w, pgdata)
			return

		} else { /// TP Update failedfmt.Println("RouteEditTradeProfile. TP Update with code failed:", err)
			fmt.Println("RouteCreateTradeProfile. TP Update with code failed:", err)
			errdata := message{
				Error:       "Sorry your trade profile failed to create: " + fmt.Sprintf("%v", err),
				Information: "Please refresh the page and cross check your entry...........",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, errdata)
		}

	}

}

// RouteEditTradeProfile is use to edit trade profile
func RouteEditTradeProfile(w http.ResponseWriter, r *http.Request) {
	var data tradeProfileData
	data.InheritSubDesc = "Do Not Inherit"
	var subTpList, jsonExchanges, jsonBaseMarkets, jsonTradeTypes string
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("RouteEditTradeProfile. Failed to Open DB connection:", err)
		return
	}
	defer con.Close()
	if r.Method == "GET" {

		code := r.FormValue("code")

		if code == "" { // check if code is empty
			errdata := message{
				Error:       "Sorry you do not have access to the this page.. !!!",
				Information: "Go back to telegram to Generate link to be able to access this page",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//	t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, errdata)
			return
		}

		// decode the base64 and get the actual code
		codeDecoded := deCodeBase64(code)
		if checkCode("trade_profiles", codeDecoded) == false { // check for exist or expire
			errdata := message{
				Error:       "Sorry the page you are trying to view is no longer available... or your link have expired !!!",
				Information: "Go back to telegram to Generate New link. Links are use only once.",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")
			t.Execute(w, errdata)
			return
		}
		selsql := `SELECT trade_profile,trade_type,trade_mode,stop_loss,profit_lock_start,wallet_exposure,exchange_id,
buy_order_timeout,profit_keep,sell_trigger,inherit_subscribers_from,min_cap,max_cap,base_market,partialbuy_timeout,
partialbuy_timeout_pl,sell_order_timeout,profile_privacy,profit_keep_readjust_pl,profit_keep_readjust,sell_trigger_readjust,
trade_commission,tp_owner_account_id,tp_history,json_base_markets,json_exchanges,json_trade_types FROM trade_profiles WHERE finalized = 0 AND code = $1 LIMIT 1`
		err := con.Db.QueryRow(selsql, codeDecoded).Scan(&data.TradeProfile, &data.TradeType, &data.TradeMode, &data.StopLoss, &data.ProfitLockStart,
			&data.WalletExposure, &data.ExchangeID, &data.BuyOrderTimeout, &data.ProfitKeep, &data.SellTrigger, &data.InheritSubscribersFrom,
			&data.MinCap, &data.MaxCap, &data.BaseMarket, &data.PartialBuyTimeout, &data.PartialBuyTimeoutPl, &data.SellOrderTimeout,
			&data.ProfilePrivacy, &data.ProfitKeepReadjustPl, &data.ProfitKeepReadjust, &data.SellTriggerReadjust, &data.TradeCommission,
			&data.TpOwnerID, &subTpList, &jsonBaseMarkets, &jsonExchanges, &jsonTradeTypes)

		if err != nil {
			fmt.Println("RouteEditTradeProfile: profile selection with code failed:", err)
			errdata := message{
				Error:       "Sorry the page you are trying to view is no longer available... or your link have expired !!! " + fmt.Sprintf("%v", err),
				Information: "Go back to telegram to Generate New link. Links are use only once.",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")
			t.Execute(w, errdata)
			return
		}
		//set the code
		data.Code = codeDecoded
		bMarket := getJsonBaseMarkets(jsonBaseMarkets)
		types := getJsonTradeTypes(jsonTradeTypes)
		exchanges := getJsonExchanges(jsonExchanges)
		subTp := getJsonTpHistory(subTpList)
		//set the chosen value of inheritance
		data.InheritSubDesc = subTp[int(data.InheritSubscribersFrom)]
		if data.InheritSubDesc == "" {
			data.InheritSubDesc = "Do Not Inherit"
			data.InheritSubscribersFrom = 0
		}

		pgdata := pageMainData{
			Title:          "Modify Trade Profile",
			TradeProfile:   data,
			BaseMarketData: bMarket,
			TradeTypesData: types,
			ExchangeData:   exchanges,
			SubscriberTp:   subTp,
		}

		// template from db
		templ := getTemplate("trade_profile_edit") // select the template name from db
		//fmt.Println("Gotten trade_profile_new Template From Db", templ)
		t := template.New("Trade_profile_edit")
		t, _ = t.Parse(templ)

		t.Execute(w, pgdata)
		return

	}

	if r.Method == "POST" {
		codePost := r.FormValue("code")
		//fmt.Println(codeBase64)
		//code, _ := base64.StdEncoding.DecodeString(codeBase64)
		tradeProfile := r.FormValue("tradeProfile")
		//fmt.Println("tradeProfile", tradeProfile)
		tradeType := r.FormValue("tradeType")
		//fmt.Println("tradeType", tradeType)
		tradeMode := r.FormValue("tradeMode")
		//fmt.Println("tradeMode", tradeMode)
		tradeCommission := r.FormValue("tradeCommission")
		//fmt.Println("tradeCommission", tradeCommission)
		profilePrivacy := r.FormValue("profilePrivacy")
		//fmt.Println("profilePrivacy", profilePrivacy)
		exchangeID := r.FormValue("exchange")
		//fmt.Println("exchangeID", exchangeID)
		baseMarket := r.FormValue("baseMarket")
		//fmt.Println("baseMarket", baseMarket)
		walletExposure := r.FormValue("walletExposure")
		//fmt.Println("walletExposure", walletExposure)
		minCap := r.FormValue("minCap")
		//fmt.Println("minCap", minCap)
		maxCap := r.FormValue("maxCap")
		//fmt.Println("maxCap", maxCap)
		profitLockStart := r.FormValue("profitLockStart")
		//fmt.Println("profitLockStart", profitLockStart)
		profitKeep := r.FormValue("profitKeep")
		//fmt.Println("profitKeep", profitKeep)
		profitKeepReadJstPL := r.FormValue("profitKeepReadJstPL")
		//fmt.Println("profitKeepReadJstPL", profitKeepReadJstPL)
		profitKeepReadjust := r.FormValue("profitKeepReadjust")
		//fmt.Println("ProfitKeepReadjust", profitKeepReadjust)
		buyOrderTimeout := r.FormValue("buyOrderTimeout")
		//fmt.Println("buyOrderTimeout", buyOrderTimeout)
		partialBuyTimeout := r.FormValue("partialBuyTimeout")
		//fmt.Println("partialBuyTimeout", partialBuyTimeout)
		partialBuyTimeoutPl := r.FormValue("partialBuyTimeoutPl")
		//fmt.Println("partialBuyTimeoutPl", partialBuyTimeoutPl)
		sellTrigger := r.FormValue("sellTrigger")
		//fmt.Println("sellTrigger", sellTrigger)
		sellOrderTimeout := r.FormValue("sellOrderTimeout")
		//fmt.Println("sellOrderTimeout", sellOrderTimeout)
		sellTriggerReadjust := r.FormValue("sellTriggerReadjust")
		//fmt.Println("sellTriggerReadjust", sellTriggerReadjust)
		stopLoss := r.FormValue("stopLoss")
		//fmt.Println("stopLoss", stopLoss)
		inheritSubscribersFrom := r.FormValue("inheritSubscribersFrom")
		//fmt.Println("inheritSubscribersFrom", inheritSubscribersFrom)
		//tpOwnerID := r.FormValue("tpOwnerID")
		//fmt.Println("tpOwnerID", tpOwnerID)
		submitButtonValue := r.FormValue("action")
		//fmt.Println("submitButtonValue", submitButtonValue)
		var finalized int
		var finalizedOn time.Time
		if submitButtonValue == "SAVE" {
			finalized = 0
			finalizedOn = time.Now()

		} else {
			finalized = 1
			finalizedOn = time.Now()
		}

		//fmt.Println("finalized", finalized)
		//fmt.Println("finalizedOn", finalizedOn)
		insertq := `UPDATE trade_profiles SET code = $1,trade_type = $2,trade_mode=$3,stop_loss=$4,profit_lock_start=$5,
wallet_exposure=$6,exchange_id=$7,buy_order_timeout=$8,profit_keep=$9,sell_trigger=$10,inherit_subscribers_from=$11,
min_cap=$12,max_cap=$13,base_market=$14,partialbuy_timeout=$15,partialbuy_timeout_pl=$16,sell_order_timeout=$17,profile_privacy=$18,
profit_keep_readjust_pl=$19,profit_keep_readjust=$20,sell_trigger_readjust=$21,trade_commission=$22,finalized=$24,finalized_on=$25, initiated_on=now()::timestamp WHERE code = $23
 RETURNING code`
		var retcode string
		//update with code so as not to change the numbering of d params.
		err := con.Db.QueryRow(insertq, codePost, tradeType, tradeMode, stopLoss, profitLockStart, walletExposure, exchangeID, buyOrderTimeout, profitKeep, sellTrigger,
			inheritSubscribersFrom, minCap, maxCap, baseMarket, partialBuyTimeout, partialBuyTimeoutPl, sellOrderTimeout, profilePrivacy,
			profitKeepReadJstPL, profitKeepReadjust, sellTriggerReadjust, tradeCommission, codePost, finalized, finalizedOn).Scan(&retcode)
		//	fmt.Println(insert)
		//if finalized is 0 then go back to the form to keep modification
		fmt.Println("ReCode = ", retcode, " Posted Code =", codePost)
		fmt.Println("Error is: ", err)
		if err == nil && (retcode == codePost) && finalized == 1 { /// insert was successful give user success message. and send message.

			datamsg := message{
				Sucess:      "Congratulation..... you have successfully Modified" + tradeProfile + "  trade profile !!!",
				Information: " you can Go back to telegram to start using your editted profile.",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("success")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, datamsg)
			// send the user message on telegram.
			//sendMessageToTelegram(tpOwnerID)
		} else if err == nil && (retcode == codePost) && finalized == 0 {
			//Redirect user to new page to continue modification.
			//Do not complete until finalized

			selsql := `SELECT trade_profile,trade_type,trade_mode,stop_loss,profit_lock_start,wallet_exposure,exchange_id,
buy_order_timeout,profit_keep,sell_trigger,inherit_subscribers_from,min_cap,max_cap,base_market,partialbuy_timeout,
partialbuy_timeout_pl,sell_order_timeout,profile_privacy,profit_keep_readjust_pl,profit_keep_readjust,sell_trigger_readjust,
trade_commission,tp_owner_account_id,tp_history,json_base_markets,json_exchanges,json_trade_types FROM trade_profiles WHERE finalized = 0 AND code = $1 LIMIT 1`
			err := con.Db.QueryRow(selsql, codePost).Scan(&data.TradeProfile, &data.TradeType, &data.TradeMode, &data.StopLoss, &data.ProfitLockStart,
				&data.WalletExposure, &data.ExchangeID, &data.BuyOrderTimeout, &data.ProfitKeep, &data.SellTrigger, &data.InheritSubscribersFrom,
				&data.MinCap, &data.MaxCap, &data.BaseMarket, &data.PartialBuyTimeout, &data.PartialBuyTimeoutPl, &data.SellOrderTimeout,
				&data.ProfilePrivacy, &data.ProfitKeepReadjustPl, &data.ProfitKeepReadjust, &data.SellTriggerReadjust, &data.TradeCommission,
				&data.TpOwnerID, &subTpList, &jsonBaseMarkets, &jsonExchanges, &jsonTradeTypes)

			if err != nil {
				fmt.Println("RouteEditTradeProfile profile selection with code failed:", err)
				errdata := message{
					Error:       "Sorry, your data could not be loaded for futher modifications!!! " + fmt.Sprintf("%v", err),
					Information: "Go back to telegram to Generate New link. Links are use only once.",
				}

				// template from db
				templ := getTemplate("message") // select the template name from db
				t := template.New("error")
				t, _ = t.Parse(templ)

				//template from file
				//t, _ := template.ParseFiles("static/message_page.html")
				t.Execute(w, errdata)
				return
			}

			bMarket := getJsonBaseMarkets(jsonBaseMarkets)
			types := getJsonTradeTypes(jsonTradeTypes)
			exchanges := getJsonExchanges(jsonExchanges)
			subTp := getJsonTpHistory(subTpList)
			//set the chosen value of inheritance
			data.InheritSubDesc = subTp[int(data.InheritSubscribersFrom)]
			if data.InheritSubDesc == "" {
				data.InheritSubDesc = "Do Not Inherit"
				data.InheritSubscribersFrom = 0
			}
			//assign the submitted code to the new data
			data.Code = codePost
			pgdata := pageMainData{
				Title:          "Modify Trade Profile. Last Update Successful",
				TradeProfile:   data,
				BaseMarketData: bMarket,
				TradeTypesData: types,
				ExchangeData:   exchanges,
				SubscriberTp:   subTp,
			}

			// template from db
			templ := getTemplate("trade_profile_edit") // select the template name from db
			//fmt.Println("Gotten trade_profile_new Template From Db", templ)
			t := template.New("Trade_profile_edit")
			t, _ = t.Parse(templ)

			t.Execute(w, pgdata)
			return

		} else { /// TP Update failedfmt.Println("RouteEditTradeProfile. TP Update with code failed:", err)
			fmt.Println("RouteEditTradeProfile. TP Update with code failed:", err)
			errdata := message{
				Error:       "Sorry your trade profile failed to Update: " + fmt.Sprintf("%v", err),
				Information: "Please refresh the page and cross check your entry...........",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, errdata)
			return
		}

	}

}

// RouteEditTradeProfileInstance is use to edit trade profile instance
func RouteEditTradeProfileInstance(w http.ResponseWriter, r *http.Request) {
	var tpidata tradeProfileData
	var subTpList, jsonExchanges, jsonBaseMarkets, jsonTradeTypes string
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("RouteEditTradeProfileInstance. Failed to Open DB connection:", err)
		return
	}
	defer con.Close()
	if r.Method == "GET" {
		code := r.FormValue("code")

		if code == "" { // check if code is empty
			data := message{
				Error:       "Sorry you do not have access to the this page.. !!!",
				Information: "Go back to telegram to Generate link to be able to access this page",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, data)
			return
		}
		// decode the base64 and get the actual code
		codeDecoded := deCodeBase64(code)
		if checkCode("trade_profile_instances", codeDecoded) == false { // check for exist or expire
			data := message{
				Error:       "Sorry the page you are trying to view is no longer available... or your link have expired !!!",
				Information: "Go back to telegram to Generate New link. Links are use only once.",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, data)
			return
		}
		sqlsel := `SELECT trade_profile,market,trade_type,trade_mode,stop_loss,profit_lock_start,wallet_exposure,exchange_id,
					buy_order_timeout,profit_keep,sell_trigger,inherit_subscribers_from,min_cap,max_cap,base_market,partialbuy_timeout,
					partialbuy_timeout_pl,sell_order_timeout,profile_privacy,profit_keep_readjust_pl,profit_keep_readjust,sell_trigger_readjust,
					profile_id,tp_owner_account_id,tp_history,json_base_markets,json_exchanges,json_trade_types FROM trade_profile_instances WHERE code = $1 AND finalized = 0`
		err := con.Db.QueryRow(sqlsel, codeDecoded).Scan(&tpidata.TradeProfile, &tpidata.Market, &tpidata.TradeType, &tpidata.TradeMode,
			&tpidata.StopLoss, &tpidata.ProfitLockStart, &tpidata.WalletExposure, &tpidata.ExchangeID, &tpidata.BuyOrderTimeout, &tpidata.ProfitKeep,
			&tpidata.SellTrigger, &tpidata.InheritSubscribersFrom, &tpidata.MinCap, &tpidata.MaxCap, &tpidata.BaseMarket, &tpidata.PartialBuyTimeout,
			&tpidata.PartialBuyTimeoutPl, &tpidata.SellOrderTimeout, &tpidata.ProfilePrivacy, &tpidata.ProfitKeepReadjustPl, &tpidata.ProfitKeepReadjust,
			&tpidata.SellTriggerReadjust, &tpidata.ProfileID, &tpidata.TpOwnerID, &subTpList, &jsonBaseMarkets, &jsonExchanges, &jsonTradeTypes)

		if err != nil {
			fmt.Println("RouteEditTradeProfileInstance. TPI selection failed:", err)
			data := message{
				Error:       "Sorry the page you are trying to view is no longer available... or your link have expired !!!",
				Information: "Go back to telegram to Generate New link. Links are use only once.",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, data)
			return
		}

		bMarket := getJsonBaseMarkets(jsonBaseMarkets)
		types := getJsonTradeTypes(jsonTradeTypes)
		exchanges := getJsonExchanges(jsonExchanges)
		//subTp := getJsonTpHistory(subTpList)
		//assign the code
		tpidata.Code = codeDecoded
		pgdata := pageMainData{
			Title:          "iTradeCoin: Trade Execution Instance",
			TradeProfile:   tpidata,
			BaseMarketData: bMarket,
			TradeTypesData: types,
			ExchangeData:   exchanges,
			//SubscriberTp:   subTp,
		}

		// template from db
		templ := getTemplate("trade_profile_instance") // select the template name from db
		t := template.New("Trade_profile_instance")
		t, _ = t.Parse(templ)

		//template from file
		//fmt.Println("the value is", data)
		//t, _ := template.ParseFiles("static/instance.html")

		t.Execute(w, pgdata)
		return
	}

	if r.Method == "POST" {
		codePost := r.FormValue("code")
		//fmt.Println(codeBase64)
		//code, _ := base64.StdEncoding.DecodeString(codeBase64)
		tradeProfile := r.FormValue("tradeProfile")

		profitLockStart := r.FormValue("profitLockStart")
		//fmt.Println("profitLockStart", profitLockStart)
		profitKeep := r.FormValue("profitKeep")
		//fmt.Println("profitKeep", profitKeep)
		profitKeepReadJstPL := r.FormValue("profitKeepReadJstPL")
		//fmt.Println("profitKeepReadJstPL", profitKeepReadJstPL)
		profitKeepReadjust := r.FormValue("profitKeepReadjust")
		//fmt.Println("ProfitKeepReadjust", profitKeepReadjust)
		buyOrderTimeout := r.FormValue("buyOrderTimeout")
		//fmt.Println("buyOrderTimeout", buyOrderTimeout)
		partialBuyTimeout := r.FormValue("partialBuyTimeout")
		//fmt.Println("partialBuyTimeout", partialBuyTimeout)
		partialBuyTimeoutPl := r.FormValue("partialBuyTimeoutPl")
		//fmt.Println("partialBuyTimeoutPl", partialBuyTimeoutPl)
		sellTrigger := r.FormValue("sellTrigger")
		//fmt.Println("sellTrigger", sellTrigger)
		sellOrderTimeout := r.FormValue("sellOrderTimeout")
		//fmt.Println("sellOrderTimeout", sellOrderTimeout)
		sellTriggerReadjust := r.FormValue("sellTriggerReadjust")
		//fmt.Println("sellTriggerReadjust", sellTriggerReadjust)
		stopLoss := r.FormValue("stopLoss")
		//fmt.Println("stopLoss", stopLoss)

		//fmt.Println("inheritSubscribersFrom", inheritSubscribersFrom)
		//tpOwnerID := r.FormValue("tpOwnerID")
		//fmt.Println("tpOwnerID", tpOwnerID)
		submitButtonValue := r.FormValue("action")
		//fmt.Println("submitButtonValue", submitButtonValue)
		finalized := 1
		finalizedOn := time.Now()
		if submitButtonValue == "SAVE" {
			finalized = 0
		}
		//	fmt.Println("finalized", finalized)
		//	fmt.Println("finalizedOn", finalizedOn)
		tpiuq := `UPDATE trade_profile_instances SET stop_loss=$1,profit_lock_start=$2,buy_order_timeout=$3,profit_keep=$4,sell_trigger=$5,
partialbuy_timeout=$6,partialbuy_timeout_pl=$7,sell_order_timeout=$8,
profit_keep_readjust_pl=$9,profit_keep_readjust=$10,sell_trigger_readjust=$11,finalized=$12,finalized_on=$13 WHERE code = $14 RETURNING code`
		var retcode string
		err := con.Db.QueryRow(tpiuq,
			stopLoss, profitLockStart, buyOrderTimeout, profitKeep, sellTrigger,
			partialBuyTimeout, partialBuyTimeoutPl, sellOrderTimeout,
			profitKeepReadJstPL, profitKeepReadjust, sellTriggerReadjust, finalized, finalizedOn, codePost).Scan(&retcode)

		//fmt.Println(insert)
		if err == nil && (retcode == codePost) && finalized == 1 { /// insert was successful give user success message. and send message.

			data := message{
				Sucess:      "Congratulation..... you have successfully edited/updated " + tradeProfile + "  Execution Instance !!!",
				Information: " you can Go back to telegram to start using your edited profile.",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("success")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, data)
			// send the user message on telegram.
			//sendMessageToTelegram(tpOwnerID)
		} else if err == nil && (retcode == codePost) && finalized == 0 {
			//continue editing

			sqlsel := `SELECT trade_profile,market,trade_type,trade_mode,stop_loss,profit_lock_start,wallet_exposure,exchange_id,
			buy_order_timeout,profit_keep,sell_trigger,inherit_subscribers_from,min_cap,max_cap,base_market,partialbuy_timeout,
			partialbuy_timeout_pl,sell_order_timeout,profile_privacy,profit_keep_readjust_pl,profit_keep_readjust,sell_trigger_readjust,
			profile_id,tp_owner_account_id,tp_history,json_base_markets,json_exchanges,json_trade_types FROM trade_profile_instances WHERE code = $1 AND finalized = 0`
			err := con.Db.QueryRow(sqlsel, codePost).Scan(&tpidata.TradeProfile, &tpidata.Market, &tpidata.TradeType, &tpidata.TradeMode,
				&tpidata.StopLoss, &tpidata.ProfitLockStart, &tpidata.WalletExposure, &tpidata.ExchangeID, &tpidata.BuyOrderTimeout, &tpidata.ProfitKeep,
				&tpidata.SellTrigger, &tpidata.InheritSubscribersFrom, &tpidata.MinCap, &tpidata.MaxCap, &tpidata.BaseMarket, &tpidata.PartialBuyTimeout,
				&tpidata.PartialBuyTimeoutPl, &tpidata.SellOrderTimeout, &tpidata.ProfilePrivacy, &tpidata.ProfitKeepReadjustPl, &tpidata.ProfitKeepReadjust,
				&tpidata.SellTriggerReadjust, &tpidata.ProfileID, &tpidata.TpOwnerID, &subTpList, &jsonBaseMarkets, &jsonExchanges, &jsonTradeTypes)

			if err != nil {
				fmt.Println("RouteEditTradeProfileInstance. TPI selection failed:", err)
				data := message{
					Error:       "Sorry, could not reload the page for continued modification!!!",
					Information: "Go back to telegram to Generate New link. Links are use only once.",
				}

				// template from db
				templ := getTemplate("message") // select the template name from db
				t := template.New("error")
				t, _ = t.Parse(templ)

				//template from file
				//t, _ := template.ParseFiles("static/message_page.html")

				t.Execute(w, data)
				return
			}

			bMarket := getJsonBaseMarkets(jsonBaseMarkets)
			types := getJsonTradeTypes(jsonTradeTypes)
			exchanges := getJsonExchanges(jsonExchanges)
			//	subTp := getJsonTpHistory(subTpList)
			//assign the code
			tpidata.Code = codePost
			pgdata := pageMainData{
				Title:          "iTradeCoin: Trade Execution Instance: Your Changes Saved Successfully",
				TradeProfile:   tpidata,
				BaseMarketData: bMarket,
				TradeTypesData: types,
				ExchangeData:   exchanges,
				//SubscriberTp:   subTp,
			}

			// template from db
			templ := getTemplate("trade_profile_instance") // select the template name from db
			t := template.New("Trade_profile_instance")
			t, _ = t.Parse(templ)

			//template from file
			//fmt.Println("the value is", data)
			//t, _ := template.ParseFiles("static/instance.html")

			t.Execute(w, pgdata)
			return

		} else { /// TPI update failed
			fmt.Println("RouteEditTradeProfileInstance. TPI Update failed:", err)
			data := message{
				Error:       "Sorry Trade Execution Instance Data failed to load.",
				Information: "Please refresh the page and cross check your entry...........",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")

			t.Execute(w, data)
		}

	}

}

//DoTradeProfileOnWebDB performs db delete operation and return affected row
func DoTradeProfileOnWebDB(tableName, code string) int {
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("tradeprofile.go:DoTradeProfileNew():connection failed:", err)
		return 0
	}
	defer con.Close()

	qry, err := con.Db.Exec("DELETE FROM "+tableName+" WHERE code = $1 ", code)
	if err != nil {
		fmt.Println(" DELETE failed:", err)
		return 0
	}
	res, _ := qry.RowsAffected()

	return int(res)
}

func getTemplate(templateName string) string {
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("tradeprofile.go:getTemplate:connection failed:", err)
		return TemplateErrorPage()
	}
	var tpl string
	defer con.Close()
	sqlsel := `SELECT tpl FROM templates WHERE name = $1`
	err = con.Db.QueryRow(sqlsel, templateName).Scan(&tpl)
	if err != nil {
		fmt.Println("template selection failed:", err)
		return TemplateErrorPage()
	}
	return tpl
}

func checkCode(tableName, code string) bool {
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("tradeprofile.go:getTemplate:connection failed:", err)
		return false
	}
	var tpoid int64
	defer con.Close()
	sql := "SELECT tp_owner_account_id FROM " + tableName + " WHERE code = $1 AND finalized = 0"
	err = con.Db.QueryRow(sql, code).Scan(&tpoid)
	if err != nil {
		return false
	}
	if tpoid > 0 {
		return true
	}
	return false
}

func deCodeBase64(code string) string {
	//fmt.Println(code)
	codeDecode, _ := base64.StdEncoding.DecodeString(code)

	return string(codeDecode)
}

func sendMessageToTelegram(code string) {
	//fmt.Println("Enter To Send alert")
	aID := base64.StdEncoding.EncodeToString([]byte(code))
	url := fmt.Sprintf("%v/WebSiteCallBack?aID=%v", h.TelegramURL, aID)
	http.Get(url)
	//fmt.Println("Alert Sent......")
}
