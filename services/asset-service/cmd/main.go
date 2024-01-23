package main

import (
	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlDSN  = "user:password@tcp(localhost:3306)/dbname"
	redisAddr = "localhost:6379"
)

//// Asset represents a user's asset in a specific currency.
//type Asset struct {
//	Currency string
//	Balance  float64
//	Frozen   float64
//}
//
//// UserService provides methods to interact with user assets.
//type UserService struct {
//	db    *sql.DB
//	redis *redis.Client
//}
//
//// NewUserService creates a new instance of UserService.
//func NewUserService() *UserService {
//	// Initialize MySQL connection
//	db, err := sql.Open("mysql", mysqlDSN)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Initialize Redis connection
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: redisAddr,
//		DB:   0,
//	})
//
//	// Ping to check if Redis is reachable
//	if _, err := redisClient.Ping().Result(); err != nil {
//		log.Fatal(err)
//	}
//
//	return &UserService{
//		db:    db,
//		redis: redisClient,
//	}
//}
//
//// GetUserAsset retrieves user's asset for a specific currency from MySQL and Redis.
//func (us *UserService) GetUserAsset(userID, currency string) (*Asset, error) {
//	// Try to get asset from Redis
//	redisKey := fmt.Sprintf("%s:%s", userID, currency)
//	redisResult, err := us.redis.HMGet(redisKey, "balance", "frozen").Result()
//	if err == nil && redisResult[0] != nil && redisResult[1] != nil {
//		balance, _ := redis.Float64(redisResult[0], nil)
//		frozen, _ := redis.Float64(redisResult[1], nil)
//		return &Asset{Currency: currency, Balance: balance, Frozen: frozen}, nil
//	}
//
//	// If not found in Redis, fetch from MySQL
//	query := "SELECT balance, frozen FROM user_assets WHERE user_id = ? AND currency = ?"
//	row := us.db.QueryRow(query, userID, currency)
//
//	var balance, frozen float64
//	if err := row.Scan(&balance, &frozen); err != nil {
//		return nil, err
//	}
//
//	// Cache the result in Redis
//	us.redis.HMSet(redisKey, map[string]interface{}{"balance": balance, "frozen": frozen})
//
//	return &Asset{Currency: currency, Balance: balance, Frozen: frozen}, nil
//}
//
//// GetUserAllAssets retrieves all user assets from MySQL and Redis.
//func (us *UserService) GetUserAllAssets(userID string) ([]*Asset, error) {
//	// Try to get all assets from Redis
//	redisKey := fmt.Sprintf("%s:all_assets", userID)
//	redisResult, err := us.redis.HGetAll(redisKey).Result()
//	if err == nil && len(redisResult) > 0 {
//		var assets []*Asset
//		for currency, value := range redisResult {
//			balance, _ := redis.Float64(value, nil)
//			assets = append(assets, &Asset{Currency: currency, Balance: balance})
//		}
//		return assets, nil
//	}
//
//	// If not found in Redis, fetch from MySQL
//	query := "SELECT currency, balance, frozen FROM user_assets WHERE user_id = ?"
//	rows, err := us.db.Query(query, userID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var assets []*Asset
//	for rows.Next() {
//		var currency string
//		var balance, frozen float64
//		if err := rows.Scan(&currency, &balance, &frozen); err != nil {
//			return nil, err
//		}
//		assets = append(assets, &Asset{Currency: currency, Balance: balance, Frozen: frozen})
//	}
//
//	// Cache the result in Redis
//	redisMap := make(map[string]interface{})
//	for _, asset := range assets {
//		redisMap[asset.Currency] = asset.Balance
//	}
//	us.redis.HMSet(redisKey, redisMap)
//
//	return assets, nil
//}
//
//// FreezeAsset freezes a specified amount of user's asset for a specific currency.
//func (us *UserService) FreezeAsset(userID, currency string, amount float64) error {
//	// Check if the balance is sufficient
//	asset, err := us.GetUserAsset(userID, currency)
//	if err != nil {
//		return err
//	}
//	if asset.Balance < amount {
//		return fmt.Errorf("insufficient balance to freeze")
//	}
//
//	// Update MySQL
//	query := "UPDATE user_assets SET balance = balance - ?, frozen = frozen + ? WHERE user_id = ? AND currency = ?"
//	_, err = us.db.Exec(query, amount, amount, userID, currency)
//	if err != nil {
//		return err
//	}
//
//	// Update Redis
//	redisKey := fmt.Sprintf("%s:%s", userID, currency)
//	us.redis.HIncrByFloat(redisKey, "balance", -amount)
//	us.redis.HIncrByFloat(redisKey, "frozen", amount)
//
//	return nil
//}
//
//// UnfreezeAsset unfreezes a specified amount of user's asset for a specific currency.
//func (us *UserService) UnfreezeAsset(userID, currency string, amount float64) error {
//	// Check if the frozen amount is sufficient
//	asset, err := us.GetUserAsset(userID, currency)
//	if err != nil {
//		return err
//	}
//	if asset.Frozen < amount {
//		return fmt.Errorf("insufficient frozen amount to unfreeze")
//	}
//
//	// Update MySQL
//	query := "UPDATE user_assets SET frozen = frozen - ? WHERE user_id = ? AND currency = ?"
//	_, err = us.db.Exec(query, amount, userID, currency)
//	if err != nil {
//		return err
//	}
//
//	// Update Redis
//	redisKey := fmt.Sprintf("%s:%s", userID, currency)
//	us.redis.HIncrByFloat(redisKey, "frozen", -amount)
//
//	return nil
//}
//
//// TransferAsset transfers a specified amount of asset from one user to another for a specific currency.
//func (us *UserService) TransferAsset(fromUserID, toUserID, currency string, amount float64) error {
//	// Check if the balance is sufficient
//	asset, err := us.GetUserAsset(fromUserID, currency)
//	if err != nil {
//		return err
//	}
//	if asset.Balance < amount {
//		return fmt.Errorf("insufficient balance to transfer")
//	}
//
//	// Update MySQL for the sender
//	querySender := "UPDATE user_assets SET balance = balance - ? WHERE user_id = ? AND currency = ?"
//	_, err =
