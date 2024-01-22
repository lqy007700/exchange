package rpc

import (
	context "context"
	"exchange/services/asset-service/internal"
	"exchange/services/asset-service/proto"
	"go-micro.dev/v4/client"
)

type AssetService struct {
	Asset *internal.Asset
}

func (a *AssetService) GetAsset(ctx context.Context, req *asset_service.GetAssetReq, opts ...client.CallOption) (*asset_service.Asset, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AssetService) GetAssets(ctx context.Context, in *asset_service.GetAssetsReq, opts ...client.CallOption) (*asset_service.GetAssetsResp, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AssetService) TryFreeze(ctx context.Context, in *asset_service.FreezeReq, opts ...client.CallOption) (*asset_service.CommonResp, error) {
	//TODO implement me
	panic("implement me")
}

func (a AssetService) UnFreeze(ctx context.Context, in *asset_service.FreezeReq, opts ...client.CallOption) (*asset_service.CommonResp, error) {
	//TODO implement me
	panic("implement me")
}
