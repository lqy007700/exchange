package redis

import (
	"asset-service/config"
	"asset-service/repository/model"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"go-micro.dev/v4/logger"
	"time"
)

const (
	UserAssetHashKey = "user:asset:hash:%d"
)

type AssetCache struct {
	client *redis.Client
}

func NewAssetCache() *AssetCache {
	conf := config.Conf.Redis
	c := &AssetCache{
		client: redis.NewClient(&redis.Options{
			Addr:         conf.Addr,
			Password:     conf.Password,
			DB:           conf.DB,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleConns,
			DialTimeout:  time.Duration(conf.DialTimeout),
			ReadTimeout:  time.Duration(conf.ReadTimeout),
			WriteTimeout: time.Duration(conf.WriteTimeout),
			IdleTimeout:  time.Duration(conf.IdleTimeout),
		}),
	}

	if _, err := c.client.Ping().Result(); err != nil {
		panic(err)
	} else {
		logger.Infof("redis connected.")
	}
	return c
}

func (c *AssetCache) Ping() error {
	return c.client.Ping().Err()
}

func (c *AssetCache) Close() error {
	err := c.client.Close()
	return err
}

func (c *AssetCache) GetUserAsset(userId int64, coin string) (*model.Asset, error) {
	key := fmt.Sprintf(UserAssetHashKey, userId)
	result, err := c.client.HGet(key, coin).Result()
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, errors.New("asset not found for cache")
	}

	assetResp := &model.Asset{}
	err = json.Unmarshal([]byte(result), assetResp)
	if err != nil {
		logger.Errorf("json unmarshal error: %v", err)
		return nil, err
	}
	return assetResp, nil
}

func (c *AssetCache) SetUserAsset(userId int64, coin string, val *model.Asset) error {
	key := fmt.Sprintf(UserAssetHashKey, userId)
	marshal, err := json.Marshal(val)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = c.client.HSet(key, coin, marshal).Err()
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (c *AssetCache) GetUserAssets(uid int64) ([]*model.Asset, error) {
	key := fmt.Sprintf(UserAssetHashKey, uid)
	result, err := c.client.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	var assets []*model.Asset
	if len(result) > 0 {
		for _, value := range result {
			tmp := &model.Asset{}
			err = json.Unmarshal([]byte(value), tmp)
			if err != nil {
				logger.Errorf("json unmarshal error: %v", err)
				continue
			}
			assets = append(assets, tmp)
		}
	}
	return assets, nil
}

func (c *AssetCache) CleanUserAsset(uid int64, coin string) {
	key := fmt.Sprintf(UserAssetHashKey, uid)
	err := c.client.HDel(key, coin).Err()
	if err != nil {
		logger.Errorf("uid: %d, coin: %s, cleanUserAsset error: %v", uid, coin, err)
	}
}
