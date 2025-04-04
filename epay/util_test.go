package epay

import (
	"encoding/pem"
	"log"
	"sort"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestGetSignContent(t *testing.T) {
	asserts := assert.New(t)
	{
		requestParams := map[string]string{
			"pid":          "1",
			"type":         "1",
			"out_trade_no": "1",
			"notify_url":   "1",
			"name":         "1",
			"money":        "1",
			"device":       "1",
			"sign_type":    "MD5",
			"sign":         "",
		}

		expected := "device=1&money=1&name=1&notify_url=1&out_trade_no=1&pid=1&type=1"
		actual := GetSignContent(requestParams)
		asserts.Equal(expected, actual)
	}
}

func TestMD5String(t *testing.T) {
	asserts := assert.New(t)
	{
		urlString := "device=devicev&money=moneyv"
		key := "1234567"
		expected := "3854cc9f022e0fb821bd2e002260245d"
		asserts.EqualValues(expected, MD5String(urlString, key))
	}
}

func TestMapToStruct(t *testing.T) {
	{
		var verifyInfo VerifyRes
		mapstructure.Decode(map[string]string{
			"pid":  "pidv",
			"type": "typev",
			//
			"out_trade_no": "out_trade_nov",
			//
			"notify_url": "notify_urlv",
			//
			"name": "namev",
			//
			"money": "moneyv",
			//
			"device": "devicev",
		}, &verifyInfo)
		log.Println(verifyInfo)
	}
}

// go test -v ./epay -run TestRSA2SignAndVerify
func TestRSA2SignAndVerify(t *testing.T) {
	// 正确格式的私钥
	privateKey := ``

	// 正确格式的公钥
	publicKey := ``

	// 测试1: 直接测试RSA签名和验证函数
	testStr := "test_string_for_rsa_sign"
	t.Log("测试字符串:", testStr)

	// 签名测试
	sign, err := RSASign(testStr, privateKey)
	if err != nil {
		t.Fatalf("RSA签名失败: %v", err)
	}
	t.Logf("签名结果: %s", sign)

	// 验证测试
	verified, err := RSAVerify(testStr, sign, publicKey)
	if err != nil {
		t.Fatalf("RSA验证失败: %v", err)
	}
	t.Logf("验证结果: %v", verified)

	if !verified {
		t.Fatal("签名验证未通过")
	}

	// 测试2: 测试GenerateParams函数
	t.Log("\n测试GenerateParams函数")

	// 创建测试参数
	params := map[string]string{
		"out_trade_no": "test123456",
		"money":        "100.00",
		"name":         "测试商品",
		"notify_url":   "http://example.com/notify",
		"return_url":   "http://example.com/return",
	}

	// 保存原始参数用于调试
	originalParams := make(map[string]string)
	for k, v := range params {
		originalParams[k] = v
	}

	// 生成签名参数
	GenerateParams(params, privateKey, SignTypeRSA)

	// 检查签名是否已添加
	if _, ok := params["sign"]; !ok {
		t.Fatal("签名参数缺失")
	}

	if _, ok := params["sign_type"]; !ok {
		t.Fatal("签名类型参数缺失")
	}

	t.Logf("生成的签名: %s", params["sign"])
	t.Logf("参数列表: %v", params)

	// 验证生成的签名
	// 1. 重建待签名字符串
	var keys []string
	for k := range originalParams {
		if k != "sign" && k != "sign_type" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var signParts []string
	for _, k := range keys {
		if originalParams[k] != "" {
			signParts = append(signParts, k+"="+originalParams[k])
		}
	}
	signStr := strings.Join(signParts, "&")
	t.Logf("重建的待签名字符串: %s", signStr)

	// 2. 验证签名
	verifyResult, err := RSAVerify(signStr, params["sign"], publicKey)
	if err != nil {
		t.Fatalf("验证GenerateParams签名失败: %v", err)
	}

	if !verifyResult {
		t.Fatal("GenerateParams签名验证未通过")
	} else {
		t.Log("GenerateParams签名验证通过")
	}

	// 打印RSASign函数内部过程
	t.Log("\n调试RSASign函数")
	// 解码私钥
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		t.Fatal("私钥PEM解码失败")
	} else {
		t.Log("私钥PEM解码成功")
	}

	// 此处可以进一步测试解析PKCS1或PKCS8私钥的过程
	// ...

	// 如果测试通过
	t.Log("RSA2签名测试通过")
}
