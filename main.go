package main

import (
	"fmt"
	"github.com/Technology-99/third_party/qxCrypto"
	"log"
)

const (
	BaseKey = "jfPtq78ilnWMLhHMmGhRycZgkTESwA9PIkBPe3xckSc="
	BaseIv  = "c/QOlzfYuCzWrML1"
)

func main() {

	msg := "BoWjqfBVy5RZQxkstMV7L+4BMni03YH4kTQtYEeJj9CasziGdBei3/LPhX6ORRfEsJy/mGqVI1UgDxmYag=="

	// 解密

	decryptData, err := qxCrypto.AESDecryptByGCM(msg, BaseKey, BaseIv)
	if err != nil {
		log.Printf("Failed to decrypt data: %v", err)
		panic(err)
	}
	fmt.Printf("解密数据: %s", decryptData)
}
