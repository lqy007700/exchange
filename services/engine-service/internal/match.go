package internal

import (
	"github.com/lqy007700/exchange/common/order"
	"math/big"
)

// MatchResult 成交结果
type MatchResult struct {
	matchDetails []*MatchResultDetail
}
type MatchResultDetail struct {
	Price      *big.Float
	Quantity   *big.Float
	TakerOrder *order.OrderEntity
	MakerOrder *order.OrderEntity
}

func newMatchResult() *MatchResult {
	return &MatchResult{
		matchDetails: make([]*MatchResultDetail, 0, 4),
	}
}

func (m *MatchResult) add(price *big.Float, quantity *big.Float, takerOrder *order.OrderEntity, makerOrder *order.OrderEntity) {
	m.matchDetails = append(m.matchDetails, &MatchResultDetail{
		Price:      price,
		Quantity:   quantity,
		TakerOrder: takerOrder,
		MakerOrder: makerOrder,
	})
}

// 发送撮合结果给 kafka
func (m *MatchResult) sendMatchResToQueue() {
	//if m == nil {
	//	logger.Info("match result is nil")
	//	return
	//}
	//
	//var wg sync.WaitGroup
	//
	//for _, detail := range m.matchDetails {
	//	go func(info *MatchResultDetail) {
	//		wg.Add(1)
	//		defer wg.Done()
	//		//detailJson, err := json.Marshal(info)
	//		if err != nil {
	//			logger.Errorf("marshal match result detail error: %v", err)
	//			return
	//		}
	//		//err = m.mq.Produce("match_result", detailJson)
	//		if err != nil {
	//			logger.Errorf("send match result error: %v", err)
	//			return
	//		}
	//	}(detail)
	//}
	//wg.Wait()
}
