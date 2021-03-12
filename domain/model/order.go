package model

import "time"

type Order struct {
	ID          int64         `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	OrderCode   string        `gorm:"uniqueIndex;not null;size:255" json:"order_code"`
	PayStatus   int32         `json:"pay_status"`
	ShipStatus  int32         `json:"ship_status"`
	Price       float64       `json:"price"`
	OrderDetail []OrderDetail `gorm:"foreignKey:OrderID" json:"order_detail"`
	CreateAt    time.Time
	UpdateAt    time.Time
}
