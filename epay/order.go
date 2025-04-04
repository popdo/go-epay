package epay

// 创建订单
func (c *Client) CreateOrder(args *CreateOrderArgs) (string, map[string]string, error) {
	if c.Config.PublicKey != "" {
		return c.V2CreateOrder(args)
	}
	return c.V1CreateOrder(args)
}

// API创建订单
func (c *Client) ApiCreateOrder(args *ApiCreateOrderArgs) (*ApiCreateOrderRes, error) {
	if c.Config.PublicKey != "" {
		return c.V2ApiCreateOrder(args)
	}
	return c.V1ApiCreateOrder(args)
}

// 单个订单查询
func (c *Client) QueryOrder(tradeNo, outTradeNo string) (*ApiOrderQueryRes, error) {
	if c.Config.PublicKey != "" {
		return c.V2QueryOrder(tradeNo, outTradeNo)
	}
	return c.V1QueryOrder(tradeNo, outTradeNo)
}
