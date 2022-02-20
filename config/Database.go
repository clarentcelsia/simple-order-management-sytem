package config

import (
	"database/sql"
	"net/http"

	MResponse "restaurant/model/response"

	logger "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type DBJSON struct {
	DBStats sql.DBStats
	DBName  string
}

func GetDBHealthCheck(c *gin.Context) {
	//Create new list of DBJson using make
	newStatArr := make([]DBJSON, 0)
	dbconn, err := GetConnection()
	if err != nil {
		//Serialize struct as JSON
		c.JSON(http.StatusOK, MResponse.DBHealthResponse{
			DBStatus: http.StatusInternalServerError,
			Message:  "Database is not running",
		})
		logger.Info("PostgreDB is not running")
		return
	}

	newStatArr = append(newStatArr, DBJSON{
		DBStats: dbconn.Stats(),
		DBName:  "dbrestaurant",
	})

	/*
		if there're more than 1 database
		newStatArr = append(newStatArr, dbJson2)
	*/

	c.JSON(http.StatusOK, MResponse.PageResponse{
		Status:  http.StatusOK,
		Message: "Connection Found",
		Items:   newStatArr,
	})

}
