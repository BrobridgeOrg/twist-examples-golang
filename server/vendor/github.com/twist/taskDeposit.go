package twist

import (
	"fmt"
	"log"

	"encoding/json"

	"github.com/gofiber/fiber"
)

func Deposit(ctx *fiber.Ctx) {

	switch ctx.Get("twist-phrase") {
	case "confirm":
		fmt.Println("===Deposit-Confirm===")
		if ctx.Get("twist-task-id") == "" {
			log.Fatal("Required task ID")
		}
		// Getting task state
		task := GetTask(ctx.Get("twist-task-id"))

		// JSON FORM Read
		var taskJSON map[string]interface{}
		json.Unmarshal([]byte(task), &taskJSON)
		taskState := taskJSON["payload"].(map[string]interface{})
		// Execute to update database

		DataBalance[taskState["user"].(string)] += taskState["balance"].(int)

		// Response Fiber!!
		ctx.SendString("{user:" + taskState["user"].(string) + ",wallet:" + string(DataBalance["User"]) + "}")

	case "cancel":
		fmt.Println("===Deposit-Cancel===")

		if ctx.Get("twist-task-id") == "" {
			log.Println("Required task ID")
		}

		// Getting task state
		task := GetTask(ctx.Get("twist-task-id"))

		// JSON FORM Read
		var taskJSON map[string]interface{}
		json.Unmarshal([]byte(task), &taskJSON)
		taskState := taskJSON["payload"].(map[string]interface{})

		// rollback if confirmed already
		// Task need to be JSON
		if taskState["status"] == "CONFIRMED" {
			DataBalance[taskState["user"].(string)] -= taskState["balance"].(int)
		}

		// Response clear

		ctx.SendString("")

	default:
		fmt.Println("===Deposit-Try===")

		taskState := new(Task)
		if err := ctx.BodyParser(taskState); err != nil {
			log.Fatal(err)
		}

		if taskState.User == "" {
			log.Fatal("Required user")
		}

		// ?? if state is confuse ??
		if taskState.Balance == 0 {
			log.Fatal("Required balance")
		}

		// Prepare a task state
		newTask := Task{
			User:    taskState.User,
			Balance: taskState.Balance}

		// Check user whether does exists or not
		user := DataUser[newTask.User]
		if !user {
			log.Fatal("User not Alive")
		}

		// Twist lifecycle
		taskResp := CreateTask("{actions:{confirm:{type: 'rest',method: 'post',uri: 'http://127.0.0.1:3000/deposit'},cancel:{type: 'rest',method: 'post',uri: 'http://127.0.0.1:3000/deposit'}},payload:" + "{User:" + string(taskState.User) + ",Balance:" + string(taskState.Balance) + "}" + ",timeout: 30000}")
		ctx.SendString(taskResp)

	}
}
