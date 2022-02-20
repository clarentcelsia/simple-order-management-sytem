package model

import "restaurant/model"

type (
	Menu struct {
		MenuID          string     `json:"menu_id"`
		MenuCategoryID  int        `json:"menu_category_id"`
		MenuName        string     `json:"menu_name"`
		MenuDescription string     `json:"menu_description"`
		MenuPrice       int        `json:"menu_price"`
		Base            model.Base `json:"base"`
	}

	//Food or Beverage
	MenuCategory struct {
		MenuCategoryID   int    `json:"menu_category_id"`
		MenuCategoryName string `json:"menu_category_name"`
	}
)
