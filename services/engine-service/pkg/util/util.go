package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GetNum() string {
	// 获取当前时间的年月日时分秒
	timestamp := time.Now().Format("20060102150405")

	// 生成一个随机数作为订单号的一部分
	randomNumber := rand.Intn(1000) // 生成一个范围在 0 到 999 之间的随机数

	// 拼接年月日时分秒和随机数，组成订单号
	return fmt.Sprintf("%s%d", timestamp, randomNumber)
}
