package sequencer

import (
	"github.com/jinzhu/gorm"
	"go-micro.dev/v4/logger"
	"math/big"
)

type Seq struct {
	Id  *big.Int `gorm:"column:id"`
	db  *gorm.DB
	inc *big.Int
}

func NewSeq(db *gorm.DB) *Seq {
	s := &Seq{
		db:  db,
		inc: big.NewInt(1),
	}
	s.GetLastId()
	return s
}

func (s *Seq) GetLastId() {
	var lastId *big.Int
	result := s.db.Table("seq").Select("id").Order("seq_id desc").First(&lastId)
	if result.Error != nil {
		logger.Errorf("get last id error: %v", result.Error)
	}

	s.Id = lastId
}

func (s *Seq) GetId() {
	s.Id.Add(s.Id, s.inc)
}
