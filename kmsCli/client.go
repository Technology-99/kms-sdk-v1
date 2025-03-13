package kmsCli

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Technology-99/kms-sdk-v1/kmsConfig"
	"github.com/Technology-99/third_party/middleware"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	KmsClientStatusNoInit = iota + 1
	KmsClientStatusInit
	KmsClientStatusReady
)

type KmsClient struct {
	*http.Client
	Config  *kmsConfig.KmsConfig
	Context context.Context
	Status  int
}

func NewKmsClient(ctx context.Context, conf *kmsConfig.KmsConfig) *KmsClient {
	httpClient := &http.Client{
		Timeout: kmsConfig.DefaultTimeout * time.Millisecond,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return &KmsClient{
		Config:  conf,
		Client:  httpClient,
		Context: ctx,
		Status:  KmsClientStatusReady,
	}
}

func (cli *KmsClient) WithContext(ctx context.Context) *KmsClient {
	cli.Context = ctx
	return cli
}

// note: 添加将requestID继承到下个服务的能力
func (cli *KmsClient) WithRequestId(requestId string) *KmsClient {
	cli.Context = context.WithValue(cli.Context, middleware.CtxRequestID, requestId)
	return cli
}

func (cli *KmsClient) WithTimeout(timeout time.Duration) *KmsClient {
	cli.Client.Timeout = timeout
	return cli
}

func (cli *KmsClient) EasyNewRequest(ctx context.Context, relativePath string, method string, sendBody interface{}) func() ([]byte, error) {
	apiUrl := fmt.Sprintf("%s://%s/kms/%s/apis%s?token=%s", cli.Config.Protocol, cli.Config.Host, cli.Config.Version, relativePath, cli.Config.Token)
	if cli.Context.Value(middleware.CtxRequestID) != nil {
		logx.Infof("requestID: %s, EasyNewRequest url: %s", cli.Context.Value(middleware.CtxRequestID), apiUrl)
	} else {
		logx.Infof("EasyNewRequest url: %s", apiUrl)
	}
	return cli.NewRequest(ctx, apiUrl, method, cli.GenHeaders(), sendBody)
}

func (cli *KmsClient) NewRequest(
	ctx context.Context, // 新增 context 参数
	url string, // URL
	method string, // HTTP 方法
	headers *map[string]string, // 请求头
	sendBody interface{}) func() ([]byte, error) { // 返回闭包函数

	var (
		body []byte
		err  error
	)

	// 创建一个 channel 来控制请求完成或超时
	c := make(chan struct{})
	go func() {
		defer close(c) // 保证 goroutine 退出时关闭 channel

		sendBodyJson := ""

		if sendBody != nil {
			// 将发送体序列化为 JSON
			sendBodyBt, marshalErr := json.Marshal(sendBody)
			if marshalErr != nil {
				err = marshalErr
				return
			}
			sendBodyJson = string(sendBodyBt)
		}

		logx.Infof("request sendBodyJson: %v", sendBodyJson)

		// 使用 context 控制请求
		req, reqErr := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer([]byte(sendBodyJson)))
		if reqErr != nil {
			err = reqErr
			return
		}

		// 设置请求头
		if headers != nil {
			for k, v := range *headers {
				req.Header.Set(k, v)
			}
		}

		// 发送请求
		var res *http.Response
		res, err = cli.Client.Do(req)
		if err != nil {
			return
		}
		defer res.Body.Close()

		// 读取响应体
		body, err = ioutil.ReadAll(res.Body)
	}()

	return func() ([]byte, error) {
		select {
		case <-c: // 请求完成
			return body, err
		case <-ctx.Done(): // 请求超时或取消
			return nil, ctx.Err()
		}
	}
}

func (cli *KmsClient) GenHeaders() *map[string]string {
	// note: 先处理请求头
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	return &headers
}
