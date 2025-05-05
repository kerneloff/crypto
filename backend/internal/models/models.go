package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Wallet struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	UserID    uint            `json:"user_id"`
	Currency  string          `json:"currency"`
	Balance   decimal.Decimal `json:"balance" gorm:"type:decimal(20,8)"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Order struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	UserID    uint            `json:"user_id"`
	Symbol    string          `json:"symbol"`
	Type      string          `json:"type"` // limit, market
	Side      string          `json:"side"` // buy, sell
	Price     decimal.Decimal `json:"price" gorm:"type:decimal(20,8)"`
	Quantity  decimal.Decimal `json:"quantity" gorm:"type:decimal(20,8)"`
	Status    string          `json:"status"` // new, filled, partially_filled, canceled
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Trade struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	OrderID   uint            `json:"order_id"`
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price" gorm:"type:decimal(20,8)"`
	Quantity  decimal.Decimal `json:"quantity" gorm:"type:decimal(20,8)"`
	Side      string          `json:"side"` // buy, sell
	CreatedAt time.Time       `json:"created_at"`
}

type OrderBookEntry struct {
	Price    decimal.Decimal `json:"price"`
	Quantity decimal.Decimal `json:"quantity"`
}

type OrderBook struct {
	Symbol string          `json:"symbol"`
	Bids   []OrderBookEntry `json:"bids"`
	Asks   []OrderBookEntry `json:"asks"`
}

type MarketData struct {
	Symbol    string          `json:"symbol"`
	LastPrice decimal.Decimal `json:"last_price"`
	Volume    decimal.Decimal `json:"volume"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Change    decimal.Decimal `json:"change"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
} 