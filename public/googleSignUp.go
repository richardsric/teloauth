package public

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	h "github.com/richardsric/teloauth/helper"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "",
		ClientID:     "467979088936-k2rjas1p33j1ui907te9p5enl411vmsq.apps.googleusercontent.com",
		ClientSecret: "TfiQIjKJmxv2TD9k9fJUTq7l",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
)

//HandleGoogleSignUp deals with google signup
func HandleGoogleSignUp(w http.ResponseWriter, r *http.Request) {
	state := strings.Trim(r.URL.Query().Get("state"), " ") //r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if code == "" && state != "" {
		result, err := h.DBSelect("SELECT state FROM account_signup WHERE state=$1", state)
		if err != nil {
			fmt.Println("incorrect link, state doesn't exist")
			fmt.Fprintln(w, invalidResp)
			fmt.Println("select error due to ", err)
			return
		}
		if result == nil {
			fmt.Fprintln(w, "state doesn't exist in the db")
			return
		}
		//	fmt.Println("state selected is ", result.(string))
		check := h.DBSelectRow("SELECT authorized,reg_status FROM account_signup WHERE state=$1", state)
		if check.ErrorMsg != "" {
			fmt.Println("select error due to ", check.ErrorMsg)
			return
		}
		if int(check.Columns["authorized"].(int64)) == 0 && int(check.Columns["reg_status"].(int64)) == 1 {

			url := fmt.Sprintf("%v/SignUpAlready?state=%v", h.TelegramURL, state)
			http.Get(url)
			fmt.Fprintln(w, alreadyAuthResp)
			return
		} else if int(check.Columns["authorized"].(int64)) == 1 && int(check.Columns["reg_status"].(int64)) == 1 {

			url := fmt.Sprintf("%v/SignUpAuthorized?state=%v", h.TelegramURL, state)
			http.Get(url)
			fmt.Fprintln(w, authorisedResp)
			return
		}

		gs1 := h.DBSelectRow("SELECT signup_url FROM callback_url LIMIT 1")
		if gs1.ErrorMsg != "" {
			fmt.Println("Signup URL fetch error:", gs1.ErrorMsg)
			os.Exit(2)
		}
		rURL := gs1.Columns["signup_url"].(string)
		googleOauthConfig.RedirectURL = rURL
		url := googleOauthConfig.AuthCodeURL(state)
		//	fmt.Println("Google Redirect URL=", h.GoogleCallBackURL)
		//	fmt.Printf("Google OauthConfig = %+v", googleOauthConfig)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	} else {
		state2 := r.URL.Query().Get("state")
		//	fmt.Println("state from call back is ", state2)
		if state2 != state {
			fmt.Fprintln(w, errorResp)
			fmt.Println("no state match; possible csrf OR cookies not enabled")
			return
		}
		//code2 := r.URL.Query().Get("code") //code := r.FormValue("code")
		//fmt.Println("code from call back is ", code)

		token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println("there was an issue getting your token due to ", err)
			fmt.Fprintln(w, errorResp)
			return
		}
		//	fmt.Println("token is ", token)
		//checks if token is valid
		if !token.Valid() {
			fmt.Fprintln(w, "retreived invalid token")
			return
		}
		//using the gotten access token to query googleapis for userinfo
		//	fmt.Printf("access token - %v\nExpiry - %v\n", token.AccessToken, token.Expiry)
		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			//fmt.Fprintf(w,"Error occured pls try again")
			fmt.Printf("Call To Google api didn't work with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		//	fmt.Printf("Call To Google Api to get user info succeded and returned gmail signed in details as '%s'\n", contents)

		var u user
		err = json.Unmarshal(contents, &u)
		if err != nil {
			fmt.Fprintln(w, "unmarshalling user details failed due to :", err)
			return
		}
		if len(u.GivenName) > 2 && len(u.FamilyName) > 2 && len(u.Email) > 2 {
			res := h.DBModify("UPDATE account_signup SET first_name = $1,last_name = $2,email = $3,email_id = $4,reg_status=$5, finalized = 1, finalized_on=now()::timestamp WHERE state = $6",
				u.GivenName, u.FamilyName, u.Email, u.Id, 1, state)
			if res.AffectedRows < 0 || res.ErrorMsg != "" {
				fmt.Fprintln(w, "db insertion of user details failed due to", res.ErrorMsg)
				return
			}
		}
		fmt.Fprintln(w, firstAuthResp)

		url := fmt.Sprintf("%v/SignUpSuccess?state=%v", h.TelegramURL, state)
		http.Get(url)

	}
}
