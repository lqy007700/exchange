package order

// Status 委托单状态
type Status int32

const (
	Pending          Status = iota // 等待成交
	FullyFilled                    // 完全成交
	PartialFilled                  // 部分成交
	PartialCancelled               // 部分成交后取消
	FullyCancelled                 // 完全取消
)

// Direction 买卖方向
type Direction int

const (
	Buy Direction = iota
	Sell
)
