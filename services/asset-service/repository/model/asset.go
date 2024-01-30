package model

import "math/big"

type Asset struct {
	Coin      string     `json:"coin"`
	Available *big.Float `json:"available"`
	Frozen    *big.Float `json:"frozen"`
}
