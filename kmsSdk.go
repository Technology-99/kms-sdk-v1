package kmsSdk

import (
	"context"
	"encoding/json"
	"github.com/Technology-99/kms-sdk-v1/kmsCli"
	"github.com/Technology-99/kms-sdk-v1/kmsConfig"
	"github.com/Technology-99/kms-sdk-v1/kmsTypes"
	"github.com/Technology-99/third_party/qxCrypto"
	"github.com/Technology-99/third_party/response"
	"log"
	"net/http"
)

const (
	KmsParserStatusNoInit = iota + 1
	KmsParserStatusInit
	KmsParserStatusReady
)

type (
	KmsParser interface {
		Status() int
		WithContext(ctx context.Context) KmsParser
		CreateAesKey() (kmsTypes.CreateAesKeyResp, error)
		BatchEncrypt(req *kmsTypes.BatchEncryptDataReq) (kmsTypes.BatchEncryptDataResp, error)
		BatchDecrypt(req *kmsTypes.BatchDecryptDataReq) (kmsTypes.BatchDecryptDataResp, error)
		Encrypt(req *kmsTypes.EncryptDataReq) (kmsTypes.EncryptDataResp, error)
		Decrypt(req *kmsTypes.DecryptDataReq) (kmsTypes.DecryptDataResp, error)
		DecryptUnAutoDecode(req *kmsTypes.DecryptDataReq) (kmsTypes.DecryptDataResp, error)
		CompareData(req *kmsTypes.CompareAesKeyEncryptDataReq) (kmsTypes.CompareAesKeyEncryptDataResp, error)
	}
	defaultKmsParser struct {
		cli    *kmsCli.KmsClient
		status int
	}
)

func EasyKmsParser(Token, AesKey, AesIv string) KmsParser {
	defaultConfig := kmsConfig.DefaultKmsConfig().WithToken(Token).WithTransferAesKey(AesKey).WithTransferAesIv(AesIv)
	ctx := context.Background()
	return &defaultKmsParser{
		cli:    kmsCli.NewKmsClient(ctx, defaultConfig),
		status: KmsParserStatusNoInit,
	}
}

func NewKmsParser(c *kmsConfig.KmsConfig) KmsParser {
	ctx := context.Background()
	return &defaultKmsParser{
		cli:    kmsCli.NewKmsClient(ctx, c),
		status: KmsParserStatusNoInit,
	}
}

func (m *defaultKmsParser) CompareData(req *kmsTypes.CompareAesKeyEncryptDataReq) (kmsTypes.CompareAesKeyEncryptDataResp, error) {
	result := kmsTypes.CompareAesKeyEncryptDataResp{}
	reqFn := m.cli.EasyNewRequest(context.Background(), "/aesGcm/compareAesKeyEncryptData", http.MethodPost, req)
	res, err := reqFn()
	if err != nil {
		log.Printf("BatchEncrypt request error: %v", err)
		return result, err
	}
	log.Printf("BatchEncrypt response: %s", res)
	_ = json.Unmarshal(res, &result)
	if result.Code != response.SUCCESS {
		log.Printf("kms sdk errlog: Encrypt fail: %v", result)
		return result, err
	}
	return result, nil
}

func (m *defaultKmsParser) BatchEncrypt(req *kmsTypes.BatchEncryptDataReq) (kmsTypes.BatchEncryptDataResp, error) {
	result := kmsTypes.BatchEncryptDataResp{}
	reqFn := m.cli.EasyNewRequest(context.Background(), "/aesGcm/batchEncrypt", http.MethodPost, req)
	res, err := reqFn()
	if err != nil {
		log.Printf("BatchEncrypt request error: %v", err)
		return result, err
	}
	log.Printf("BatchEncrypt response: %s", res)
	_ = json.Unmarshal(res, &result)
	if result.Code != response.SUCCESS {
		log.Printf("kms sdk errlog: Encrypt fail: %v", result)
		return result, err
	}
	return result, nil
}

