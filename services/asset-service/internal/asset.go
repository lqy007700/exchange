package internal

import (
	"asset-service/asset-service/repository/model"
	"asset-service/asset-service/repository/mysql"
	"asset-service/asset-service/repository/redis"
	"github.com/golang/groupcache/singleflight"
	"go-micro.dev/v4/logger"
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

var (
	single singleflight.Group
)

const (
	SingleFlightGetUserAssetKey = "single:flight:get:user:asset"
)

type AssetService struct {
	db    *mysql.AssetDB
	redis *redis.AssetCache
}

func NewAssetService(db *mysql.AssetDB, redisC *redis.AssetCache) *AssetService {
	return &AssetService{
		db:    db,
		redis: redisC,
	}
}

// GetUserAsset 获取用户资产信息
// 1. 从redis中读取
// 2. redis中不存在，从mysql中读取，写入redis
// 3. 数据不存在返回错误
func (u *AssetService) GetUserAsset(userId int64, coin string) (*model.Asset, error) {
	// 读缓存
	asset, err := u.redis.GetUserAsset(userId, coin)
	if err != nil {
		return nil, err
	}

	// todo 缓存穿透
	if asset != nil {
		return asset, nil
	}

	// 防止单机缓存穿透
	do, err := single.Do(SingleFlightGetUserAssetKey, func() (interface{}, error) {
		// 读DB写缓存
		userAsset, err := u.db.GetUserAsset(userId, coin)
		if err != nil {
			return nil, err
		}

		if userAsset == nil {
			return &model.Asset{}, nil
		}
		_ = u.redis.SetUserAsset(userId, coin, userAsset)
		return userAsset, nil
	})
	if err != nil {
		logger.Errorf("SingleFlightGetUserAssetKey err: %+v", err)
		return nil, err
	}

	return do.(*model.Asset), nil
}

// GetUserAssets 获取用户所有资产信息
// 1. 只从redis中读取
// 2. 异步 redis中不存在，从mysql中读取，写入redis
func (u *AssetService) GetUserAssets(uid int64) ([]*model.Asset, error) {
	assets, err := u.redis.GetUserAssets(uid)
	if err != nil {
		return nil, err
	}

	// todo 缓存不存在用户资产信息 如何处理
	// 1 异步从db获取 写入缓存 ？
	// 2 同步从db获取 写入缓存 ？
	// 3 定时任务从db获取 写入缓存 ？
	if len(assets) == 0 {

	}

	return assets, nil
}
