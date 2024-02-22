package internal

import (
	"container/heap"
	"fmt"
	"github.com/lqy007700/exchange/common/order"
	"math/big"
	"testing"
)

func TestCommon_Peek(t *testing.T) {
	a := make([]*order.OrderEntity, 0)
	a = append(a, &order.OrderEntity{
		ID:    "1",
		Price: big.NewFloat(1),
	})
	a = append(a, &order.OrderEntity{
		ID:    "2",
		Price: big.NewFloat(2),
	})
	a = append(a, &order.OrderEntity{
		ID:    "3",
		Price: big.NewFloat(3),
	})

	buy := &SellBook{Common{data: a}}
	heap.Init(buy)

	fmt.Println(buy.Peek())

}