func (m *defaultKmsParser) BatchDecrypt(req *kmsTypes.BatchDecryptDataReq) (kmsTypes.BatchDecryptDataResp, error) {
	result := kmsTypes.BatchDecryptDataResp{}
	reqFn := m.cli.EasyNewRequest(context.Background(), "/aesGcm/batchDecrypt", http.MethodPost, req)
	res, err := reqFn()
	if err != nil {
		log.Printf("batchDecrypt request error: %v", err)
		return result, err
	}
	_ = json.Unmarshal(res, &result)
	if result.Code != response.SUCCESS {
		log.Printf("batchDecrypt failed: %v", result)
		return result, err
	}

	for k, v := range result.Data.Result {
		aesDecryptResult, err := qxCrypto.AESDecryptByGCM(v.Data, m.cli.Config.TransferAesKey, m.cli.Config.TransferAesIv)
		if err != nil {
			log.Printf("aes decrypt fail: %v", err)
			return result, err
		}
		result.Data.Result[k] = kmsTypes.ModelBatchItem{
			Key:  v.Key,
			Data: string(aesDecryptResult),
		}
	}
	return result, nil
}

func (m *defaultKmsParser) Encrypt(req *kmsTypes.EncryptDataReq) (kmsTypes.EncryptDataResp, error) {
	result := kmsTypes.EncryptDataResp{}
	reqFn := m.cli.EasyNewRequest(context.Background(), "/aesGcm/encrypt", http.MethodPost, req)
	res, err := reqFn()
	if err != nil {
		log.Printf("Encrypt request error: %v", err)
		return result, err
	}
	log.Printf("Encrypt response: %s", res)
	_ = json.Unmarshal(res, &result)
	if result.Code != response.SUCCESS {
		log.Printf("kms sdk errlog: Encrypt fail: %v", result)
		return result, err
	}
	return result, nil
}

func (m *defaultKmsParser) Decrypt(req *kmsTypes.DecryptDataReq) (kmsTypes.DecryptDataResp, error) {
	result := kmsTypes.DecryptDataResp{}
	reqFn := m.cli.EasyNewRequest(context.Background(), "/aesGcm/decrypt", http.MethodPost, req)
	res, err := reqFn()
	if err != nil {
		log.Printf("Decrypt request error: %v", err)
		return result, err
	}
	_ = json.Unmarshal(res, &result)
	if result.Code != response.SUCCESS {
		log.Printf("kms sdk errlog: Decrypt fail: %v", result)
		return result, err
	}

	aesDecryptResult, err := qxCrypto.AESDecryptByGCM(result.Data.Data, m.cli.Config.TransferAesKey, m.cli.Config.TransferAesIv)
	if err != nil {
		log.Printf("aes decrypt fail: %v", err)
		return result, err
	}
	result.Data.Data = string(aesDecryptResult)
	return result, nil
}

func (m *defaultKmsParser) DecryptUnAutoDecode(req *kmsTypes.DecryptDataReq) (kmsTypes.DecryptDataResp, error) {
	result := kmsTypes.DecryptDataResp{}
	reqFn := m.cli.EasyNewRequest(context.Background(), "/aesGcm/decrypt", http.MethodPost, req)
	res, err := reqFn()
	if err != nil {
		log.Printf("Decrypt request error: %v", err)
		return result, err
	}
	_ = json.Unmarshal(res, &result)
	if result.Code != response.SUCCESS {
		log.Printf("kms sdk errlog: Decrypt fail: %v", result)
		return result, err
	}
	return result, nil
}

func (m *defaultKmsParser) CreateAesKey() (kmsTypes.CreateAesKeyResp, error) {
	result := kmsTypes.CreateAesKeyResp{}
	reqFn := m.cli.EasyNewRequest(context.Background(), "/aesGcm/createAesKey", http.MethodPost, nil)
	res, err := reqFn()
	if err != nil {
		log.Printf("CreateAesKey request error: %v", err)
		return result, err
	}
	_ = json.Unmarshal(res, &result)
	if result.Code != response.SUCCESS {
		log.Printf("kms sdk errlog: CreateAesKey fail: %v", result)
		return result, err
	}
	return result, nil
}

func (m *defaultKmsParser) WithContext(ctx context.Context) KmsParser {
	m.cli.WithContext(ctx)
	return m
}

func (m *defaultKmsParser) Status() int {
	return m.status
}
