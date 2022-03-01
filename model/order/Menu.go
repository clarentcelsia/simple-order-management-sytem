package model

import "restaurant/model"

type (
	Menu struct {
		MenuID          string     `json:"menu_id"`
		MenuCategoryID  int        `json:"menu_category_id" example:"37"`
		MenuName        string     `json:"menu_name" example:"Ginger bread" swaggertype:"string"`
		MenuDescription string     `json:"menu_description" example:"Sweet sour ginger bread" swaggertype:"string"`
		MenuPrice       int        `json:"menu_price" example:"12000"`
		Base            model.Base `json:"base"`
	}

	//Food or Beverage
	MenuCategory struct {
		MenuCategoryID   int    `json:"menu_category_id"`
		MenuCategoryName string `json:"menu_category_name" example:"Appetizer" swaggertype:"string"`
	}
)
