package internal

import (
	"container/heap"
	"github.com/lqy007700/exchange/common/order"
)

type WrapHeap interface {
	heap.Interface
	Peek() any
}

// Common
// SellBook BuyBook 组合 Common 避免重复代码
type Common struct {
	data []*order.OrderEntity
}

func (s *Common) Peek() any {
	if s.Len() > 0 {
		return s.data[0]
	}
	return nil
}

func (s *Common) Len() int {
	return len(s.data)
}

func (s *Common) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

func (s *Common) Push(x any) {
	s.data = append(s.data, x.(*order.OrderEntity))
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

// Less 价格优先 价格相同时间优先
func (s *SellBook) Less(i, j int) bool {
	if s.data[i].Price.Cmp(s.data[j].Price) == 0 {
		return s.data[i].CreateAt.Before(s.data[j].CreateAt)
	}
	return s.data[i].Price.Cmp(s.data[j].Price) < 0
}

// BuyBook 买盘 价格从高到低
type BuyBook struct {
	Common
}

// Less 价格优先 时间优先
func (b *BuyBook) Less(i, j int) bool {
	if b.data[i].Price.Cmp(b.data[j].Price) == 0 {
		return b.data[i].CreateAt.Before(b.data[j].CreateAt)
	}
	return b.data[i].Price.Cmp(b.data[j].Price) > 0
}
