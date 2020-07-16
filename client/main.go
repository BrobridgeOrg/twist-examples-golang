package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func main() {

	configHost()
	fmt.Println(viper.GetString("host.serviceHost"))
	fmt.Println(viper.GetString("host.twistHost"))

	fmt.Println("===========ORIGINAL===========")
	printWalletInfo()
	fmt.Println("*** CREATE TRANSACTION ***")
	TxID := createTransaction()

	// TRY
	fmt.Println("*** TRY ***")
	deductJSON := deduct(TxID)
	registerTasks(TxID, deductJSON)

	depositJSON := deposit(TxID)
	registerTasks(TxID, depositJSON)

	// CONFIRM

	doConfirm(TxID)
	fmt.Println("==========RESULTS==========")
	printWalletInfo()
	fmt.Println("*** Transaction was successfully processed ***")
}

type txID struct {
	TransactionID string
}

type errorState struct {
	Error string
}

type registerState struct {
	Reason  string
	Success string
}

func configHost() {
	// From the environment
	viper.SetEnvPrefix("TWIST_Example")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// From config file
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Warn("No configuration file was loaded")
	}
}
func printWalletInfo() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", viper.GetString("host.serviceHost")+"/api/v1/wallets", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	if err != nil {
		log.Fatal(err)
		return
	}
}

func createTransaction() string {
	client := &http.Client{}
	var jsonStr = []byte(`{"timeout":3000}`)
	req, err := http.NewRequest("POST", viper.GetString("host.twistHost")+"/api/v1/transactions", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var id txID
	json.Unmarshal([]byte(string(body)), &id)
	fmt.Println("Create tx with id:" + id.TransactionID)
	if err != nil {
		log.Fatal(err)
	}
	return id.TransactionID
}

func deduct(txid string) string {
	fmt.Println("[TRY]deduct 100 from fred")
	client := &http.Client{}
	var jsonStr = []byte(`{"user":"fred","balance":100}`)
	req, err := http.NewRequest("POST", viper.GetString("host.serviceHost")+"/api/v1/deduct", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Twist-Transaction-ID", txid)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println(string(body))

	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func deposit(txid string) string {
	fmt.Println("[TRY]deposit 100 to armani's wallet")
	client := &http.Client{}
	var jsonStr = []byte(`{"user":"armani","balance":100}`)
	req, err := http.NewRequest("POST", viper.GetString("host.serviceHost")+"/api/v1/deposit", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Twist-Transaction-ID", txid)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println(string(body))

	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func registerTasks(txid string, taskJSON string) {
	fmt.Println("[RegisterTasks]")
	client := &http.Client{}
	var jsonStr = []byte(`{"tasks":[` + taskJSON + `]}`)

	req, err := http.NewRequest("PUT", viper.GetString("host.twistHost")+"/api/v1/transactions/"+txid, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	var rgstState registerState
	json.Unmarshal([]byte(string(body)), &rgstState)
	if rgstState.Success == "false" {
		fmt.Println("Register False")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func doConfirm(txid string) {
	fmt.Println("*** CONFIRM ***")
	fmt.Println("Entering CONFIRM phase for:" + txid)

	client := &http.Client{}
	var jsonStr = []byte(``)
	req, err := http.NewRequest("POST", viper.GetString("host.twistHost")+"/api/v1/transactions/"+txid, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var errJSON errorState
	json.Unmarshal([]byte(string(body)), &errJSON)
	if errJSON.Error == "EOF" {
		fmt.Println("Failed to confirm")
		doCancel(txid)
	}

	if err != nil {
		log.Fatal(err)
	}

}

func doCancel(txid string) {
	fmt.Println("*** CANCEL ***")
	fmt.Println("Entering Cancel phase for:" + txid)

	client := &http.Client{}
	var jsonStr = []byte("")
	req, err := http.NewRequest("DELETE", viper.GetString("host.twistHost")+"/api/v1/transactions/"+txid, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	fmt.Println("Transaction was canceled successfully")

	if err != nil {
		log.Fatal(err)
	}

}
