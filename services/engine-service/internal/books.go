package internal

// Common
// SellBook BuyBook 组合 Common 避免重复代码
type Common struct {
	data []*TrustOrder
}

func (s *Common) Len() int {
	return len(s.data)
}

func (s *Common) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

func (s *Common) Push(x any) {
	s.data = append(s.data, x.(*TrustOrder))
}

func (s *Common) Pop() any {
	if s.Len() <= 0 {
		return nil
	}
	old := s.data
	n := len(old)
	x := old[n-1]
	s.data = old[0 : n-1]
	return x
}

// SellBook 卖盘 价格从低到高
type SellBook struct {
	Common
}

func (s *SellBook) Less(i, j int) bool {
	return s.data[i].Price.Cmp(s.data[j].Price) < 0
}

// BuyBook 买盘 价格从高到低
type BuyBook struct {
	Common
}

func (b *BuyBook) Less(i, j int) bool {
	return b.data[i].Price.Cmp(b.data[j].Price) > 0
}
