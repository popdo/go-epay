package epay

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/samber/lo"
)

const (
	// v1
	V1CreateUrl    = "/submit.php" // v1 跳转支付
	V1ApiCreateUrl = "/mapi.php"   // v1 API支付
	V1QueryUrl     = "/api.php"    // v1 查询订单
)

// v1 创建订单
func (c *Client) V1CreateOrder(args *CreateOrderArgs) (string, map[string]string, error) {
	// see https://payment.moe/doc.html
	requestParams := map[string]string{
		"pid":          c.Config.PartnerID,
		"type":         args.Type,
		"out_trade_no": args.ServiceTradeNo,
		"notify_url":   args.NotifyUrl.String(),
		"return_url":   args.ReturnUrl.String(),
		"name":         args.Name,
		"money":        args.Money,
		"sign":         "",
		"sign_type":    SignTypeMD5,
	}
	if args.Param != "" {
		requestParams["param"] = args.Param
	}

	u, err := url.Parse(c.BaseUrl.String())
	if err != nil {
		return "", nil, err
	}
	u.Path = path.Join(u.Path, V1CreateUrl)

	return u.String(), GenerateParams(requestParams, c.Config.Key), nil
}

// ApiPurchase API接口创建订单
func (c *Client) V1ApiCreateOrder(args *ApiCreateOrderArgs) (*ApiCreateOrderRes, error) {
	// 构建请求参数
	requestParams := map[string]string{
		"pid":          c.Config.PartnerID,
		"type":         args.Type,
		"out_trade_no": args.ServiceTradeNo,
		"notify_url":   args.NotifyUrl.String(),
		"name":         args.Name,
		"money":        args.Money,
		"clientip":     args.ClientIP,
		"sign":         "",
		"sign_type":    SignTypeMD5,
	}

	// 添加可选参数
	if args.Device != "" {
		requestParams["device"] = string(args.Device)
	}
	if args.Param != "" {
		requestParams["param"] = args.Param
	}
	if args.ReturnUrl != nil {
		requestParams["return_url"] = args.ReturnUrl.String()
	}

	// 生成签名
	signParams := GenerateParams(requestParams, c.Config.Key)

	// 构建API接口URL
	apiUrl, err := url.Parse(c.BaseUrl.String())
	if err != nil {
		return nil, err
	}
	apiUrl.Path = path.Join(apiUrl.Path, V1ApiCreateUrl)

	// 发送POST请求
	resp, err := http.PostForm(apiUrl.String(), url.Values(lo.MapValues(signParams, func(v string, _ string) []string {
		return []string{v}
	})))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析JSON响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ApiCreateOrderRes
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// QueryOrder 查询单个订单
func (c *Client) V1QueryOrder(tradeNo, outTradeNo string) (*ApiOrderQueryRes, error) {
	// 构建请求参数
	queryUrl, err := url.Parse(c.BaseUrl.String())
	if err != nil {
		return nil, err
	}
	queryUrl.Path = path.Join(queryUrl.Path, V1QueryUrl)

	// 设置查询参数
	query := queryUrl.Query()
	query.Add("act", "order")
	query.Add("pid", c.Config.PartnerID)
	query.Add("key", c.Config.Key) // 使用商户密钥

	// 至少需要传入一个订单号
	if tradeNo != "" {
		query.Add("trade_no", tradeNo)
	} else if outTradeNo != "" {
		query.Add("out_trade_no", outTradeNo)
	} else {
		return nil, errors.New("必须提供系统订单号或商户订单号")
	}

	queryUrl.RawQuery = query.Encode()

	// 发送GET请求
	resp, err := http.Get(queryUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析JSON响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ApiOrderQueryRes
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
