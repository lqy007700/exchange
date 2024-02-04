package order

import (
	"math/big"
	"time"
)

// OrderEntity 实体，用于委托单和成交单
type OrderEntity struct {
	ID               string     `json:"id"`
	UserID           int64      `json:"user_id"`
	Symbol           string     `json:"symbol"`
	Quantity         *big.Float `json:"quantity"`
	UnfilledQuantity *big.Float `json:"unfilled_quantity"`
	Price            *big.Float `json:"price"`
	Status           Status     `json:"status"`
	CreateAt         time.Time  `json:"create_at"`
	Direction        Direction  `json:"direction"`
}
