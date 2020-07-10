package twist

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber"
)

func DeductTry(ctx *fiber.Ctx) {
	fmt.Println("===Deduct-Try===")
	taskState := new(Task)
	if err := ctx.BodyParser(taskState); err != nil {
		log.Fatal(err)
	}
	if taskState.User == "" {
		log.Fatal("Required user")
	}

	if taskState.Balance == 0 {
		log.Fatal("Required balance")
	}
	// Check user whether does exists or not
	user := DataUser[taskState.User]
	if user == false {
		log.Fatal("User not Alive")
	}

	// Check balances
	if DataBalance[taskState.User]-DataReserve[taskState.User] < taskState.Balance {
		log.Fatal("Balance in wallet is not engough")
	}

	// Do reserved
	DataReserve[taskState.User] += taskState.Balance

	// Twist lifecycle
	taskResponse := CreateTask(`{"task":{"actions":{"confirm":{"type":"rest","method":"put","uri":"` + serviceHost + `/deduct"},"cancel":{"type":"rest","method":"delete","uri":"` + serviceHost + `/deduct"}},"payload":"{\"user\":\"` + taskState.User + `\",\"balance\":` + strconv.Itoa(taskState.Balance) + `}","timeout":30000}}`)

	// Response
	// how to use fiber write body

	ctx.SendString(taskResponse)
}

func DeductConfirm(ctx *fiber.Ctx) {
	fmt.Println("===Deduct-Confirm===")

	if ctx.Get("twist-task-id") == "" {
		log.Println("Required task ID")
	}
	// Getting task state
	task := GetTask(ctx.Get("twist-task-id"))

	// JSON FORM Read
	var taskJSON map[string]interface{}
	json.Unmarshal([]byte(task), &taskJSON)
	taskStateString := taskJSON["payload"].(string)

	var taskStateJSON map[string]interface{}
	json.Unmarshal([]byte(taskStateString), &taskStateJSON)

	// Execute to update database

	user := taskStateJSON["user"].(string)
	DataBalance[user] -= int(taskStateJSON["balance"].(float64))
	DataReserve[user] -= int(taskStateJSON["balance"].(float64))

	// Response
	// how to use fiber write body
	fmt.Println("Need")
	ctx.SendString(`{"user":"` + taskStateJSON["user"].(string) + `","wallet":"` + strconv.Itoa(int(taskStateJSON["balance"].(float64))) + `"}`)
}

func DeductCancel(ctx *fiber.Ctx) {
	fmt.Println("===Deduct-Cancel===")
	if ctx.Get("twist-task-id") == "" {
		log.Println("Required task ID")
	}

	// Getting task state
	task := GetTask(ctx.Get("twist-task-id"))

	// JSON FORM Read
	var taskJSON map[string]interface{}
	json.Unmarshal([]byte(task), &taskJSON)
	taskStateString := taskJSON["payload"].(string)

	var taskStateJSON map[string]interface{}
	json.Unmarshal([]byte(taskStateString), &taskStateJSON)

	if taskStateJSON["status"] == "CONFIRMED" {
		// Rollback if confirmed already
		DataBalance[taskStateJSON["user"].(string)] += int(taskStateJSON["balance"].(float64))
	} else {
		// Release reserved resources
		DataReserve[taskStateJSON["user"].(string)] -= int(taskStateJSON["balance"].(float64))
	}

	ctx.SendString("")

}
