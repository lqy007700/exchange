package redis

import (
	"encoding/json"
	"engine-service/config"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/lqy007700/exchange/common/order"
	"github.com/pkg/errors"
	"go-micro.dev/v4/logger"
	"time"
)

const BaseTimestamp = 1704038400 // 2024-01-01

// 使用 redis 有序集合存储订单簿，方便前端展示深度图
// 撮合在初始化添加到内存中，撮合完成后更新到 redis 中
// score 为价格 + 时间，value 为订单 OrderEntity 的 json 字符串
// score  =  价格 + （订单时间 - BaseTimestamp） / 100000
const (
	// OrderBooksZsetBuy 买盘
	OrderBooksZsetBuy = "order:book:buy:%s" // coin_pair
	// OrderBooksZsetSell 卖盘
	OrderBooksZsetSell = "order:book:sell:%s"
)

type BooksCache struct {
	client *redis.Client
}

func NewBooksCache() *BooksCache {
	conf := config.Conf.Redis
	c := &BooksCache{
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

func (c *BooksCache) Ping() error {
	return c.client.Ping().Err()
}

func (c *BooksCache) Close() error {
	err := c.client.Close()
	return err
}

func (c *BooksCache) GetBooks(coinPair string, direction order.Direction) ([]*order.OrderEntity, error) {
	var (
		orders []*order.OrderEntity
		key    string
	)

	switch direction {
	case order.Buy:

		key = fmt.Sprintf(OrderBooksZsetBuy, coinPair)
		break
	case order.Sell:
		key = fmt.Sprintf(OrderBooksZsetSell, coinPair)
		break
	default:
		return nil, errors.New("unknown direction")
	}

	val := c.client.Get(key)
	if errors.Is(val.Err(), redis.Nil) {
		return nil, nil
	}
	if val.Err() != nil {
		logger.Errorf("get books error: %v", val.Err())
		return nil, val.Err()
	}

	err := json.Unmarshal([]byte(val.Val()), &orders)
	if err != nil {
		logger.Errorf("json unmarshal error: %v", err)
		return nil, err
	}
	return orders, nil
}

func (c *BooksCache) SetBooks(coinPair string, direction order.Direction, data []*order.OrderEntity) error {
	if len(data) == 0 || data == nil {
		return nil
	}

	var key string

	switch direction {
	case order.Buy:
		key = fmt.Sprintf(OrderBooksZsetBuy, coinPair)
		break
	case order.Sell:
		key = fmt.Sprintf(OrderBooksZsetSell, coinPair)
		break
	default:
		return errors.New("unknown direction")
	}

	str, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("json marshal error: %v", err)
		return err
	}
	err = c.client.Set(key, str, 0).Err()
	if err != nil {
		logger.Errorf("set books error: %v", err)
		return err
	}
	return nil
}
