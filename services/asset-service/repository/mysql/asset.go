package mysql

import (
	"asset-service/asset-service/config"
	"asset-service/asset-service/repository/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"math/big"
)

const (
	assetTableName = "user_%s"
)

type AssetDB struct {
	DB *gorm.DB
}

func New(conf *config.Config) *AssetDB {
	repo := &AssetDB{DB: newDB(conf)}
	return repo
}

func (r *AssetDB) Close() {
	_ = r.DB.Close()
}

// GetUserAsset 从DB中获取用户资产信息
func (r *AssetDB) GetUserAsset(userId int64, coin string) (*model.Asset, error) {
	table := fmt.Sprintf(assetTableName, coin)
	resp := &model.Asset{}
	err := r.DB.Table(table).Where("user_id = ?", userId).Find(resp).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return resp, nil
}

// Transfer 转账
func (r *AssetDB) Transfer(fromUid int64, toUid int64, coin string, amount *big.Float) error {
	tableName := fmt.Sprintf(assetTableName, coin)

	// 开始事务
	tx := r.DB.Begin()
	if tx.Error != nil {
		fmt.Println("Failed to start transaction:", tx.Error)
		return tx.Error
	}

	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("panic: %v", e)
			tx.Rollback()
		}
	}()

	transferSql := fmt.Sprintf(`
        UPDATE %s SET available = available - ? WHERE user_id = ? AND available >= ?;
        UPDATE %s SET available = available + ? WHERE user_id = ?;
    `, tableName, tableName)

	err := tx.Exec(transferSql, amount, fromUid, amount, amount, toUid).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("Sql %s Failed to transfer: %v", transferSql, err)
		return err
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		logger.Errorf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}

// Frozen 冻结
func (r *AssetDB) Frozen(uid int64, coin string, amount *big.Float) error {
	err := r.DB.Table(fmt.Sprintf(assetTableName, coin)).
		Where("user_id = ? and available >= ?", uid, amount).
		Update("available", gorm.Expr("available - ?", amount)).
		Update("frozen", gorm.Expr("frozen + ?", amount)).Error
	if err != nil {
		logger.Errorf("Failed to frozen: %v", err)
	}
	return err
}

// Unfrozen 解冻
func (r *AssetDB) Unfrozen(uid int64, coin string, amount *big.Float) error {
	err := r.DB.Table(fmt.Sprintf(assetTableName, coin)).
		Where("user_id = ? and frozen >= ?", uid, amount).
		Update("available", gorm.Expr("available + ?", amount)).
		Update("frozen", gorm.Expr("frozen - ?", amount)).Error
	if err != nil {
		logger.Errorf("Failed to Unfrozen: %v", err)
	}
	return err
}
