package kmsConfig

import (
	"kms-sdk-v1/kmsTypes"
	"time"
)

const (
	DefaultTimeout = 2000
)

type KmsConfig struct {
	Host           string `json:",default=kms.csvw88.com"`
	Protocol       string `json:",default=https"`
	Version        string `json:",default=v2"`
	Token          string
	TransferAesKey string
	TransferAesIv  string

	Debug    bool          `json:",default=false"`
	Timeout  time.Duration `json:",default=2000"`
	Deadline int64         `json:",default=5"`
}

func DefaultKmsConfig() (config *KmsConfig) {
	config = &KmsConfig{
		Host:     "kms.csvw88.com",
		Protocol: kmsTypes.KmsProtocolHttps,
		Version:  kmsTypes.KmsVersionV1,
		Debug:    false,
		Timeout:  DefaultTimeout,
		Deadline: 5,
	}
	return config
}

func NewKmsConfig(c KmsConfig) *KmsConfig {
	return &c
}

func (c *KmsConfig) WithDebug(Debug bool) *KmsConfig {
	c.Debug = Debug
	return c
}

func (c *KmsConfig) WithTimeout(Timeout time.Duration) *KmsConfig {
	c.Timeout = Timeout
	return c
}

func (c *KmsConfig) WithDeadline(Deadline int64) *KmsConfig {
	c.Deadline = Deadline
	return c
}

func (c *KmsConfig) WithToken(Token string) *KmsConfig {
	c.Token = Token
	return c
}

func (c *KmsConfig) WithTransferAesKey(TransferAesKey string) *KmsConfig {
	c.TransferAesKey = TransferAesKey
	return c
}

func (c *KmsConfig) WithTransferAesIv(TransferAesIv string) *KmsConfig {
	c.TransferAesIv = TransferAesIv
	return c
}
