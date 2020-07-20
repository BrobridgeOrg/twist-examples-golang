package api

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func CreateTask(task string) string {

	client := &http.Client{}
	var jsonStr = []byte(task)
	req, err := http.NewRequest("POST", viper.GetString("host.twistHost")+"/api/v1/tasks", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

func GetTask(taskID string) string {

	client := &http.Client{}
	req, err := http.NewRequest("GET", viper.GetString("host.twistHost")+"/api/v1/tasks/"+taskID, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

func CancelTask(taskID string) string {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", viper.GetString("host.twistHost")+"/api/v1/tasks/"+taskID, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
