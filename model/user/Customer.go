package model

import (
	mBase "restaurant/model"
)

type (
	Customer struct {
		CustomerID      string     `json:"customer_id"`
		CustomerName    string     `json:"customer_name" example:"Nana" swaggertype:"string"`
		CustomerEmail   string     `json:"customer_email" example:"nana@gmail.com" swaggertype:"string"`
		CustomerPhone   string     `json:"customer_phone" example:"0212292012" swaggertype:"string"`
		CustomerDOB     string     `json:"customer_dob" example:"2022-01-22" swaggertype:"string"`
		CustomerAddress string     `json:"customer_address" example:"Rangola street 21" swaggertype:"string"`
		Base            mBase.Base `json:"base"`
		IsDeleted       bool       `json:"is_deleted"`
	}
)
