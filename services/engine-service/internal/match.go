package internal

import "math/big"

// MatchResult 成交结果
type MatchResult struct {
	matchDetails []*MatchResultDetail
}
type MatchResultDetail struct {
	Price      *big.Float
	Quantity   *big.Float
	TakerOrder *TrustOrder
	MakerOrder *TrustOrder
}

func newMatchResult() *MatchResult {
	return &MatchResult{
		matchDetails: make([]*MatchResultDetail, 0, 4),
	}
}

func (m *MatchResult) add(price *big.Float, quantity *big.Float, takerOrder *TrustOrder, makerOrder *TrustOrder) {
	m.matchDetails = append(m.matchDetails, &MatchResultDetail{
		Price:      price,
		Quantity:   quantity,
		TakerOrder: makerOrder,
		MakerOrder: makerOrder,
	})
}
