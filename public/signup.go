package public

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	h "github.com/richardsric/teloauth/helper"
)

//SelectOption Selects an option
func SelectOption(val string, option string) string {
	if val == option {
		return "selected"
	}
	return ""
}

// RouteSignUp is use to edit trade profile
func RouteSignUp(w http.ResponseWriter, r *http.Request) {

	Fmap := template.FuncMap{
		"SelectOption": SelectOption,
	}
	var data accountInfo
	data.CountryName = "NIGERIA"
	data.CountryCode = "234"
	var jsonCountryCode, jsonAccountTypes string
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("RouteSignUp. Failed to Open DB connection:", err)
		return
	}
	defer con.Close()
	if r.Method == "GET" {

		stateDecoded := r.FormValue("state")

		if stateDecoded == "" { // check if code is empty
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
		//stateDecoded := deCodeBase64(state)
		if checkState("account_signup", stateDecoded) == false { // check for exist or expire
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
		selsql := `SELECT json_country_code,json_trade_types,reg_type FROM account_signup WHERE state = $1 LIMIT 1`
		err := con.Db.QueryRow(selsql, stateDecoded).Scan(&jsonCountryCode, &jsonAccountTypes, &data.RegType)
		//	fmt.Println(err)
		if err != nil {
			fmt.Println("RouteSignUp profile selection with code failed:", err)
			errdata := message{
				Error:       "Sorry the page you are trying to view is no longer available... or your link have expired !!! " + fmt.Sprintf("%v", err),
				Information: "Go back to telegram to Generate New link. Links are use only once.",
			}

			// template from db
			templ := getTemplate("message") // select the template name from db
			t := template.New("error")
			t, _ = t.Parse(templ)
			//t, _ = t.Parse(templ)

			//template from file
			//t, _ := template.ParseFiles("static/message_page.html")
			t.Execute(w, errdata)
			return
		}
		//set the code
		data.State = stateDecoded
		cc := getJsonCountryCode(jsonCountryCode)
		types := getJsonAccountType(jsonAccountTypes)

		//set the country name of the retrieved country code
		data.CountryName = cc[data.CountryCode]
		//start choosing template
		//	fmt.Println("Registration Type:", data.RegType)
		if data.RegType == 0 {
			//	normal registration
			pgdata := signPageData{
				Title:        "Register A New Account",
				AccountInfo:  data,
				CountryCode:  cc,
				AccountTypes: types,
			}
			//fmt.Println(pgdata)
			// template from db
			templ := getTemplate("sign_up") // select the template name from db
			//fmt.Println("Gotten trade_profile_new Template From Db", templ)
			t := template.New("Sign_up")
			t, _ = t.Funcs(Fmap).Parse(templ)

			t.Execute(w, pgdata)
			return
		}
		if data.RegType == 1 {
			//managed account
			pgdata := signPageData{
				Title:        "Register A Managed Account",
				AccountInfo:  data,
				CountryCode:  cc,
				AccountTypes: types,
			}
			//fmt.Println(pgdata)
			// template from db
			templ := getTemplate("signupm") // select the template name from db
			//fmt.Println("Gotten trade_profile_new Template From Db", templ)
			t := template.New("Signupm")
			t, _ = t.Funcs(Fmap).Parse(templ)

			t.Execute(w, pgdata)
			return

		}
	}

	if r.Method == "POST" {
		statePost := r.FormValue("state")
		//	fmt.Printf("statePost", statePost)
		phone := r.FormValue("phone")
		accountType := r.FormValue("accountType")
		//fmt.Println("accountType", accountType)
		fname := r.FormValue("fname")
		//fmt.Println("fname", fname)
		lname := r.FormValue("lname")
		//fmt.Println("lname", lname)
		email := r.FormValue("email")
		//fmt.Println("email", email)
		facebookID := r.FormValue("facebook")
		//fmt.Println("facebookID", facebookID)
		whatsapp := r.FormValue("whatapp")
		//fmt.Println("whatsapp", whatsapp)
		countryCode := r.FormValue("countryCode")
		regType := r.FormValue("regType")
		//fmt.Println("countryCode", countryCode)
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
		var insertq, retstate string
		var err error
		if regType == "1" {
			insertq = `UPDATE account_signup SET phone = $1,account_type = $2,email=$3,facebook_id=$4,whatsapp=$5,
		country_code=$6,finalized=$7,finalized_on=$8,first_name=$10,last_name=$11, initiated_on=now()::timestamp WHERE state = $9 RETURNING state`

			err = con.Db.QueryRow(insertq, phone, accountType, email, facebookID, whatsapp, countryCode, finalized,
				finalizedOn, statePost, fname, lname).Scan(&retstate)
		} else {
			insertq = `UPDATE account_signup SET phone = $1,account_type = $2,email=$3,facebook_id=$4,whatsapp=$5,
			country_code=$6,finalized=$7,first_name=$9,last_name=$10, initiated_on=now()::timestamp WHERE state = $8 RETURNING state`

			err = con.Db.QueryRow(insertq, phone, accountType, email, facebookID, whatsapp, countryCode, 0, statePost, fname, lname).Scan(&retstate)
		}

		//	fmt.Println(insert)
		//if finalized is 0 then go back to the form to keep modification
		//	fmt.Println("ReCode = ", retcode, " Posted Code =", codePost)
		//	fmt.Println("Error is: ", err)
		if err == nil && (retstate == statePost) && finalized == 1 { /// insert was successful give user success message. and send message.
			url := fmt.Sprintf("%s?state=%v", h.GoogleCallBackURL, statePost)
			//fmt.Println("The reg Type Is..............", regType)
			if regType == "1" {
				datamsg := message{
					Sucess:      "Congratulation!\nYou have successfully Created New Managed Account for:\n" + lname + " " + fname + "!",
					Information: "You can Go back to telegram to continue your trading experience.",
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
			} else {
				// Other rgistration..... take the person to do google sign up
				http.Redirect(w, r, url, 301)
			}
		} else if err == nil && (retstate == statePost) && finalized == 0 {
			//Redirect user to new page to continue modification.
			//Do not complete until finalized

			selsql := `SELECT first_name,last_name,email,phone,account_type,allow_api_withdraw,facebook_id,whatsapp,
			country_code,json_country_code,json_trade_types FROM account_signup WHERE state = $1 LIMIT 1`
			err := con.Db.QueryRow(selsql, statePost).Scan(&data.FirstName, &data.LastName, &data.Email, &data.Phone,
				&data.AccountType, &data.AllowAPIWithdraw,
				&data.FacebookID, &data.Whatsapp, &data.CountryCode, &jsonCountryCode, &jsonAccountTypes)
			//	fmt.Println(err)
			if err != nil {
				fmt.Println("RouteSignUp profile selection with code failed:", err)
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

			cc := getJsonCountryCode(jsonCountryCode)
			types := getJsonAccountType(jsonAccountTypes)
			//set the country name of the retrieved country code
			data.CountryName = cc[data.CountryCode]
			//assign the submitted code to the new data
			//	fmt.Println("Registration Type:", data.RegType)
			data.State = statePost
			if regType == "0" {
				//normal signup
				//start choosing templates
				pgdata := signPageData{
					Title:        "New Account. Last Update Successful",
					AccountInfo:  data,
					CountryCode:  cc,
					AccountTypes: types,
				}
				//fmt.Println(pgdata)
				// template from db
				templ := getTemplate("sign_up") // select the template name from db
				//fmt.Println("Gotten trade_profile_new Template From Db", templ)
				t := template.New("Sign_up")
				t, _ = t.Funcs(Fmap).Parse(templ)

				t.Execute(w, pgdata)
				//end choosing templates

				return
			}
			if regType == "1" {
				//managed signup
				//start choosing templates
				pgdata := signPageData{
					Title:        "Managed Account. Last Update Successful",
					AccountInfo:  data,
					CountryCode:  cc,
					AccountTypes: types,
				}
				//fmt.Println(pgdata)
				// template from db
				templ := getTemplate("signupm") // select the template name from db
				//fmt.Println("Gotten trade_profile_new Template From Db", templ)
				t := template.New("Signupm")
				t, _ = t.Funcs(Fmap).Parse(templ)

				t.Execute(w, pgdata)
				//end choosing templates

				return
			}

		} else {
			fmt.Println("RouteSignUp. Update with code failed:", err)
			errdata := message{
				Error:       "Sorry your Account failed to create: " + fmt.Sprintf("%v", err),
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

func checkState(tableName, code string) bool {
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("tradeprofile.go:getTemplate:connection failed:", err)
		return false
	}
	var tpoid int64
	defer con.Close()
	sql := "SELECT count(*) FROM " + tableName + " as c WHERE state = $1 AND finalized = 0"
	err = con.Db.QueryRow(sql, code).Scan(&tpoid)
	//fmt.Println(err)
	if err != nil {
		return false
	}
	if tpoid > 0 {
		return true
	}
	return false
}
