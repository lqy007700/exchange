package mysql

import "asset-service/asset-service/repository/model"

func (r *AssetDB) getId() (int64, error) {
	res := &model.Events{}
	first := r.DB.Order("seq_id desc").First(res)
	if first.Error != nil {
		return 0, first.Error
	}

	return res.SeqId, nil
}
