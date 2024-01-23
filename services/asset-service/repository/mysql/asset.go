package mysql

import (
	"asset-service/asset-service/config"
	"asset-service/asset-service/repository/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
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
