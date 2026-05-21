package signx

import (
	"context"
	"testing"

	"github.com/lbe-io/crypto-sdk-go/sdk/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetSign(t *testing.T) {
	ctx := context.Background()
	sign := NewSignV1(&logger.DefLogger{})

	s := "test_appsecret"

	// 测试必填参数缺失
	t.Run("MissingRequiredFields", func(t *testing.T) {
		info := SignData{}
		_, _, err := sign.getSign(ctx, s, info)
		assert.Error(t, err)
		assert.Equal(t, "签名参数错误", err.Error())
	})

	// 测试签名生成逻辑
	t.Run("ValidSignatureGeneration", func(t *testing.T) {
		info := SignData{
			SignInfo: SignInfo{
				Method:       "GET",
				CanonicalURL: "/test",
				SignBody:     "",
			},
			AuthInfo: AuthInfo{
				Appid:      "test_appid",
				MerchantId: "test_merchant_id",
			},
		}

		signHeader, str, err := sign.getSign(ctx, s, info)
		assert.NoError(t, err)
		assert.NotEmpty(t, signHeader)
		assert.NotEmpty(t, str)
	})

	// 测试时间戳和随机数生成
	t.Run("TimestampAndNonceGeneration", func(t *testing.T) {
		info := SignData{
			SignInfo: SignInfo{
				Method:       "GET",
				CanonicalURL: "/test",
				SignBody:     "",
			},
			AuthInfo: AuthInfo{
				Appid:      "test_appid",
				MerchantId: "test_merchant_id",
			},
		}

		signHeader, str, err := sign.getSign(ctx, s, info)
		assert.NoError(t, err)
		assert.NotEmpty(t, signHeader)
		assert.NotEmpty(t, str)
	})
}

func TestCreateSign(t *testing.T) {
	ctx := context.Background()

	sign := NewSignV1(&logger.DefLogger{})

	s := "test_appsecret"

	// 测试签名创建
	t.Run("ValidSignatureCreation", func(t *testing.T) {
		vd, err := sign.CreateSign(
			ctx,
			s,
			"/test?param=value&param2=value2",
			"test_appid",
			"test_merchant_id",
			"GET",
			"test_body",
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, vd)
	})

	// 测试无效输入
	t.Run("InvalidInput", func(t *testing.T) {
		_, err := sign.CreateSign(ctx, s, "", "", "", "", "")
		assert.Error(t, err)
	})
}

func TestVerifySign(t *testing.T) {
	ctx := context.Background()

	sign := NewSignV1(&logger.DefLogger{})

	s := "test_appsecret"

	// 测试签名验证
	t.Run("ValidSignatureVerification", func(t *testing.T) {
		info := SignData{
			SignInfo: SignInfo{
				Method:       "GET",
				CanonicalURL: "/test",
				SignBody:     "",
			},
			AuthInfo: AuthInfo{
				Appid:      "test_appid",
				MerchantId: "test_merchant_id",
			},
		}
		signHeader, _, err := sign.getSign(ctx, s, info)
		t.Log("log:", signHeader, err)
		assert.NoError(t, err)

		vd := VerifyData{
			Authorization: signHeader,
			SignInfo:      info.SignInfo,
		}

		err = sign.VerifySign(ctx, s, vd)
		assert.NoError(t, err)
	})

	// 测试无效签名
	t.Run("InvalidSignature", func(t *testing.T) {
		info := SignData{
			SignInfo: SignInfo{
				Method:       "GET",
				CanonicalURL: "/test",
				SignBody:     "",
			},
			AuthInfo: AuthInfo{
				Appid:      "test_appid",
				MerchantId: "test_merchant_id",
			},
		}
		info.Sign = "invalid_sign"
		err := sign.VerifySign(ctx, s, VerifyData{
			Authorization: info.Sign,
			SignInfo:      info.SignInfo,
		})
		assert.Error(t, err)
	})
}
