package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var twistHost string = "http://192.168.1.174:32529"
var serviceHost string = "http://192.168.1.151:3000"

func main() {

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

	// Error Need but not write
	// doCancel(transactionID)

	// CONFIRM

	doConfirm(TxID)
	fmt.Println("==========RESULTS==========")
	printWalletInfo()
	fmt.Println("*** Transaction was successfully processed ***")
}

func printWalletInfo() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", serviceHost+"/api/v1/wallets", nil)
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
	req, err := http.NewRequest("POST", twistHost+"/api/v1/transactions", bytes.NewBuffer(jsonStr))
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

type txID struct {
	TransactionID string
}

func deduct(txid string) string {
	fmt.Println("[TRY]deduct 100 from fred")
	client := &http.Client{}
	var jsonStr = []byte(`{"user":"fred","balance":100}`)
	req, err := http.NewRequest("POST", serviceHost+"/api/v1/deduct", bytes.NewBuffer(jsonStr))
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
	req, err := http.NewRequest("POST", serviceHost+"/api/v1/deposit", bytes.NewBuffer(jsonStr))
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

	req, err := http.NewRequest("PUT", twistHost+"/api/v1/transactions/"+txid, bytes.NewBuffer(jsonStr))
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

	if err != nil {
		log.Fatal(err)
	}
}

func doConfirm(txid string) {
	fmt.Println("*** CONFIRM ***")
	fmt.Println("Entering CONFIRM phase for:" + txid)

	client := &http.Client{}
	var jsonStr = []byte(``)
	req, err := http.NewRequest("POST", twistHost+"/api/v1/transactions/"+txid, bytes.NewBuffer(jsonStr))
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

	fmt.Println(string(body))

	if err != nil {
		log.Fatal(err)
	}

}
