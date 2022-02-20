package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	config "restaurant/config"
	cMenu "restaurant/controller/menu"
	mOrder "restaurant/model/order"
	util "restaurant/util"

	log "github.com/sirupsen/logrus"
)

func CreateOrder(c *gin.Context) {
	db, errdb := config.GetConnection()
	if errdb != nil {
		log.Info("Connection not found")
		log.Fatal(errdb)
	}

	// var order mOrder.Order
	// var orderDetail mOrder.OrderDetail
	var orderJson mOrder.OrderJSON
	errbind := c.BindJSON(&orderJson)

	if errbind != nil {
		log.Error(errbind)
	} else {
		stmtOrder := `INSERT INTO orders(order_id, customer_id, order_date) VALUES ($1, $2, $3) RETURNING *`
		stmtOrderDetail := `INSERT INTO orderdetail(order_detail_id, order_id, menu_id, current_price, qty, subtotal) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
		orderId := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)
		errOrder := db.QueryRow(stmtOrder, orderId, orderJson.CustomerID, time.Now().Format(time.RFC3339Nano)).Scan(&orderJson.OrderID, &orderJson.CustomerID, &orderJson.OrderDate)
		if errOrder != nil {
			log.Error(errOrder)
		}

		for i := 0; i < len(orderJson.OrderDetails); i++ {
			orderDetailId := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)
			menu, errreq := cMenu.GetMenuByID(c, db, orderJson.OrderDetails[i].MenuID)
			if errreq != nil {
				log.Error(errreq)
				util.HandleError(c, errreq, util.GET_REQUEST_FAILED, "Can't find the given id, please recheck the id", http.StatusNotFound)
			} else {
				subtotal := menu.MenuPrice * orderJson.OrderDetails[i].Qty
				errOrderDetail := db.QueryRow(stmtOrderDetail, orderDetailId, orderJson.OrderID, orderJson.OrderDetails[i].MenuID, menu.MenuPrice, orderJson.OrderDetails[i].Qty, subtotal).
					Scan(&orderJson.OrderDetails[i].OrderDetailID, &orderJson.OrderDetails[i].OrderID, &orderJson.OrderDetails[i].MenuID, &orderJson.OrderDetails[i].CurrentMenuPrice, &orderJson.OrderDetails[i].Qty, &orderJson.OrderDetails[i].Subtotal)
				if errOrderDetail != nil {
					log.Error(errOrderDetail)
				}
			}
		}

		OrderParam := map[string]interface{}{
			"order_id":      orderJson.OrderID,
			"customer_id":   orderJson.CustomerID,
			"order_date":    orderJson.OrderDate,
			"order_details": orderJson.OrderDetails,
		}

		util.HandleSuccess(c, OrderParam, util.INSERT_RES_SUCCEESS, http.StatusCreated)

	}

}

func GetOrders(c *gin.Context) {
	var orders []mOrder.OrderJSON
	db, errdb := config.GetConnection()
	if errdb != nil {
		log.Error("Connection not found")
		log.Fatal(errdb)
	}
	query := `SELECT * FROM orders`
	queryDetailOrder := `SELECT * FROM orderdetail WHERE order_id = $1`
	rows, errrows := db.Query(query)
	if errrows != nil {
		log.Error("Data not found")
		log.Fatal(errrows)
	}

	for rows.Next() {
		var order mOrder.OrderJSON
		if errscan := rows.Scan(&order.OrderID, &order.CustomerID, &order.OrderDate); errscan != nil {
			log.Error("Error while scanning data")
			log.Fatal(errscan)
		}

		orderDetailRows, errrows2 := db.Query(queryDetailOrder, order.OrderID)
		if errrows2 != nil {
			log.Error("Data not found")
			log.Fatal(errrows)
		}

		for orderDetailRows.Next() {
			var orderDetail mOrder.OrderDetail
			if errscandetail := orderDetailRows.Scan(
				&orderDetail.OrderDetailID,
				&orderDetail.OrderID,
				&orderDetail.MenuID,
				&orderDetail.CurrentMenuPrice,
				&orderDetail.Qty,
				&orderDetail.Subtotal,
			); errscandetail != nil {
				log.Error("Error while scanning data")
				log.Fatal(errscandetail)
			}
			order.OrderDetails = append(order.OrderDetails, orderDetail)
			// orderDetails = append(orderDetails, orderDetail)
		}

		orders = append(orders, order)
		defer orderDetailRows.Close()
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	util.HandleSuccess(c, orders, util.GET_REQUEST_SUCCESS, http.StatusOK)
}
