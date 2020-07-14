package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"twistserver/app/datastore"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type CreateDeductRequest struct {
	User    string `json:"user"`
	Balance int    `json:"balance"`
}

func InitDeductAPI(r *gin.Engine) {

	// Try
	r.POST("/api/v1/deduct", func(c *gin.Context) {

		var request CreateDeductRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if request.User == "" {
			log.Fatal("Need user")
		}
		if request.Balance == 0 {
			log.Fatal("Require balance")
		}
		if datastore.DataUser[request.User] == false {
			log.Fatal("User is not alive")
		}

		if datastore.DataBalance[request.User]-datastore.DataReserve[request.User] < request.Balance {
			log.Fatal("Balance in wallet is not engough")
		}

		// Do reserved
		datastore.DataReserve[request.User] += request.Balance

		// Repair
		jsonData := CreateTask(`{"task":{"actions":{"confirm":{"type":"rest","method":"put","uri":"` + viper.GetString("host.serviceHost") + `/deduct"},"cancel":{"type":"rest","method":"delete","uri":"` + viper.GetString("host.serviceHost") + `/deduct"}},"payload":"{\"user\":\"` + request.User + `\",\"balance\":` + strconv.Itoa(request.Balance) + `}","timeout":30000}}`)

		c.Data(http.StatusOK, "application/json", []byte(jsonData))
	})

	// Confirm
	r.PUT("/api/v1/deduct", func(c *gin.Context) {

		if c.Request.Header.Get("twist-task-id") == "" {
			log.Fatal("Need task ID")
		}
		task := GetTask(c.Request.Header.Get("twist-task-id"))

		// JSON FORM Read
		var taskJSON map[string]interface{}
		json.Unmarshal([]byte(task), &taskJSON)
		taskStateString := taskJSON["payload"].(string)

		var taskStateJSON map[string]interface{}
		json.Unmarshal([]byte(taskStateString), &taskStateJSON)

		// Execute to update database

		user := taskStateJSON["user"].(string)
		datastore.DataBalance[user] -= int(taskStateJSON["balance"].(float64))
		datastore.DataReserve[user] -= int(taskStateJSON["balance"].(float64))
		c.JSON(http.StatusOK, gin.H{
			"user":   taskStateJSON["user"].(string),
			"wallet": strconv.Itoa(int(taskStateJSON["balance"].(float64))),
		})
	})

	// Cancel
	r.DELETE("/api/v1/deduct", func(c *gin.Context) {
		if c.Request.Header.Get("twist-task-id") == "" {

		}

		task := GetTask(c.Request.Header.Get("twist-task-id"))

		// JSON FORM Read
		var taskJSON map[string]interface{}
		json.Unmarshal([]byte(task), &taskJSON)
		taskStateString := taskJSON["payload"].(string)

		var taskStateJSON map[string]interface{}
		json.Unmarshal([]byte(taskStateString), &taskStateJSON)

		if taskStateJSON["status"] == "CONFIRMED" {
			// Rollback if confirmed already
			datastore.DataBalance[taskStateJSON["user"].(string)] += int(taskStateJSON["balance"].(float64))
		} else {
			// Release reserved resources
			datastore.DataReserve[taskStateJSON["user"].(string)] -= int(taskStateJSON["balance"].(float64))
		}

		c.JSON(http.StatusOK, gin.H{})
	})

}
