package models

type OrderItem struct {
	ID         uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID    uint    `json:"order_id"`
	ProductID  uint    `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
}
