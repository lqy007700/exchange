package rpc

import (
	"asset-service/internal"
	asset_service "asset-service/proto"
	context "context"
	"go-micro.dev/v4/logger"
)

type AssetService struct {
	Asset *internal.AssetService
	Order *internal.OrderService
}

func (a *AssetService) GetAsset(ctx context.Context, req *asset_service.GetAssetReq, resp *asset_service.Asset) error {
	// todo 验参

	userAsset, err := a.Asset.GetUserAsset(req.Uid, req.Coin)
	if err != nil {
		return err
	}

	resp.Available = userAsset.Available.String()
	resp.Frozen = userAsset.Frozen.String()
	resp.Coin = userAsset.Coin
	return nil
}

func (a *AssetService) GetAssets(ctx context.Context, req *asset_service.GetAssetsReq, resp *asset_service.GetAssetsResp) error {
	resp.Assets = make([]*asset_service.Asset, 0)
	resp.Assets = append(resp.Assets, &asset_service.Asset{
		Coin:      "BTC",
		Available: "100",
		Frozen:    "0",
	})
	logger.Info(111111111)
	return nil
}

func (a *AssetService) Freeze(ctx context.Context, req *asset_service.FreezeReq, resp *asset_service.CommonResp) error {
	//TODO implement me
	panic("implement me")
}

func (a *AssetService) UnFreeze(ctx context.Context, req *asset_service.FreezeReq, resp *asset_service.CommonResp) error {
	//TODO implement me
	panic("implement me")
}

//
//func (a *AssetService) GetAsset(ctx context.Context, req *asset_service.GetAssetReq, opts ...client.CallOption) (*asset_service.Asset, error) {
//	// todo 验参
//	asset, err := a.Asset.GetUserAsset(req.Uid, req.Coin)
//	if err != nil {
//		return nil, err
//	}
//
//	return &asset_service.Asset{
//		Coin:      req.Coin,
//		Available: asset.Available.String(),
//		Frozen:    asset.Frozen.String(),
//	}, nil
//}
//
//func (a *AssetService) GetAssets(ctx context.Context, req *asset_service.GetAssetsReq, opts ...client.CallOption) (*asset_service.GetAssetsResp, error) {
//	// todo 验参
//	asset, err := a.Asset.GetUserAssets(req.Uid)
//	if err != nil {
//		return nil, err
//	}
//
//	var resp []*asset_service.Asset
//	for _, v := range asset {
//		resp = append(resp, &asset_service.Asset{
//			Coin:      v.Coin,
//			Available: v.Available.String(),
//			Frozen:    v.Frozen.String(),
//		})
//	}
//	return &asset_service.GetAssetsResp{Assets: resp}, nil
//}
//
//func (a *AssetService) Freeze(ctx context.Context, in *asset_service.FreezeReq, opts ...client.CallOption) (*asset_service.CommonResp, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (a AssetService) UnFreeze(ctx context.Context, in *asset_service.FreezeReq, opts ...client.CallOption) (*asset_service.CommonResp, error) {
//	//TODO implement me
//	panic("implement me")
//}
