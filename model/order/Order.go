package model

type (
	Order struct {
		OrderID    string `json:"order_id"`
		CustomerID string `json:"customer_id" example:"C0291jjda23" swaggertype:"string"`
		OrderDate  string `json:"order_date"`
	}

	OrderDetail struct {
		OrderDetailID    string `json:"order_detail_id"`
		OrderID          string `json:"-"`
		MenuID           string `json:"menu_id"  example:"M012jidw" swaggertype:"string"`
		CurrentMenuPrice int    `json:"current_price"`
		Qty              int    `json:"qty"  example:"2"`
		Subtotal         int    `json:"subtotal"`
	}

	OrderJSON struct {
		OrderID      string        `json:"order_id"`
		CustomerID   string        `json:"customer_id" example:"C012333nc2" swaggertype:"string"`
		OrderDate    string        `json:"order_date"`
		OrderDetails []OrderDetail `json:"order_details"`
	}
)
