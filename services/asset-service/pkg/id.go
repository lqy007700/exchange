package pkg

import "time"

func GetId(seqId int64) uint64 {
	// 获取当前时间
	currentTime := time.Now()

	// 获取年份和月份
	year := int64(currentTime.Year())
	month := int64(currentTime.Month())
	return uint64(seqId*10000 + (year*100 + month))
}
