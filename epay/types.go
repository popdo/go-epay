package epay

import "net/url"

var (
	PC     DeviceType = "pc"     // PC PC端
	MOBILE DeviceType = "mobile" // MOBILE 移动端
)

type Config struct {
	PartnerID string
	Key       string // MD5密钥或RSA私钥
	PublicKey string // 平台公钥(用于验证签名)
}
type Client struct {
	Config  *Config
	BaseUrl *url.URL
}

type DeviceType string

type PurchaseArgs struct {
	// 支付类型
	Type string
	// 商家订单号
	ServiceTradeNo string
	// 商品名称
	Name string
	// 金额
	Money string
	// 设备类型
	Device    DeviceType
	NotifyUrl *url.URL
	ReturnUrl *url.URL
}

// API支付请求参数
type ApiPurchaseArgs struct {
	// 支付类型
	Type string
	// 商家订单号
	ServiceTradeNo string
	// 商品名称
	Name string
	// 金额
	Money string
	// 客户端IP
	ClientIP string
	// 设备类型 (可选)
	Device DeviceType
	// 业务扩展参数 (可选)
	Param string
	// 异步通知地址
	NotifyUrl *url.URL
	// 跳转地址 (可选)
	ReturnUrl *url.URL
}

// API支付响应
type ApiPurchaseRes struct {
	// 返回状态码 1成功，其他失败
	Code int `json:"code"`
	// 返回信息
	Message string `json:"msg"`
	// 订单号
	TradeNo string `json:"trade_no"`
	// 支付跳转URL (三选一)
	PayURL string `json:"payurl"`
	// 二维码链接 (三选一)
	QRCode string `json:"qrcode"`
	// 小程序跳转URL (三选一)
	URLScheme string `json:"urlscheme"`
}

// VerifyRes 验证结果
type VerifyRes struct {
	// 支付类型
	Type string
	// 易支付订单号
	TradeNo string `mapstructure:"trade_no"`
	// 商家订单号
	ServiceTradeNo string `mapstructure:"out_trade_no"`
	// 商品名称
	Name string
	// 金额
	Money string
	// 订单支付状态
	TradeStatus string `mapstructure:"trade_status"`
	// 签名检验
	VerifyStatus bool `mapstructure:"-"`
}
