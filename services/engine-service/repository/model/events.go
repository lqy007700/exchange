package model

import "time"

type Events struct {
	SeqId    int64         `json:"seq_id"`
	PreId    int64         `json:"pre_id"`
	CreateAt time.Duration `json:"create_at"`
}
