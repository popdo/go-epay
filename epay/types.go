package epay

import "net/url"

const StatusTradeSuccess = "TRADE_SUCCESS"

// 支付接口类型（预留外部使用）
const (
	MethodWeb    = "web"    // 通用网页支付
	MethodWap    = "wap"    // H5支付
	MethodQrcode = "qrcode" // 扫码支付
	MethodJsapi  = "jsapi"  // JSAPI支付
	MethodMinipg = "minipg" // 小程序支付
)

// 支付发起类型（预留外部使用）
const (
	PayTypeJump     = "jump"     // 跳转支付
	PayTypeQrcode   = "qrcode"   // 二维码支付
	PayTypeJsapi    = "jsapi"    // JSAPI支付
	PayTypeScan     = "scan"     // 扫码支付结果
	PayTypeWxplugin = "wxplugin" // 微信收银台支付
	PayTypeWxapp    = "wxapp"    // 微信小程序跳转支付
)

type DeviceType string

var (
	PC     DeviceType = "pc"     // PC PC端
	MOBILE DeviceType = "mobile" // MOBILE 移动端
	QQ     DeviceType = "qq"     // QQ QQ内置浏览器
	WECHAT DeviceType = "wechat" // WECHAT 微信内置浏览器
	ALIPAY DeviceType = "alipay" // ALIPAY 支付宝
)

type Config struct {
	PartnerID string // 商户ID
	Key       string // MD5密钥或RSA私钥
	PublicKey string // 平台公钥(用于验证签名)
}
type Client struct {
	Config  *Config
	BaseUrl *url.URL
}

type CreateOrderArgs struct {
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
	// 业务扩展参数
	Param string
}

// API支付请求参数
type ApiCreateOrderArgs struct {
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

	// 新增V2特有字段
	Method    string //  (V2必填)
	AuthCode  string // 扫码支付授权码 仅被扫支付需要传
	SubOpenID string // 用户OpenID 仅JSAPI支付需要传
	SubAppID  string // 公众号AppID 仅JSAPI支付需要传
}

// API支付响应
type ApiCreateOrderRes struct {
	// 返回状态码 v1是1成功，其他失败
	// 返回状态码 v2是0成功，其他失败
	Code int `json:"code"`
	// 返回信息
	Message string `json:"msg"`
	// 订单号
	TradeNo string `json:"trade_no"`
	// V1特有字段
	PayURL    string `json:"payurl,omitempty"`    // 支付跳转URL (三选一)
	QRCode    string `json:"qrcode,omitempty"`    // 二维码链接 (三选一)
	URLScheme string `json:"urlscheme,omitempty"` // 小程序跳转URL (三选一)

	// V2特有字段
	PayType   string `json:"pay_type,omitempty"`  // 发起支付类型
	PayInfo   string `json:"pay_info,omitempty"`  // 发起支付参数
	Timestamp string `json:"timestamp,omitempty"` // 时间戳
	Sign      string `json:"sign,omitempty"`      // 签名
	SignType  string `json:"sign_type,omitempty"` // 签名类型
}

// OrderQueryRes 查询订单响应
type ApiOrderQueryRes struct {
	// 返回状态码 1成功，其他失败
	Code int `json:"code"`
	// 返回信息
	Message string `json:"msg"`
	// 易支付订单号
	TradeNo string `json:"trade_no"`
	// 商户订单号
	ServiceTradeNo string `json:"out_trade_no"`
	// 第三方订单号
	ApiTradeNo string `json:"api_trade_no"`
	// 支付方式
	Type string `json:"type"`
	// 商户ID
	// PID int `json:"pid"`
	// 创建订单时间
	AddTime string `json:"addtime"`
	// 完成交易时间
	EndTime string `json:"endtime"`
	// 商品名称
	Name string `json:"name"`
	// 金额
	Money string `json:"money"`
	// 支付状态 1支付成功，0未支付
	Status int `json:"status"`
	// 业务扩展参数
	Param string `json:"param"`
	// 支付者账号
	Buyer string `json:"buyer"`

	// V2特有字段
	RefundMoney string `json:"refundmoney,omitempty"` // 已退款金额
	ClientIP    string `json:"clientip,omitempty"`    // 用户IP
	Timestamp   string `json:"timestamp,omitempty"`   // 时间戳
	Sign        string `json:"sign,omitempty"`        // 签名
	SignType    string `json:"sign_type,omitempty"`   // 签名类型
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
