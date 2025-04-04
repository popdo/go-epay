package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/popdo/go-epay/epay"
)

func main() {
	var err error
	// 运行签名测试 ==================
	err = runRSASignTest()
	if err != nil {
		log.Fatalf("测试失败: %v", err)
	}
	log.Println("测试完成")

	// 测试V1订单创建 ==================
	err = testCreateV1Order()
	if err != nil {
		log.Fatalf("创建订单失败: %v", err)
	}
	log.Println("订单创建测试完成")

	// 测试V2订单创建 ==================
	err = testCreateV2Order()
	if err != nil {
		log.Fatalf("创建订单失败: %v", err)
	}
	log.Println("订单创建测试完成")
}

func runRSASignTest() error {
	// 正确格式的商户私钥
	privateKey := ``

	// 正确格式的商户公钥
	publicKey := ``

	// 第1步: 测试简单字符串签名
	fmt.Println("===== 测试简单字符串签名验证 =====")
	testStr := "test_string_for_rsa_sign"
	sign, err := epay.RSASign(testStr, privateKey)
	if err != nil {
		return fmt.Errorf("RSA签名失败: %v", err)
	}
	fmt.Printf("签名结果: %s\n", sign)

	verified, err := epay.RSAVerify(testStr, sign, publicKey)
	if err != nil {
		return fmt.Errorf("RSA验证失败: %v", err)
	}
	fmt.Printf("验证结果: %v\n\n", verified)

	baseUrl := "http://" // 修改为实际支付网关地址
	notifyUrl, _ := url.Parse(baseUrl + "/api/user/epay/notify")
	returnUrl, _ := url.Parse(baseUrl + "/panel/topup")
	// 第2步: 测试实际支付参数签名和验证
	fmt.Println("===== 测试支付参数签名验证 =====")
	requestParams := map[string]string{
		"pid":          "1000", // 商户ID
		"method":       "web",
		"type":         "wxpay",
		"out_trade_no": "8412317576584121",
		"notify_url":   notifyUrl.String(),
		"return_url":   returnUrl.String(),
		"name":         "测试商品",
		"money":        "1.00",
		"clientip":     "127.0.0.1",
		"timestamp":    fmt.Sprintf("%d", time.Now().Unix()),
	}

	// 生成带签名的参数
	signParams := epay.GenerateParams(requestParams, privateKey, epay.SignTypeRSA)
	fmt.Println("签名后参数:", signParams)
	fmt.Printf("签名值: %s\n", signParams["sign"])

	// 第3步: 使用Client验证签名
	fmt.Println("\n===== 使用Client验证签名 =====")

	client, err := epay.NewClient(&epay.Config{
		PartnerID: "1000",
		Key:       privateKey,
		PublicKey: publicKey, // 平台公钥（用于验签）
	}, baseUrl)

	if err != nil {
		return fmt.Errorf("创建客户端失败: %v", err)
	}

	verifyInfo, err := client.Verify(signParams)
	if err == nil && verifyInfo.VerifyStatus {
		fmt.Println("验签成功")
	} else {
		fmt.Println("验签失败:", err)
	}
	fmt.Printf("验签结果: %+v\n", verifyInfo)

	return nil
}
func testCreateV2Order() error {
	fmt.Println("\n===== 测试V2创建订单 =====")

	// 正确格式的商户私钥
	privateKey := ``

	// 正确格式的商户公钥
	publicKey := ``

	// 创建客户端
	baseUrl := "http://vpay.gpt.ge" // 修改为实际支付网关地址
	client, err := epay.NewClient(&epay.Config{
		PartnerID: "1000", // 修改为你的商户ID
		Key:       privateKey,
		PublicKey: publicKey, // v2接口需要提供平台公钥
	}, baseUrl)

	if err != nil {
		return fmt.Errorf("创建客户端失败: %v", err)
	}

	// 创建通知URL
	notifyUrl, _ := url.Parse(baseUrl + "/api/user/epay/notify")
	returnUrl, _ := url.Parse(baseUrl + "/panel/topup")

	// 测试API支付创建订单
	orderNo := fmt.Sprintf("TEST%d", time.Now().Unix())

	orderArgs := &epay.ApiCreateOrderArgs{
		Device:     epay.PC,
		Type:       "wxpay", // 支付方式
		OutTradeNo: orderNo, // 商户订单号
		NotifyURL:  notifyUrl,
		ReturnURL:  returnUrl,
		Name:       "V2测试商品",
		Money:      "2.5", // 测试金额
		ClientIP:   "127.0.0.1",
	}
	if publicKey != "" {
		orderArgs.Method = epay.MethodWeb // 会根据device判断，自动 返回跳转url/二维码/小程序跳转url等
	}
	result, err := client.ApiCreateOrder(orderArgs)

	if err != nil {
		return fmt.Errorf("创建订单失败: %v", err)
	}

	fmt.Printf("API创建订单结果:\n- 状态码: %d\n- 消息: %s\n", result.Code, result.Message)

	// V2接口成功返回码是0
	if result.Code == 0 {
		fmt.Println("创建订单成功!")
		if result.PayURL != "" {
			fmt.Println("支付URL:", result.PayURL)
		}
		if result.QRCode != "" {
			fmt.Println("二维码链接:", result.QRCode)
		}
		if result.PayInfo != "" {
			fmt.Println("支付参数:", result.PayInfo)
		}
	}

	return nil
}
func testCreateV1Order() error {
	fmt.Println("\n===== 测试V1创建订单 =====")

	// 正确格式的商户私钥
	privateKey := ``

	// 创建客户端
	baseUrl := "http://vpay.gpt.ge" // 修改为实际支付网关地址
	client, err := epay.NewClient(&epay.Config{
		PartnerID: "1000", // 修改为你的商户ID
		Key:       privateKey,
	}, baseUrl)

	if err != nil {
		return fmt.Errorf("创建客户端失败: %v", err)
	}

	// 创建通知URL
	notifyUrl, _ := url.Parse(baseUrl + "/api/user/epay/notify")
	returnUrl, _ := url.Parse(baseUrl + "/panel/topup")

	// 测试API支付创建订单
	orderNo := fmt.Sprintf("TEST%d", time.Now().Unix())

	orderArgs := &epay.ApiCreateOrderArgs{
		Device:     epay.PC,
		Type:       "wxpay", // 支付方式
		OutTradeNo: orderNo, // 商户订单号
		NotifyURL:  notifyUrl,
		ReturnURL:  returnUrl,
		Name:       "V1测试商品",
		Money:      "2.5", // 测试金额
		ClientIP:   "127.0.0.1",
	}

	result, err := client.ApiCreateOrder(orderArgs)

	if err != nil {
		return fmt.Errorf("创建订单失败: %v", err)
	}

	fmt.Printf("API创建订单结果:\n- 状态码: %d\n- 消息: %s\n", result.Code, result.Message)

	// V2接口成功返回码是0
	if result.Code == 0 {
		fmt.Println("创建订单成功!")
		if result.PayURL != "" {
			fmt.Println("支付URL:", result.PayURL)
		}
		if result.QRCode != "" {
			fmt.Println("二维码链接:", result.QRCode)
		}
		if result.PayInfo != "" {
			fmt.Println("支付参数:", result.PayInfo)
		}
	}

	return nil
}
