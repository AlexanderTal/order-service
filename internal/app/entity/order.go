package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

const (
	tableNameOrder     = "orders"
	tableNameOrderItem = "order_items"
)

type Order struct {
	ID            int64       `gorm:"column:id;primaryKey;autoIncrement"`
	GUID          uuid.UUID   `gorm:"column:guid;type:uuid;not null;unique"`
	UserGUID      *uuid.UUID  `gorm:"column:user_guid;type:uuid"`
	TotalPrice    int64       `gorm:"column:total_price;not null"`
	DeliveryPrice int64       `gorm:"column:delivery_price;not null;default:0"`
	Currency      string      `gorm:"column:currency;type:varchar(3);not null"`
	Status        string      `gorm:"column:status;type:varchar(32);not null"`
	CreatedAt     time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time   `gorm:"column:updated_at;autoUpdateTime"`
	Items         []OrderItem `gorm:"foreignKey:OrderGUID;references:GUID"`
}

func (Order) TableName() string { return tableNameOrder }

type OrderItem struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	GUID        uuid.UUID `gorm:"column:guid;type:uuid;not null;unique"`
	OrderGUID   uuid.UUID `gorm:"column:order_guid;type:uuid;not null"`
	ProductGUID uuid.UUID `gorm:"column:product_guid;type:uuid;not null"`
	Quantity    int       `gorm:"column:quantity;not null"`
	UnitPrice   int64     `gorm:"column:unit_price;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (OrderItem) TableName() string { return tableNameOrderItem }

////////////////////////////////////////////////////////////////////////////////
///// HTTP REQUEST & RESPONSE //////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type RequestOrderCreate struct {
	UserGUID *uuid.UUID               `json:"user_guid" binding:"omitempty"`
	Currency string                   `json:"currency"  binding:"required,len=3"`
	Items    []RequestOrderItemCreate `json:"items"     binding:"required,min=1,dive"`
}

type RequestOrderItemCreate struct {
	ProductGUID uuid.UUID `json:"product_guid" binding:"required"`
	Quantity    int       `json:"quantity"     binding:"required,gt=0"`
	UnitPrice   int64     `json:"unit_price"   binding:"required,gt=0"`
}

type RequestOrderUpdate struct {
	Status string `json:"status" binding:"required"`
}

type RequestOrderList struct {
	Status   *string    `json:"status"    binding:"omitempty"`
	UserGUID *uuid.UUID `json:"user_guid" binding:"omitempty"`
}

type ResponseOrderCreate struct {
	GUID          uuid.UUID           `json:"guid"`
	UserGUID      *uuid.UUID          `json:"user_guid"`
	TotalPrice    int64               `json:"total_price"`
	DeliveryPrice int64               `json:"delivery_price"`
	Currency      string              `json:"currency"`
	Status        string              `json:"status"`
	CreatedAt     time.Time           `json:"created_at"`
	Items         []ResponseOrderItem `json:"items"`
}

type ResponseOrderGet struct {
	GUID          uuid.UUID           `json:"guid"`
	UserGUID      *uuid.UUID          `json:"user_guid"`
	TotalPrice    int64               `json:"total_price"`
	DeliveryPrice int64               `json:"delivery_price"`
	Currency      string              `json:"currency"`
	Status        string              `json:"status"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Items         []ResponseOrderItem `json:"items"`
}

type ResponseOrderUpdate struct {
	GUID          uuid.UUID           `json:"guid"`
	UserGUID      *uuid.UUID          `json:"user_guid"`
	TotalPrice    int64               `json:"total_price"`
	DeliveryPrice int64               `json:"delivery_price"`
	Currency      string              `json:"currency"`
	Status        string              `json:"status"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Items         []ResponseOrderItem `json:"items"`
}

type ResponseOrderItem struct {
	GUID        uuid.UUID `json:"guid"`
	ProductGUID uuid.UUID `json:"product_guid"`
	Quantity    int       `json:"quantity"`
	UnitPrice   int64     `json:"unit_price"`
}

type ResponseOrderList struct {
	Data []ResponseOrderListItem `json:"data"`
}

type ResponseOrderListItem struct {
	GUID          uuid.UUID  `json:"guid"`
	UserGUID      *uuid.UUID `json:"user_guid"`
	TotalPrice    int64      `json:"total_price"`
	DeliveryPrice int64      `json:"delivery_price"`
	Currency      string     `json:"currency"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
