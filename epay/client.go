package epay

import "net/url"

var _ Service = (*Client)(nil)

// 易支付API
type Service interface {
	// 创建订单
	CreateOrder(args *CreateOrderArgs) (string, map[string]string, error)
	// API创建订单
	ApiCreateOrder(args *ApiCreateOrderArgs) (*ApiCreateOrderRes, error)
	// 查询订单
	QueryOrder(tradeNo, outTradeNo string) (*ApiOrderQueryRes, error)
	// Verify 验证回调参数是否符合签名
	Verify(params map[string]string) (*VerifyRes, error)
}

// 创建一个新的易支付客户端
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
