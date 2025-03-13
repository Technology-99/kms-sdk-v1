package kmsTypes

type BootstrapReq struct {
}

type BootstrapResp struct {
	Code      int32  `json:"code"`
	Msg       string `json:"msg"`
	Path      string `json:"path"`
	RequestID string `json:"requestId"`
	Data      string `json:"data"`
}

type CreateAesKeyReq struct {
	Token string `form:"token"`
}

type CreateAesKeyResp struct {
	Code      int32                `json:"code"`
	Msg       string               `json:"msg"`
	Path      string               `json:"path"`
	RequestID string               `json:"requestId"`
	Data      CreateAesKeyRespData `json:"data"`
}

type CreateAesKeyRespData struct {
	Status string `json:"status"`
	Key    string `json:"key"`
}

type DecryptDataReq struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}

type DecryptDataResp struct {
	Code      int32               `json:"code"`
	Msg       string              `json:"msg"`
	Path      string              `json:"path"`
	RequestID string              `json:"requestId"`
	Data      ModelAesEncryptData `json:"data"`
}

type EncryptDataReq struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}

type EncryptDataResp struct {
	Code      int32               `json:"code"`
	Msg       string              `json:"msg"`
	Path      string              `json:"path"`
	RequestID string              `json:"requestId"`
	Data      ModelAesEncryptData `json:"data"`
}

type HealthzResp struct {
	Code      int32  `json:"code"`
	Msg       string `json:"msg"`
	Path      string `json:"path"`
	RequestID string `json:"requestId"`
	Data      string `json:"data"`
}

type ModelAesEncryptData struct {
	Status string `json:"status"`
	Key    string `json:"key"`
	Data   string `json:"data"`
}

type PingReq struct {
	Ping string `form:"ping"`
}

type PingResp struct {
	Code      int32  `json:"code"`
	Msg       string `json:"msg"`
	Path      string `json:"path"`
	RequestID string `json:"requestId"`
	Data      string `json:"data"`
}
