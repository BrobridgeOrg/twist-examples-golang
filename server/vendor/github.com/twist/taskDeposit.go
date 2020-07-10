package twist

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber"
)

func Deposit(ctx *fiber.Ctx) {
	fmt.Println(ctx.Get("twist-phrase"))
	switch ctx.Get("twist-phrase") {
	case "confirm":
		fmt.Println("===Deposit-Confirm===")
		if ctx.Get("twist-task-id") == "" {
			log.Fatal("Required task ID")
		}
		// Getting task state
		task := GetTask(ctx.Get("twist-task-id"))
		fmt.Println(task)
		// JSON FORM Read
		var taskJSON map[string]interface{}
		json.Unmarshal([]byte(task), &taskJSON)
		taskStateString := taskJSON["payload"].(string)

		var taskStateJSON map[string]interface{}
		json.Unmarshal([]byte(taskStateString), &taskStateJSON)
		// Execute to update database

		DataBalance[taskStateJSON["user"].(string)] += int(taskStateJSON["balance"].(float64))
		// Response Fiber!!
		ctx.SendString(`{"user":"` + taskStateJSON["user"].(string) + `","wallet":"` + strconv.Itoa(DataBalance[taskStateJSON["user"].(string)]) + `"}`)

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
		taskStateString := taskJSON["payload"].(string)

		var taskStateJSON map[string]interface{}
		json.Unmarshal([]byte(taskStateString), &taskStateJSON)

		// rollback if confirmed already
		// Task need to be JSON
		if taskStateJSON["status"] == "CONFIRMED" {
			DataBalance[taskStateJSON["user"].(string)] -= int(taskStateJSON["balance"].(float64))
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
		taskResp := CreateTask(`{"task":{"actions":{"confirm":{"type":"rest","method":"post","uri":"` + serviceHost + `/deposit"},"cancel":{"type":"rest","method":"post","uri":"` + serviceHost + `/deposit"}},"payload":"{\"user\":\"armani\",\"balance\":100}","timeout":30000}}`)
		fmt.Println(taskResp)
		ctx.SendString(taskResp)

	}
}
