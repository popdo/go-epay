package epay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/samber/lo"
)

const (
	V2CreateUrl    = "/api/pay/submit" // v2 跳转支付
	V2ApiCreateUrl = "/api/pay/create" // v2 API支付
	V2QueryUrl     = "/api/pay/query"  // v2 查询订单
)

// 创建订单
func (c *Client) V2CreateOrder(args *CreateOrderArgs) (string, map[string]string, error) {
	requestParams := map[string]string{
		"pid":          c.Config.PartnerID,
		"type":         args.Type,
		"out_trade_no": args.OutTradeNo,
		"notify_url":   args.NotifyUrl.String(),
		"return_url":   args.ReturnUrl.String(),
		"name":         args.Name,
		"money":        args.Money,
		"timestamp":    strconv.FormatInt(time.Now().Unix(), 10),
	}

	// 可选参数
	if args.Param != "" {
		requestParams["param"] = args.Param
	}

	u, err := url.Parse(c.BaseUrl.String())
	if err != nil {
		return "", nil, err
	}
	u.Path = path.Join(u.Path, V2CreateUrl)

	return u.String(), GenerateParams(requestParams, c.Config.Key, SignTypeRSA), nil
}

// API创建订单
func (c *Client) V2ApiCreateOrder(args *ApiCreateOrderArgs) (*ApiCreateOrderRes, error) {
	// 构建请求参数
	requestParams := map[string]string{
		"pid":          c.Config.PartnerID,
		"method":       args.Method, // 接口类型：web/wap/qrcode/jsapi/minipg
		"type":         args.Type,
		"out_trade_no": args.OutTradeNo,
		"notify_url":   args.NotifyURL.String(),
		"return_url":   args.ReturnURL.String(),
		"name":         args.Name,
		"money":        args.Money,
		"clientip":     args.ClientIP,
		"timestamp":    fmt.Sprintf("%d", time.Now().Unix()),
	}

	// 添加可选参数
	if args.Device != "" {
		requestParams["device"] = string(args.Device)
	}
	if args.Param != "" {
		requestParams["param"] = args.Param
	}
	if args.AuthCode != "" {
		requestParams["auth_code"] = args.AuthCode
	}
	if args.SubOpenID != "" {
		requestParams["sub_openid"] = args.SubOpenID
	}
	if args.SubAppID != "" {
		requestParams["sub_appid"] = args.SubAppID
	}

	// 生成签名
	signParams := GenerateParams(requestParams, c.Config.Key, SignTypeRSA)

	// 构建API接口URL
	apiUrl, err := url.Parse(c.BaseUrl.String())
	if err != nil {
		return nil, err
	}
	apiUrl.Path = path.Join(apiUrl.Path, V2ApiCreateUrl)

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

	// if result.Code != 0 {
	// 	// 记录响应的完整内容用于调试
	// 	responseBody := string(body)
	// 	log.Printf("创建订单失败，服务器响应: %s", responseBody)
	// }

	return &result, nil
}

// 查询单个订单
func (c *Client) V2QueryOrder(tradeNo, outTradeNo string) (*ApiOrderQueryRes, error) {
	// 构建请求参数
	requestParams := map[string]string{
		"pid":       c.Config.PartnerID,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	}

	// 至少需要传入一个订单号
	if tradeNo != "" {
		requestParams["trade_no"] = tradeNo
	} else if outTradeNo != "" {
		requestParams["out_trade_no"] = outTradeNo
	} else {
		return nil, errors.New("必须提供系统订单号或商户订单号")
	}

	// 生成签名
	signParams := GenerateParams(requestParams, c.Config.Key, SignTypeRSA)

	// 构建API接口URL
	apiUrl, err := url.Parse(c.BaseUrl.String())
	if err != nil {
		return nil, err
	}
	apiUrl.Path = path.Join(apiUrl.Path, V2QueryUrl)

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

	var result ApiOrderQueryRes
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
