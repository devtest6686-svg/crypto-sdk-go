package client

import (
	"testing"
)

func TestNewClient(t *testing.T) {

	config := SDKConfig{
		ServerHost: "https://api.example.com",
		MerchantID: "your_merchant_id",
		APIKey:     "your_api_key",
		APISecret:  "your_api_secret",
	}

	sdk, err := NewClient(config)
	if err != nil {
		t.Fatalf("初始化失败: %v", err)
	}
	if sdk.config.ServerHost != config.ServerHost || sdk.config.MerchantID != config.MerchantID || sdk.config.APIKey != config.APIKey || sdk.config.APISecret != config.APISecret {
		t.Fatalf("参数未正确赋值")
	}
}

func TestNewClient_EmptyParams(t *testing.T) {

	config := SDKConfig{
		ServerHost: "https://api.example.com",
		MerchantID: "your_merchant_id",
		APIKey:     "your_api_key",
		APISecret:  "your_api_secret",
	}

	_, err := NewClient(config)
	if err != nil {
		t.Fatalf("缺少参数时应返回错误")
	}
}
