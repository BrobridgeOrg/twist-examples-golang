package twist

import (
	"encoding/json"
	"fmt"
	"log"

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
	taskResponse := CreateTask("{actions:{confirm:{type: 'rest',method: 'put',uri: 'http://172.20.10.10:3000/deduct'},cancel:{type: 'rest',method: 'delete',uri: 'http://172.20.10.10:3000/deduct'}},payload:{User:" + taskState.User + ",Balance:" + string(taskState.Balance) + "},timeout: 30000,}")

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
	taskState := taskJSON["payload"].(map[string]interface{})

	// Execute to update database
	user := taskState["User"].(string)
	DataBalance[user] -= taskState["balance"].(int)
	DataReserve[user] -= taskState["balance"].(int)

	// Response
	// how to use fiber write body

	ctx.SendString("{user:" + taskState["User"].(string) + ",wallet:" + taskState["Balance"].(string) + "}")
}

func DeductCancel(ctx *fiber.Ctx) {
	fmt.Println("===Deduct-Cancel===")
	fmt.Println("===DeductCancel===")
	if ctx.Get("twist-task-id") == "" {
		log.Println("Required task ID")
	}

	// Getting task state
	task := GetTask(ctx.Get("twist-task-id"))

	// JSON FORM Read
	var taskJSON map[string]interface{}
	json.Unmarshal([]byte(task), &taskJSON)
	taskState := taskJSON["payload"].(map[string]interface{})

	if taskState["status"] == "CONFIRMED" {
		// Rollback if confirmed already
		DataBalance[taskState["user"].(string)] += taskState["Balance"].(int)
	} else {
		// Release reserved resources
		DataReserve[taskState["user"].(string)] -= taskState["Balance"].(int)
	}

	ctx.SendString("")

}
