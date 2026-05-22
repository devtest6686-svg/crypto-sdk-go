package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lbe-io/crypto-sdk-go/sdk/logger"
	"github.com/lbe-io/crypto-sdk-go/sdk/signx"
)

// SDKConfig 用于存储SDK初始化参数
type SDKConfig struct {
	ServerHost string
	MerchantID string
	APIKey     string
	APISecret  string
}

// CryptoSDK 是SDK的主结构体
// 封装了所有对外暴露的方法
type CryptoClient struct {
	config     SDKConfig
	options    *clientOptions
	httpClient *http.Client
}

// NewCryptoSDK 初始化SDK
// 传入 server_host, merchant_id, apikey, apisecret
func NewClient(config SDKConfig, options ...ClientOptions) (*CryptoClient, error) {
	if config.ServerHost == "" || config.MerchantID == "" || config.APIKey == "" || config.APISecret == "" {
		return nil, fmt.Errorf("所有参数均不能为空")
	}

	ops := defaultClientOptions()
	for _, v := range options {
		v(ops)
	}

	client := &http.Client{
		Timeout: ops.timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
		},
	}
	return &CryptoClient{config: config, options: ops, httpClient: client}, nil
}

// 通用请求方法
func (sdk *CryptoClient) doRequest(ctx context.Context, method, path string, reqBody interface{}, respBody interface{}) error {
	url := fmt.Sprintf("%s%s", sdk.config.ServerHost, path)

	reqBy, err := json.Marshal(reqBody)
	if err != nil {
		sdk.options.logger.Errorf(ctx, "doRequest Marshal err:%v", err)
		return err
	}

	// 参数签名
	sg := signx.NewSignV1(sdk.options.logger)
	vd, err := sg.CreateSign(ctx, sdk.config.APISecret, url, sdk.config.APIKey, sdk.config.MerchantID, method, string(reqBy))
	if err != nil {
		sdk.options.logger.Errorf(ctx, "doRequest CreateSign err:%v", err)
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBy))
	if err != nil {
		sdk.options.logger.Errorf(ctx, "doRequest NewRequest err:%v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", vd.Authorization)
	resp, err := sdk.httpClient.Do(req)
	if err != nil {
		sdk.options.logger.Errorf(ctx, "doRequest httpClient.Do err:%v", err)
		return err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		sdk.options.logger.Errorf(ctx, "doRequest io.ReadAll err:%v", err)
		return err
	}
	// 响应验签
	respSign := resp.Header.Get("Authorization")

	// 打印响应值
	sdk.options.logger.Infof(ctx, "doRequest resp body:%s auth:%s", string(respBytes), respSign)

	if respSign == "" {
		sdk.options.logger.Errorf(ctx, "doRequest Authorization is nil")
		return fmt.Errorf("签名响应错误")
	}
	err = sg.VerifySign(ctx, sdk.config.APISecret, signx.VerifyData{
		Authorization: respSign,
		SignInfo: signx.SignInfo{
			Method:       method,
			CanonicalURL: url,
			SignBody:     string(respBytes),
		},
	})

	if err != nil {
		sdk.options.logger.Errorf(ctx, "doRequest VerifySign err:%v", err)
		return err
	}

	if respBody != nil {
		return json.Unmarshal(respBytes, respBody)
	}
	return nil
}

// GET请求
func (sdk *CryptoClient) Get(ctx context.Context, path string, req, resp interface{}) error {
	return sdk.doRequest(ctx, http.MethodGet, path, req, resp)
}

func (sdk *CryptoClient) GetPage(ctx context.Context, path string, req, resp interface{}) error {
	return sdk.doRequest(ctx, http.MethodGet, path, req, resp)
}

// POST请求
func (sdk *CryptoClient) Post(ctx context.Context, path string, req, resp interface{}) error {
	return sdk.doRequest(ctx, http.MethodPost, path, req, resp)
}

func (sdk *CryptoClient) PostPage(ctx context.Context, path string, req, resp interface{}) error {
	return sdk.doRequest(ctx, http.MethodPost, path, req, resp)
}

// 实体定义示例
type ExampleRequest struct {
	Param1 string `json:"param1"`
	Param2 int    `json:"param2"`
}

type ExampleResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// 回调处理
func (sdk *CryptoClient) HandleCallback(r *http.Request, callbackUri string) ([]byte, error) {

	if callbackUri == "" {
		callbackUri = r.RequestURI
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		sdk.options.logger.Errorf(r.Context(), "HandleCallback io.ReadAll err:%v", err)
		return nil, err
	}

	sdk.options.logger.Infof(r.Context(), "HandleCallback body:%s", string(body))

	sign := r.Header.Get("Authorization")

	sdk.options.logger.Infof(r.Context(), "HandleCallback sign:%s", sign)

	err = signx.NewSignV1(sdk.options.logger).VerifySign(r.Context(), sdk.config.APISecret, signx.VerifyData{
		Authorization: sign,
		SignInfo: signx.SignInfo{
			Method:       "POST",
			CanonicalURL: callbackUri,
			SignBody:     string(body),
		},
	})
	if err != nil {
		sdk.options.logger.Errorf(r.Context(), "HandleCallback VerifySign err:%v", err)
		return nil, err
	}

	return body, nil
}

// Options 获取客户端配置选项
func (sdk *CryptoClient) TraceIdField() string {
	return sdk.options.traceIdField
}

// Logger 获取日志记录器
func (sdk *CryptoClient) Logger() logger.ILogger {
	return sdk.options.logger
}
