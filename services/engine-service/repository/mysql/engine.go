package mysql

import (
	"github.com/jinzhu/gorm"
)

const (
	assetTableName = ""
)

type EngineDB struct {
	DB *gorm.DB
}

func New() *EngineDB {
	repo := &EngineDB{DB: newDB()}
	return repo
}

func (r *EngineDB) Close() {
	_ = r.DB.Close()
}

//// GetOrderBooks  获取订单簿
//func (r *EngineDB) GetOrderBooks() (*model.Asset, error) {
//	table := fmt.Sprintf(assetTableName, coin)
//	resp := &model.Asset{}
//	err := r.DB.Table(table).Where("user_id = ?", userId).Find(resp).Error
//	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		return nil, err
//	}
//	return resp, nil
//}
