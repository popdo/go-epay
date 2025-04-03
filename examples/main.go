package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/popdo/go-epay/epay"
	"github.com/samber/lo"
)

func main() {
	baseUrl := "http://localhost:8080"
	client, err := epay.NewClient(&epay.Config{
		PartnerID: "1000",
		Key:       "KEY",
		PublicKey: "PLATFORM_PUBLIC_KEY", // 平台公钥（用于验签）
	}, baseUrl)

	if err != nil {
		log.Panicln(err)
	}
	notify, _ := url.Parse(baseUrl + "/verify")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		url, params, err := client.Purchase(&epay.PurchaseArgs{
			Type:           "wxpay",
			ServiceTradeNo: "8412317576584121",
			Name:           "test",
			Money:          "0.01",
			Device:         epay.PC,
			NotifyUrl:      notify,
			ReturnUrl:      notify,
		})
		if err != nil {
			log.Println(err)
			return
		}

		html := "<form id='alipaysubmit' name='alipaysubmit' action='" + url + "' method='POST'>"
		for key, value := range params {
			html += "<input type='hidden' name='" + key + "' value='" + value + "'/>"
		}
		html += "<input type='submit'>POST</form>"

		writer.Header().Set("Content-Type", "text/html")
		writer.Write([]byte(html))
	})

	// 新增API支付示例
	mux.HandleFunc("/api_pay", func(writer http.ResponseWriter, request *http.Request) {
		clientIP := request.RemoteAddr
		if ip := request.Header.Get("X-Real-IP"); ip != "" {
			clientIP = ip
		} else if ip = request.Header.Get("X-Forwarded-For"); ip != "" {
			clientIP = strings.Split(ip, ",")[0]
		}

		result, err := client.ApiPurchase(&epay.ApiPurchaseArgs{
			Type:           "wxpay",
			ServiceTradeNo: "API" + time.Now().Format("20060102150405"),
			Name:           "API支付测试",
			Money:          "0.01",
			ClientIP:       clientIP,
			Device:         epay.PC,
			NotifyUrl:      notify,
			ReturnUrl:      notify,
		})

		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(result)

		// 如果需要显示二维码或跳转
		if result.Code == 1 {
			if result.QRCode != "" {
				log.Println("生成二维码:", result.QRCode)
			} else if result.PayURL != "" {
				log.Println("跳转支付URL:", result.PayURL)
			} else if result.URLScheme != "" {
				log.Println("小程序URL:", result.URLScheme)
			}
		}
	})

	mux.HandleFunc("/verify", func(writer http.ResponseWriter, request *http.Request) {
		params := lo.Reduce(lo.Keys(request.URL.Query()), func(r map[string]string, t string, i int) map[string]string {
			r[t] = request.URL.Query().Get(t)
			return r
		}, map[string]string{})

		verifyInfo, err := client.Verify(params)
		if err == nil && verifyInfo.VerifyStatus {
			writer.Write([]byte("success"))
		} else {
			writer.Write([]byte("fail"))
		}

		if verifyInfo.TradeStatus == epay.StatusTradeSuccess {
			log.Println(verifyInfo)
		}
	})
	http.ListenAndServe(":8080", mux)
}
