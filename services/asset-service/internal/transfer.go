package internal

import (
	"github.com/pkg/errors"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"math/big"
)

type Transfer int

const (
	AvailableToAvailable Transfer = iota // 两个账户转账
	AvailableToFrozen                    // 冻结
	FrozenToAvailable                    // 解冻
	AddToAvailable                       // 充值
	DecToAvailable                       // 提现
)

// Transfer tryTransfer 转账
// 划转 冻结 解冻 都是对用户资产的操作，放到一个 func 管理
func (u *AssetService) tryTransfer(t Transfer, fromUid int64, toUid int64, coin string, amount *big.Float) error {
	if amount.Sign() == 0 {
		logger.Infof("amount is zero")
		return nil
	}

	if amount.Sign() < 0 {
		return errors.New("amount must be positive")
	}

	// 获取用户资产
	fromAsset, err := u.GetUserAsset(fromUid, coin)
	if err != nil {
		return err
	}

	switch t {
	case AvailableToAvailable:
		if fromAsset.Available.Cmp(amount) < 0 {
			return errors.New("from user available is not enough")
		}
		return u.transferDB(fromUid, toUid, coin, amount)
	case AvailableToFrozen:
		if fromAsset.Available.Cmp(amount) < 0 {
			return errors.New("from user available is not enough")
		}
		return u.frozenDB(fromUid, coin, amount)
	case FrozenToAvailable:
		if fromAsset.Frozen.Cmp(amount) < 0 {
			return errors.New("from user frozen is not enough")
		}
		return u.unfrozen(fromUid, coin, amount)
	case AddToAvailable:
		return u.addToAvailable(fromUid, coin, amount)
	case DecToAvailable:
		if fromAsset.Available.Cmp(amount) < 0 {
			return errors.New("from user available is not enough")
		}
		return u.decToAvailable(fromUid, coin, amount)
	default:
		return errors.New("unknown transfer type")
	}
}

// transferDB 转账
func (u *AssetService) transferDB(fromUid int64, toUid int64, coin string, amount *big.Float) error {
	return u.db.Transfer(fromUid, toUid, coin, amount)
}

// frozenDB 冻结
func (u *AssetService) frozenDB(uid int64, coin string, amount *big.Float) error {
	return u.db.Frozen(uid, coin, amount)
}

func (u *AssetService) unfrozen(uid int64, coin string, amount *big.Float) error {
	return u.db.Unfrozen(uid, coin, amount)
}

func (u *AssetService) addToAvailable(uid int64, coin string, amount *big.Float) error {
	return nil
}

func (u *AssetService) decToAvailable(uid int64, coin string, amount *big.Float) error {
	return nil
}
