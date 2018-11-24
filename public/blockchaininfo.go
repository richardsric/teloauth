package public

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	humanize "github.com/dustin/go-humanize"
)

var keyCode = "13ba4210-f827-499d-9f80-3dfc9a0e702f"
var btcAddress = "1HkRzK3hsAQUkghd5wqWpqtRpCs9tWz4n5"
var callBackURL = "http://itradecoin.ngrok.io/blockchaincallback?secret="
var reqURL = "https://api.blockchain.info/v2/receive/balance_update"

// Visit https://blockchain.info/api/api_receive

//BlockChainUpdateCallBack handles callbacks from blockchain address monitoring
func BlockChainUpdateCallBack(w http.ResponseWriter, r *http.Request) {

	dbsecret := "webdb"

	//get transaction_hash
	transactionHash := r.URL.Query().Get("transaction_hash")

	//get address
	address := r.URL.Query().Get("address")

	//get confirmations
	confirmations := r.URL.Query().Get("confirmations")

	//get value in satoshi and convert to BTC
	gottenValueStr := r.URL.Query().Get("value")
	if gottenValueStr == "" {
		gottenValueStr = "0"
	}

	val, err := strconv.ParseFloat(gottenValueStr, 64)

	if err != nil {
		fmt.Println("gotten value conversion error due to ", err)
	}
	//val is in satoshi. Convert to BTC
	valueInBTC := val / 100000000

	//get address
	callSecret := r.URL.Query().Get("secret")
	itx := r.URL.Query().Get("itx")
	inv := r.URL.Query().Get("invoice")

	//get other params we passed to it
	msg1 := fmt.Sprint("Notification Received:\n", "transaction_hash: ", transactionHash, "\naddress: ", address, "\nconfirmations: ", confirmations, "\nvalue: Ƀ ", humanize.FormatFloat("#,###.########", valueInBTC), "\nTxID: ", itx, "\nInvoiceNumber: ", inv)

	if dbsecret == callSecret {
		//send telegram message
		if valueInBTC > 0 {
			msg := "<b>BTC Payment Alert!</b>\n" + msg1
			SendServiceStatusIM(msg)
		}
		//callback is valid
		fmt.Println("Valid Notification Received:\n", "transaction_hash: ", transactionHash, "\naddress: ", address, "\nconfirmations: ", confirmations, "\nvalue: Ƀ ", humanize.FormatFloat("#,###.########", valueInBTC), "\nsecret: ", callSecret)
		fmt.Fprint(w, "*ok*")
		return
	}
	//send telegram message
	if valueInBTC > 0 {
		msg := "<b>BTC Payment Alert!</b>\n" + msg1
		SendServiceStatusIM(msg)
	}
	//respond with this information
	fmt.Fprint(w, "*ok*")
	return
}

//DeleteMonitoredAddress removes a monitored address from being monitored
func DeleteMonitoredAddress(id int) bool {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%d?key=%s", reqURL, id, keyCode), nil)
	if err != nil {
		// handle err
		fmt.Println("Delete Address Monitor Request", err)
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		fmt.Println("Delete Address Monitor Request DO Request", err)
		return false
	}
	defer resp.Body.Close()

	resbody, _ := ioutil.ReadAll(resp.Body)
	var m map[string]interface{}

	err = json.Unmarshal(resbody, &m)
	//fmt.Println("response: \n", string(resbody))
	if err != nil {
		fmt.Println("Delete Address Monitor JSON fail:", err)
		return false
	}
	if m["deleted"] != nil {
		del := m["deleted"].(bool)
		return del
	}
	return false
}
