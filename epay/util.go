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
)

// 确定签名类型
const (
	SignTypeMD5 = "MD5"
	SignTypeRSA = "RSA"
)

// RSASign 使用SHA256WithRSA算法生成签名
func RSASign(data string, privateKeyContent string) (string, error) {
	// 添加PKCS#8格式头尾并格式化
	privateKey := "-----BEGIN RSA PRIVATE KEY-----\n"
	keyLen := len(privateKeyContent)
	for i := 0; i < keyLen; i += 64 {
		end := i + 64
		if end > keyLen {
			end = keyLen
		}
		privateKey += privateKeyContent[i:end] + "\n"
	}
	privateKey += "-----END RSA PRIVATE KEY-----"

	// 解析私钥
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("private key error")
	}

	// 使用PKCS#8格式解析
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// 尝试PKCS#1格式解析
		rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
		privateKeyInterface = rsaPrivateKey
	}

	rsaPrivateKey, ok := privateKeyInterface.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("private key type error")
	}

	// 使用SHA256WithRSA算法
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// RSAVerify 使用RSA公钥验证签名
func RSAVerify(urlString, sign, publicKeyContent string) (bool, error) {
	// 添加公钥头尾并格式化
	publicKey := "-----BEGIN PUBLIC KEY-----\n"
	keyLen := len(publicKeyContent)
	for i := 0; i < keyLen; i += 64 {
		end := i + 64
		if end > keyLen {
			end = keyLen
		}
		publicKey += publicKeyContent[i:end] + "\n"
	}
	publicKey += "-----END PUBLIC KEY-----"

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

// MD5String 生成 加盐(商户 key) MD5 字符串
func MD5String(urlString string, key string) string {
	digest := md5.Sum([]byte(urlString + key))
	return fmt.Sprintf("%x", digest)
}

// GenerateParams 生成加签参数
func GenerateParams(params map[string]string, key string, signType string) map[string]string {
	// 复制一份参数，避免修改原始数据
	newParams := make(map[string]string)
	for k, v := range params {
		newParams[k] = v
	}
	// 生成待签名字符串
	signContent := GetSignContent(newParams)
	// 根据签名类型生成签名
	var sign string
	var err error
	if signType == SignTypeRSA {
		sign, err = RSASign(signContent, key)
		if err != nil {
			return newParams
		}
	} else if signType == SignTypeMD5 {
		sign = MD5String(signContent, key)
	}

	newParams["sign"] = sign
	newParams["sign_type"] = signType

	return newParams

}

// GetSignContent 获取待签名字符串，与PHP端逻辑保持一致
func GetSignContent(params map[string]string) string {
	var keys []string
	for k, v := range params {
		// 跳过空值、sign和sign_type
		if v == "" || k == "sign" || k == "sign_type" {
			continue
		}
		keys = append(keys, k)
	}

	// 按键名字典序排序
	sort.Strings(keys)

	var signParts []string
	for _, k := range keys {
		signParts = append(signParts, k+"="+params[k])
	}

	return strings.Join(signParts, "&")
}
