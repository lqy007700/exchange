package order

import (
	"math/big"
	"time"
)

// OrderEntity 实体，用于委托单和成交单
type OrderEntity struct {
	ID               string     `json:"id"`
	UserID           int64      `json:"user_id"`
	CoinFrom         string     `json:"coin_from"`
	CoinTo           string     `json:"coin_to"`
	Direction        Direction  `json:"direction"`
	Price            *big.Float `json:"price"`
	Quantity         *big.Float `json:"quantity"`
	UnfilledQuantity *big.Float `json:"unfilled_quantity"`
	Status           Status     `json:"status"`
	CreateAt         time.Time  `json:"create_at"`
}
