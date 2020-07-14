package api

import (
	"net/http"
	"strconv"
	"twistserver/app/datastore"

	"github.com/gin-gonic/gin"
)

func InitWalletsAPI(r *gin.Engine) {
	r.GET("/api/v1/wallets", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"wallets": gin.H{
				"fred": gin.H{
					"balance":  strconv.Itoa(datastore.DataBalance["fred"]),
					"reserved": strconv.Itoa(datastore.DataReserve["fred"]),
				},
				"armani": gin.H{
					"balance":  strconv.Itoa(datastore.DataBalance["armani"]),
					"reserved": strconv.Itoa(datastore.DataReserve["armani"]),
				},
			},
		})
	})
}
