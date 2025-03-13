package main

import (
	"context"
	"github.com/Technology-99/kms-sdk-v1/kmsParser"
	"github.com/Technology-99/kms-sdk-v1/kmsTypes"
	"github.com/Technology-99/third_party/qxCrypto"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
)

func main() {

	const (
		// note: 联系系统管理员获取
		Token = "8d50c4ad-0b92-4506-b892-ff93fd6690b5"
		// note: 联系系统管理员获取, base64模式的key和iv
		BaseKey = "JlYRyK6N8tXpkXlgygIJ/j9EUQojkf+MavcV6Pdpl94="
		BaseIv  = "/g6I/zs8dajEdZML"
	)

	msg := "kms密钥管理服务真滴强"

	newCtx := context.Background()
	newCtx = context.WithValue(newCtx, "requestId", "123456")
	// note: 创建一对aes密钥串
	kmsParse := kmsParser.EasyKmsParser(Token, BaseKey, BaseIv)
	createDataResult, err := kmsParse.WithContext(newCtx).CreateAesKey()
	if err != nil {
		log.Printf("Failed to create data: %v", err)
		return
	}
	logx.Infof("创建的密钥串: %s", createDataResult.Data.Key)

	// 加密
	encryptDataResult, err := kmsParse.WithContext(newCtx).Encrypt(&kmsTypes.EncryptDataReq{
		Data: msg,
		Key:  createDataResult.Data.Key,
	})
	if err != nil {
		log.Printf("Failed to encrypt data: %v", err)
		return
	}
	logx.Infof("加密结果: %s", encryptDataResult.Data.Data)

	// 解密
	decryptDataResult, err := kmsParse.WithContext(newCtx).Decrypt(&kmsTypes.DecryptDataReq{
		Data: encryptDataResult.Data.Data,
		Key:  createDataResult.Data.Key,
	})
	if err != nil {
		log.Printf("Failed to encrypt data: %v", err)
		return
	}
	logx.Infof("解密结果: %s", decryptDataResult.Data)

	decryptDataUnAutoDecodeResult, err := kmsParse.WithContext(newCtx).DecryptUnAutoDecode(&kmsTypes.DecryptDataReq{
		Data: encryptDataResult.Data.Data,
		Key:  createDataResult.Data.Key,
	})
	if err != nil {
		log.Printf("Failed to encrypt data: %v", err)
		return
	}
	logx.Infof("解密结果: %s", decryptDataUnAutoDecodeResult.Data)

	// note: 自行利用本地保存的aes密钥进行解密
	aesDecryptResult, err := qxCrypto.AESDecryptByGCM(decryptDataUnAutoDecodeResult.Data.Data, BaseKey, BaseIv)
	if err != nil {
		logx.Errorf("aes decrypt fail: %v", err)
		return
	}
	logx.Infof("aes解密结果: %s", aesDecryptResult)

}
