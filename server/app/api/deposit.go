package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"twistserver/app/datastore"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type CreateDepositRequest struct {
	User    string `json:"user"`
	Balance int    `json:"balance"`
}

func InitDepositAPI(r *gin.Engine) {

	// Try
	r.POST("/api/v1/deposit", func(c *gin.Context) {
		switch c.Request.Header.Get("twist-phrase") {
		case "confirm":
			if c.Request.Header.Get("twist-task-id") == "" {
				c.JSON(http.StatusOK, gin.H{
					"error": "Need task ID",
				})
			}

			task := GetTask(c.Request.Header.Get("twist-task-id"))

			// JSON FORM Read
			var taskJSON map[string]interface{}
			json.Unmarshal([]byte(task), &taskJSON)
			taskStateString := taskJSON["payload"].(string)

			var taskStateJSON map[string]interface{}
			json.Unmarshal([]byte(taskStateString), &taskStateJSON)
			// Execute to update database

			datastore.DataBalance[taskStateJSON["user"].(string)] += int(taskStateJSON["balance"].(float64))
			c.JSON(http.StatusOK, gin.H{
				"user":   taskStateJSON["user"].(string),
				"wallet": strconv.Itoa(datastore.DataBalance[taskStateJSON["user"].(string)]),
			})

		case "cancel":
			if c.Request.Header.Get("twist-task-id") == "" {
				c.JSON(http.StatusOK, gin.H{
					"error": "Need task ID",
				})
			}

			task := GetTask(c.Request.Header.Get("twist-task-id"))

			// JSON FORM Read
			var taskJSON map[string]interface{}
			json.Unmarshal([]byte(task), &taskJSON)
			taskStateString := taskJSON["payload"].(string)

			var taskStateJSON map[string]interface{}
			json.Unmarshal([]byte(taskStateString), &taskStateJSON)

			// rollback if confirmed already
			// Task need to be JSON
			if taskStateJSON["status"] == "CONFIRMED" {
				datastore.DataBalance[taskStateJSON["user"].(string)] -= int(taskStateJSON["balance"].(float64))
			}
			c.JSON(http.StatusOK, gin.H{})
		default:
			var request CreateDepositRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if request.User == "" {
				c.JSON(http.StatusOK, gin.H{
					"error": "Need User",
				})
			}
			if request.Balance < 0 {
				c.JSON(http.StatusOK, gin.H{
					"error": "Require balance",
				})
			}
			if datastore.DataUser[request.User] == false {
				c.JSON(http.StatusOK, gin.H{
					"error": "User Not Alive",
				})
			}
			taskResp := CreateTask(`{"task":{"actions":{"confirm":{"type":"rest","method":"post","uri":"` + viper.GetString("host.serviceHost") + `/api/v1/deposit"},"cancel":{"type":"rest","method":"post","uri":"` + viper.GetString("host.serviceHost") + `/api/v1/deposit"}},"payload":"{\"user\":\"armani\",\"balance\":100}","timeout":30000}}`)
			c.Data(http.StatusOK, "application/json", []byte(taskResp))
		}
	})
}
