package epay

import "net/url"

var _ Service = (*Client)(nil)

// Service 易支付API
type Service interface {
	// Purchase 生成支付链接和参数
	Purchase(args *PurchaseArgs) (string, map[string]string, error)
	// 新增 API支付方法
	ApiPurchase(args *ApiPurchaseArgs) (*ApiPurchaseRes, error)
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
