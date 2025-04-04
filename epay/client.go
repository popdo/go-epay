package epay

import "net/url"

var _ Service = (*Client)(nil)

// Service 易支付API
type Service interface {
	// v1 创建订单
	V1CreateOrder(args *CreateOrderArgs) (string, map[string]string, error)
	// v1 API创建订单
	V1ApiCreateOrder(args *ApiCreateOrderArgs) (*ApiCreateOrderRes, error)
	// v1 查询订单
	V1QueryOrder(tradeNo, outTradeNo string) (*ApiOrderQueryRes, error)
	// Verify 验证回调参数是否符合签名
	Verify(params map[string]string) (*VerifyRes, error)
}

// NewClient 创建一个新的易支付客户端
func NewClient(config *Config, baseUrl string) (*Client, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	return &Client{
		Config:  config,
		BaseUrl: u,
	}, nil
}
