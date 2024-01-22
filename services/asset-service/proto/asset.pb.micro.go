// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/asset.proto

package asset_service

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for AssetService service

func NewAssetServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for AssetService service

type AssetService interface {
	// 获取用户资产
	GetAsset(ctx context.Context, in *GetAssetReq, opts ...client.CallOption) (*Asset, error)
	// 获取用户所有资产
	GetAssets(ctx context.Context, in *GetAssetsReq, opts ...client.CallOption) (*GetAssetsResp, error)
	// 冻结&解冻
	TryFreeze(ctx context.Context, in *FreezeReq, opts ...client.CallOption) (*CommonResp, error)
	UnFreeze(ctx context.Context, in *FreezeReq, opts ...client.CallOption) (*CommonResp, error)
}

type assetService struct {
	c    client.Client
	name string
}

func NewAssetService(name string, c client.Client) AssetService {
	return &assetService{
		c:    c,
		name: name,
	}
}

func (c *assetService) GetAsset(ctx context.Context, in *GetAssetReq, opts ...client.CallOption) (*Asset, error) {
	req := c.c.NewRequest(c.name, "AssetService.GetAsset", in)
	out := new(Asset)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assetService) GetAssets(ctx context.Context, in *GetAssetsReq, opts ...client.CallOption) (*GetAssetsResp, error) {
	req := c.c.NewRequest(c.name, "AssetService.GetAssets", in)
	out := new(GetAssetsResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assetService) TryFreeze(ctx context.Context, in *FreezeReq, opts ...client.CallOption) (*CommonResp, error) {
	req := c.c.NewRequest(c.name, "AssetService.TryFreeze", in)
	out := new(CommonResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assetService) UnFreeze(ctx context.Context, in *FreezeReq, opts ...client.CallOption) (*CommonResp, error) {
	req := c.c.NewRequest(c.name, "AssetService.UnFreeze", in)
	out := new(CommonResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AssetService service

type AssetServiceHandler interface {
	// 获取用户资产
	GetAsset(context.Context, *GetAssetReq, *Asset) error
	// 获取用户所有资产
	GetAssets(context.Context, *GetAssetsReq, *GetAssetsResp) error
	// 冻结&解冻
	TryFreeze(context.Context, *FreezeReq, *CommonResp) error
	UnFreeze(context.Context, *FreezeReq, *CommonResp) error
}

func RegisterAssetServiceHandler(s server.Server, hdlr AssetServiceHandler, opts ...server.HandlerOption) error {
	type assetService interface {
		GetAsset(ctx context.Context, in *GetAssetReq, out *Asset) error
		GetAssets(ctx context.Context, in *GetAssetsReq, out *GetAssetsResp) error
		TryFreeze(ctx context.Context, in *FreezeReq, out *CommonResp) error
		UnFreeze(ctx context.Context, in *FreezeReq, out *CommonResp) error
	}
	type AssetService struct {
		assetService
	}
	h := &assetServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&AssetService{h}, opts...))
}

type assetServiceHandler struct {
	AssetServiceHandler
}

func (h *assetServiceHandler) GetAsset(ctx context.Context, in *GetAssetReq, out *Asset) error {
	return h.AssetServiceHandler.GetAsset(ctx, in, out)
}

func (h *assetServiceHandler) GetAssets(ctx context.Context, in *GetAssetsReq, out *GetAssetsResp) error {
	return h.AssetServiceHandler.GetAssets(ctx, in, out)
}

func (h *assetServiceHandler) TryFreeze(ctx context.Context, in *FreezeReq, out *CommonResp) error {
	return h.AssetServiceHandler.TryFreeze(ctx, in, out)
}

func (h *assetServiceHandler) UnFreeze(ctx context.Context, in *FreezeReq, out *CommonResp) error {
	return h.AssetServiceHandler.UnFreeze(ctx, in, out)
}