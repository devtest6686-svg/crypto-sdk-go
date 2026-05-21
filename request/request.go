package request

import (
	"context"

	"github.com/lbe-io/crypto-sdk-go/client"
	"github.com/lbe-io/crypto-sdk-go/common/consts"
	"github.com/lbe-io/crypto-sdk-go/common/entity"
)

// ICrypto 主接口
type IRequest interface {
	TransactionDetail(ctx context.Context, req *entity.TransactionDetailReq) (*entity.TransactionDetailResp, error)
	WalletList(ctx context.Context, req *entity.WalletListReq) (*entity.WalletListResp, error)
	WalletAddressGenerate(ctx context.Context, req *entity.WalletAddressGenerateReq) (*entity.WalletAddressGenerateResp, error)
	WalletAddressList(ctx context.Context, req *entity.WalletAddressListReq) (*entity.WalletAddressListResp, error)
	WalletAddressAssets(ctx context.Context, req *entity.WalletAddressAssetsReq) (*entity.WalletAddressAssetsResp, error)
}

type crypto struct {
	client *client.CryptoClient
}

func NewRequest(client *client.CryptoClient) IRequest {
	return &crypto{
		client: client,
	}
}

// 交易详情查询
func (c *crypto) TransactionDetail(ctx context.Context, req *entity.TransactionDetailReq) (*entity.TransactionDetailResp, error) {

	resp := &entity.TransactionDetailResp{}

	err := c.client.Post(ctx, consts.PathTransactionDetail, req, resp)
	if err != nil {
		c.client.Logger().Errorf(ctx, "TransactionDetail err:%v", err)
		return nil, err
	}

	return resp, nil
}

// 钱包列表查询
func (c *crypto) WalletList(ctx context.Context, req *entity.WalletListReq) (*entity.WalletListResp, error) {

	resp := &entity.WalletListResp{}

	err := c.client.Post(ctx, consts.PathWalletList, req, resp)
	if err != nil {
		c.client.Logger().Errorf(ctx, "WalletList err:%v", err)
		return nil, err
	}

	return resp, nil
}

// 钱包地址生成
func (c *crypto) WalletAddressGenerate(ctx context.Context, req *entity.WalletAddressGenerateReq) (*entity.WalletAddressGenerateResp, error) {

	resp := &entity.WalletAddressGenerateResp{}

	err := c.client.Post(ctx, consts.PathWalletAddressGenerate, req, resp)
	if err != nil {
		c.client.Logger().Errorf(ctx, "WalletAddressGenerate err:%v", err)
		return nil, err
	}

	return resp, nil
}

// 钱包地址列表查询
func (c *crypto) WalletAddressList(ctx context.Context, req *entity.WalletAddressListReq) (*entity.WalletAddressListResp, error) {

	resp := &entity.WalletAddressListResp{}

	err := c.client.Post(ctx, consts.PathWalletAddressList, req, resp)
	if err != nil {
		c.client.Logger().Errorf(ctx, "WalletAddressList err:%v", err)
		return nil, err
	}

	return resp, nil
}

// 钱包地址资产查询
func (c *crypto) WalletAddressAssets(ctx context.Context, req *entity.WalletAddressAssetsReq) (*entity.WalletAddressAssetsResp, error) {

	resp := &entity.WalletAddressAssetsResp{}

	err := c.client.Post(ctx, consts.PathWalletAddressAssets, req, resp)
	if err != nil {
		c.client.Logger().Errorf(ctx, "WalletAddressAssets err:%v", err)
		return nil, err
	}

	return resp, nil
}
