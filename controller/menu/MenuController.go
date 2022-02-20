package controller

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"restaurant/config"
	mOrder "restaurant/model/order"
	"restaurant/util"
)

func CreateMenuCategory(c *gin.Context) {
	db, errconn := config.GetConnection()
	if errconn != nil {
		log.Info("Connection not found")
		log.Fatal(errconn)
	}

	var MenuCategoryParam map[string]interface{}
	var menuCategory mOrder.MenuCategory
	errbind := c.BindJSON(&menuCategory) //&value referred to memory address of menucategory

	if errbind != nil {
		log.Error(util.INSERT_RES_FAILED)
		util.HandleError(c, errbind, util.MISSING_REQUEST, "Make sure you make correct request", http.StatusBadRequest)
	} else {
		query := `INSERT INTO menucategory (menu_category_name) VALUES ($1) RETURNING *`
		err := db.QueryRow(query, menuCategory.MenuCategoryName).Scan(&menuCategory.MenuCategoryID, &menuCategory.MenuCategoryName)
		if err != nil {
			log.Error(err)
		} else {
			MenuCategoryParam = map[string]interface{}{
				"menu_category_id":   menuCategory.MenuCategoryID,
				"menu_category_name": menuCategory.MenuCategoryName,
			}
			util.HandleSuccess(c, MenuCategoryParam, util.INSERT_RES_SUCCEESS, http.StatusCreated)
		}
	}
}

func CreateMenu(c *gin.Context) {
	db, errconn := config.GetConnection()
	if errconn != nil {
		log.Info("Connection not found")
		log.Fatal(errconn)
	}

	var MenuParam map[string]interface{}
	var menu mOrder.Menu
	errbind := c.BindJSON(&menu)

	if errbind != nil {
		log.Error(util.INSERT_RES_FAILED)
		util.HandleError(c, errbind, util.MISSING_REQUEST, "Make sure you make correct request", http.StatusBadRequest)
	} else {
		query := `INSERT INTO menu(menu_id, menu_category_id, menu_name, menu_description, menu_price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`
		id := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)
		if errPersist := menu.Base.PrePersist(); errPersist != nil {
			log.Error(errPersist)
		}
		err := db.QueryRow(query, id, menu.MenuCategoryID, menu.MenuName, menu.MenuDescription, menu.MenuPrice, menu.Base.CreatedAt, menu.Base.UpdatedAt).
			Scan(&menu.MenuID, &menu.MenuCategoryID, &menu.MenuName, &menu.MenuDescription, &menu.MenuPrice, &menu.Base.CreatedAt, &menu.Base.UpdatedAt)
		if err != nil {
			log.Error(err)
		} else {
			MenuParam = map[string]interface{}{
				"menu_id":          menu.MenuID,
				"menu_category_id": menu.MenuCategoryID,
				"menu_name":        menu.MenuName,
				"menu_description": menu.MenuDescription,
				"menu_price":       menu.MenuPrice,
				"crt_at":           menu.Base.CreatedAt,
				"upd_at":           menu.Base.UpdatedAt,
			}
			util.HandleSuccess(c, MenuParam, util.INSERT_RES_SUCCEESS, http.StatusCreated)
		}
	}
}

func FindByID(c *gin.Context) {
	db, errconn := config.GetConnection()
	if errconn != nil {
		log.Error("Connection not found")
		log.Fatal(errconn)
	} else {
		paramID := c.Param("id")
		menu, err := GetMenuByID(c, db, paramID)
		if err != nil {
			log.Error(err)
			util.HandleError(c, err, util.GET_REQUEST_FAILED, "Can't find the given id", http.StatusNotFound)
		} else {
			util.HandleSuccess(c, menu, util.GET_REQUEST_SUCCESS, http.StatusOK)
		}
	}
}

func GetMenuByID(c *gin.Context, db *sql.DB, id string) (mOrder.Menu, error) {
	var menu mOrder.Menu
	query := `SELECT * FROM menu WHERE menu_id=$1`
	err := db.QueryRow(query, id).Scan(&menu.MenuID, &menu.MenuCategoryID, &menu.MenuName, &menu.MenuDescription, &menu.MenuPrice, &menu.Base.CreatedAt, &menu.Base.UpdatedAt)
	if err != nil {
		log.Error(err)
	}
	return menu, nil
}

func GetMenus(c *gin.Context) {
	db, errconn := config.GetConnection()
	if errconn != nil {
		log.Fatal(errconn)
	} else {
		var (
			menus []mOrder.Menu
			menu  mOrder.Menu
		)

		query := `SELECT * FROM menu`
		rows, errdb := db.Query(query)
		if errdb != nil {
			log.Error(errdb)
			util.HandleError(c, errdb, util.GET_REQUEST_FAILED, "Please check the code", http.StatusInternalServerError)
		} else {
			for rows.Next() {
				err := rows.Scan(&menu.MenuID, &menu.MenuCategoryID, &menu.MenuName, &menu.MenuDescription, &menu.MenuPrice, &menu.Base.CreatedAt, &menu.Base.UpdatedAt)
				if err != nil {
					log.Error("Error while scanning data")
					log.Fatal(err)
				}
				menus = append(menus, menu)
			}
		}
		defer rows.Close()

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		util.HandleSuccess(c, menus, util.GET_REQUEST_SUCCESS, http.StatusOK)
	}
}
