package order

import "fmt"

type Order struct {
	ID    string      `json:"id" example:"al888888-9999-0000-aaaa-bbbbbbbbbbbb"`
	Items []OrderItem `json:"order_item"`
	Table int32       `json:"table" example:"1"`
	Total float64     `json:"total" example:"100.0"`
}

type OrderItem struct {
	OrderID   string  `json:"order_id" example:"al888888-9999-0000-aaaa-bbbbbbbbbbbb"`
	Name      string  `json:"name" example:"fried rice"`
	Quantity  int     `json:"quantity" example:"1"`
	UnitPrice float64 `json:"unit_price" example:"100.0"`
}

func (o *Order) CalculateTotal() float64 {
	if len(o.Items) == 0 {
		o.Total = 0
		return 0
	}
	var sum float64
	for _, it := range o.Items {
		sum += it.UnitPrice * float64(it.Quantity)
	}
	o.Total = sum
	return sum
}

func SeedOrder(index int) Order {
	orderID := fmt.Sprintf("ORD-%03d", index)
	return Order{
		ID:    orderID,
		Table: int32(index),
		Items: []OrderItem{
			{OrderID: orderID, Name: "Pho", Quantity: 2, UnitPrice: 40000},
			{OrderID: orderID, Name: "Com chien", Quantity: 1, UnitPrice: 25000},
		},
	}
}
