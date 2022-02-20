package model

type (
	Order struct {
		OrderID    string `json:"order_id"`
		CustomerID string `json:"customer_id"`
		OrderDate  string `json:"order_date"`
	}

	OrderDetail struct {
		OrderDetailID    string `json:"order_detail_id"`
		OrderID          string `json:"-"`
		MenuID           string `json:"menu_id"`
		CurrentMenuPrice int    `json:"current_price"`
		Qty              int    `json:"qty"`
		Subtotal         int    `json:"subtotal"`
	}

	OrderJSON struct {
		OrderID      string        `json:"order_id"`
		CustomerID   string        `json:"customer_id"`
		OrderDate    string        `json:"order_date"`
		OrderDetails []OrderDetail `json:"order_details"`
	}
)
