package internal

import (
	"math/big"
)

/**
用户资产操作
- 资产查询
	- 查询用户所有资产
	- 查询用户指定币种资产
- 资产转入
- 资产转出
- 资产冻结
- 资产解冻


资产信息存储在DB 和 redis中，key为用户id，value为用户资产信息
对用户资产的操作需要读写锁
读 redis - mysql
写 mysql - del redis
*/

type Asset struct {
	Available *big.Float `json:"available"`
	Frozen    *big.Float `json:"frozen"`
}

type UserAsset struct {
	UserId int64
	Assets map[int32]*Asset
}

func (u *UserAsset) GetAsset(coinId int32) *Asset {
	// 读缓存

	// 读DB写缓存
	return u.Assets[coinId]
}
