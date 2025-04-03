package epay

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/samber/lo"
)

// 确定签名类型
const (
	SignTypeMD5 = "MD5"
	SignTypeRSA = "RSA2"
)

// DetectSignType 根据密钥长度判断签名类型
func DetectSignType(key string) string {
	if len(key) > 1000 {
		return SignTypeRSA
	}
	return SignTypeMD5
}

// RSASign 使用RSA私钥进行SHA256WithRSA签名
func RSASign(urlString string, privateKey string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("private key error")
	}

	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(urlString))
	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// RSAVerify 使用RSA公钥验证签名
func RSAVerify(urlString, sign, publicKey string) (bool, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return false, errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("invalid public key type")
	}

	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256([]byte(urlString))
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signBytes)
	return err == nil, nil
}

// ParamsFilter 过滤参数，生成签名时需删除 “sign” 和 “sign_type” 参数
func ParamsFilter(params map[string]string) map[string]string {
	return lo.PickBy[string, string](params, func(key string, value string) bool {
		return !(key == "sign" || key == "sign_type" || value == "")
	})
}

// ParamsSort 对参数进行排序，返回排序后的 keys 和 values （go 中 map 为乱序）
func ParamsSort(params map[string]string) ([]string, []string) {
	keys := lo.Keys(params)
	sort.Strings(keys)

	values := lo.Map(keys, func(key string, i int) string {
		return params[key]
	})

	return keys, values
}

// CreateUrlString 生成待签名字符串, ["a", "b", "c"], ["d", "e", "f"] => "a=d&b=e&c=f"
func CreateUrlString(keys []string, values []string) string {
	urlString := ""
	for i, key := range keys {
		urlString += key + "=" + values[i] + "&"
	}
	// trim 掉最后的 &
	return strings.TrimSuffix(urlString, "&")
}

// MD5String 生成 加盐(商户 key) MD5 字符串
func MD5String(urlString string, key string) string {
	digest := md5.Sum([]byte(urlString + key))
	return fmt.Sprintf("%x", digest)
}

// GenerateParams 生成加签参数
//
//	func GenerateParams(params map[string]string, key string) map[string]string {
//		filtered := ParamsFilter(params)
//		keys, values := ParamsSort(filtered)
//		sign := MD5String(CreateUrlString(keys, values), key)
//		params["sign"] = sign
//		params["sign_type"] = "MD5"
//		return params
//	}
//
// GenerateParams 生成加签参数
func GenerateParams(params map[string]string, key string) map[string]string {
	filtered := ParamsFilter(params)
	keys, values := ParamsSort(filtered)
	urlString := CreateUrlString(keys, values)

	signType := DetectSignType(key)

	if signType == SignTypeRSA {
		// RSA签名
		sign, err := RSASign(urlString, key)
		if err != nil {
			// 如果RSA签名失败，回退到MD5
			sign = MD5String(urlString, key)
			signType = SignTypeMD5
		}
		params["sign"] = sign
	} else {
		// MD5签名
		sign := MD5String(urlString, key)
		params["sign"] = sign
	}
	return params
}
