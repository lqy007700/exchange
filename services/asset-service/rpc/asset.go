package rpc

import (
	"asset-service/asset-service/internal"
	asset_service "asset-service/asset-service/proto"
	"context"
	"go-micro.dev/v4/client"
)

type AssetService struct {
	Asset *internal.AssetService
}

func (a *AssetService) GetAsset(ctx context.Context, req *asset_service.GetAssetReq, opts ...client.CallOption) (*asset_service.Asset, error) {
	// todo 验参
	asset, err := a.Asset.GetUserAsset(req.Uid, req.Coin)
	if err != nil {
		return nil, err
	}

	return &asset_service.Asset{
		Coin:      req.Coin,
		Available: asset.Available.String(),
		Frozen:    asset.Frozen.String(),
	}, nil
}

func (a *AssetService) GetAssets(ctx context.Context, req *asset_service.GetAssetsReq, opts ...client.CallOption) (*asset_service.GetAssetsResp, error) {
	// todo 验参
	asset, err := a.Asset.GetUserAssets(req.Uid)
	if err != nil {
		return nil, err
	}

	var resp []*asset_service.Asset
	for _, v := range asset {
		resp = append(resp, &asset_service.Asset{
			Coin:      v.Coin,
			Available: v.Available.String(),
			Frozen:    v.Frozen.String(),
		})
	}
	return &asset_service.GetAssetsResp{Assets: resp}, nil
}

func (a *AssetService) TryFreeze(ctx context.Context, in *asset_service.FreezeReq, opts ...client.CallOption) (*asset_service.CommonResp, error) {
	//TODO implement me
	panic("implement me")
}

func (a AssetService) UnFreeze(ctx context.Context, in *asset_service.FreezeReq, opts ...client.CallOption) (*asset_service.CommonResp, error) {
	//TODO implement me
	panic("implement me")
}
