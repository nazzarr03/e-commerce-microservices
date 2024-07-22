package models

type Product struct {
	ProductID uint    `json:"product_id" gorm:"primaryKey;autoIncrement"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}
