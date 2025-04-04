package epay

import "github.com/mitchellh/mapstructure"

// Verify 验证回调参数是否符合签名
// 验签流程：
// 1. 获取待签名字符串（过滤参数 -> 排序 -> 生成URL字符串）
// 2. 根据签名类型选择验证方法：
//   - SignTypeRSA: 使用平台公钥进行RSA验签（SHA256WithRSA）
//   - SignTypeMD5: 使用MD5密钥进行验证
//
// 注意：
// - 商户私钥(Key)用于请求时签名
// - 平台公钥(PublicKey)用于验证平台返回数据的签名
func (c *Client) Verify(params map[string]string) (*VerifyRes, error) {
	sign := params["sign"]
	signType := params["sign_type"]
	var verifyRes VerifyRes

	// 从 map 映射到 struct 上
	err := mapstructure.Decode(params, &verifyRes)
	if err != nil {
		return nil, err
	}

	// 准备验证签名
	urlString := GetSignContent(params)

	// 根据签名类型和是否提供PublicKey来决定验证方式
	if signType == SignTypeRSA && c.Config.PublicKey != "" {
		verified, err := RSAVerify(urlString, sign, c.Config.PublicKey)
		if err != nil {
			return nil, err
		}
		verifyRes.VerifyStatus = verified
	} else {
		// 默认MD5验证
		verifyRes.VerifyStatus = sign == MD5String(urlString, c.Config.Key)
	}
	return &verifyRes, nil
}
