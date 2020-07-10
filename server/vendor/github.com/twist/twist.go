package twist

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber"
)

var twistHost string = "http://192.168.1.174:32529"
var serviceHost string = "http://192.168.1.150:3000"

func Wallets(c *fiber.Ctx) {
	fmt.Println("===Wallets===")
	c.JSON(`
	wallets:{
		fred:{
			balance:` + strconv.Itoa(DataBalance["fred"]) + `,
			reserved:` + strconv.Itoa(DataReserve["fred"]) + `
		},
		armani:{
			balance:` + strconv.Itoa(DataBalance["armani"]) + `,
			reserved:` + strconv.Itoa(DataReserve["armani"]) + `
		}
	}`)
}

type Task struct {
	User    string `json:"user"`
	Balance int    `json:"balance"`
}
type TaskState struct {
	User    string `json:"user"`
	Balance int    `json:"balance"`
}

type Payload struct {
	User    string `json:"user"`
	Balance int    `json:"balance"`
}

func CreateAccount(name string, initBalance int) {
	DataBalance[name] = initBalance
	DataReserve[name] = 0
	DataUser[name] = true
}

var DataBalance = make(map[string]int)
var DataReserve = make(map[string]int)
var DataUser = make(map[string]bool)

//
func CreateTask(task string) string {
	fmt.Println("==CreateTask==")
	//fmt.Println(task)
	client := &http.Client{}
	var jsonStr = []byte(task)
	req, err := http.NewRequest("POST", twistHost+"/api/v1/tasks", bytes.NewBuffer(jsonStr))
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
	fmt.Println("------------CreateTask-Resp------------")
	return string(body)
}

func GetTask(taskID string) string {
	fmt.Println("==GetTask==")
	fmt.Println(taskID)
	client := &http.Client{}

	req, err := http.NewRequest("GET", twistHost+"/api/v1/tasks/"+taskID, nil)
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
	fmt.Println("------------GetTask-Resp------------")
	return string(body)
}

func CancelTask(taskID string) string {
	fmt.Println("==CancelTask==")
	fmt.Println(taskID)
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", twistHost+"/api/v1/tasks/"+taskID, nil)
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
	fmt.Println("------------CancelTask-Resp------------")
	return string(body)
}
