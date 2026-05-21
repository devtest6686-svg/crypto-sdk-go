package signx

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lbe-io/crypto-sdk-go/sdk/logger"
	"github.com/spf13/cast"
)

type SignVersion string

const (
	SignVersionV1 SignVersion = "Crypto-V1"
)

// 设计思想参考微信支付
// https://pay.weixin.qq.com/doc/v3/merchant/4012365334

type SignV1 struct {
	logger logger.ILogger
}

func NewSignV1(log logger.ILogger) *SignV1 {
	return &SignV1{
		logger: log,
	}
}

// 签名有效期
var WindowsDuration = time.Minute * 5

// 签名参数
type AuthInfo struct {
	MerchantId string `json:"merchant_id"` // 必填 商户ID Header
	Appid      string `json:"appid"`       // 必填 Header
	SignTime   int64  `json:"sign_time"`   // 毫秒时间戳 Header
	Nonce      string `json:"nonce"`       // 随机数 Header
	Sign       string `json:"sign"`        // 签名值 Header
}

// 签名信息
type SignInfo struct {
	Method       string `json:"method"`        // 必填
	CanonicalURL string `json:"canonical_url"` // 必填 参数(去掉Host) 如/user/info?name=123&age=18
	SignBody     string `json:"sign_body"`     // 请求体
}

// 验签参数
type VerifyData struct {
	Authorization string `json:"authorization"` // 必填(在Header中) 如 Crypto-V2 merchant_id:xxx&appid:xxx&nonce:xxx&sign_time:xxx&sign:xxx
	SignInfo
}

// 签名参数
type SignData struct {
	AuthInfo
	SignInfo
}

// 创建签名
func (s *SignV1) CreateSign(ctx context.Context, appsecret, canonicalURL string, appid string, merchantId string, method string, signBody string) (*VerifyData, error) {

	// 解析URL 去掉host
	u, err := url.Parse(canonicalURL)
	if err != nil {
		s.logger.Errorf(ctx, "CreateSign url.Parse err: %v", err)
		return nil, err
	}

	info := SignData{
		SignInfo: SignInfo{
			Method:       method,
			CanonicalURL: u.RequestURI(),
			SignBody:     signBody,
		},
	}
	info.Appid = appid
	info.MerchantId = merchantId
	signHeader, _, err := s.getSign(ctx, appsecret, info)
	if err != nil {
		s.logger.Errorf(ctx, "CreateSign getSign err: %v", err)
		return nil, err
	}

	return &VerifyData{
		Authorization: signHeader,
		SignInfo:      info.SignInfo,
	}, nil
}

// 验证签名
func (s *SignV1) VerifySign(ctx context.Context, appsecret string, verifyData VerifyData) error {

	if verifyData.Authorization == "" {
		return errors.New("VerifySign Authorization 不能为空")
	}

	authInfo, err := s.ParseAuthStr(ctx, verifyData.Authorization)
	if err != nil {
		s.logger.Errorf(ctx, "VerifySign ParseAuthStr err: %v", err)
		return err
	}

	signData := SignData{
		AuthInfo: *authInfo,
		SignInfo: verifyData.SignInfo,
	}

	auth, sign, err := s.getSign(ctx, appsecret, signData)
	if err != nil {
		s.logger.Errorf(ctx, "VerifySign getSign err: %v", err)
		return err
	}
	s.logger.Infof(ctx, "VerifySign signData.localAuth:%s", auth)

	if signData.Sign != sign {
		s.logger.Errorf(ctx, "VerifySign signData.Sign: %s, infoResp.Sign: %s", signData.Sign, sign)
		return errors.New("签名验证失败")
	}

	return nil
}

// 生成签名
func (s *SignV1) getSign(ctx context.Context, appsecret string, info SignData) (signHeader string, sign string, err error) {
	// 创建副本避免修改原始参数
	signData := info

	if signData.MerchantId == "" ||
		signData.Appid == "" ||
		signData.Method == "" ||
		signData.CanonicalURL == "" {
		s.logger.Errorf(ctx, "getSign info.MerchantId: %s, info.Appid: %s, info.Method: %s, info.CanonicalURL: %s",
			signData.MerchantId, signData.Appid, signData.Method, signData.CanonicalURL)
		return "", "", errors.New("签名参数错误")
	}
	// Method 转换为大写
	signData.Method = strings.ToUpper(signData.Method)

	if signData.SignTime == 0 {
		signData.SignTime = time.Now().UnixMilli()
	}
	// 有效时间正负5分钟
	if signData.SignTime < time.Now().Add(-WindowsDuration).UnixMilli() || signData.SignTime > time.Now().Add(WindowsDuration).UnixMilli() {
		s.logger.Errorf(ctx, "getSign info.SignTime: %d, time.Now().UnixMilli(): %d", signData.SignTime, time.Now().UnixMilli())
		return "", "", errors.New("签名时间错误")
	}
	// 随机数
	if signData.Nonce == "" {
		u, err := uuid.NewV7()
		if err != nil {
			s.logger.Errorf(ctx, "uuid.NewV7 err: %v", err)
			return "", "", err
		}
		signData.Nonce = u.String()
	}
	// 解析URL 去掉host
	u, err := url.Parse(signData.CanonicalURL)
	if err != nil {
		s.logger.Errorf(ctx, "getSign url.Parse err: %v", err)
		return "", "", err
	}
	signData.CanonicalURL = u.RequestURI()

	signStr := s.buildSignStr(ctx, signData)
	s.logger.Infof(ctx, "getSign GetSignature rawData: %s", signStr)

	// 使用 HMAC-SHA256 算法生成签名
	h := hmac.New(sha256.New, []byte(appsecret))
	h.Write([]byte(signStr))
	signature := hex.EncodeToString(h.Sum(nil))
	signData.Sign = signature

	auth, err := s.buildAuthStr(ctx, signData)
	if err != nil {
		s.logger.Errorf(ctx, "getSign BuildAuthorizationV2 err: %v", err)
		return "", "", err
	}

	return auth, signData.Sign, nil
}

// 构建签名字符串
func (s *SignV1) buildSignStr(ctx context.Context, info SignData) string {
	str := fmt.Sprintf("merchant_id:%s&appid:%s&method:%s&nonce:%s&sign_time:%d&canonical_url:%s&sign_body:%s",
		info.MerchantId, info.Appid, info.Method, info.Nonce, info.SignTime, info.CanonicalURL, info.SignBody)
	return str
}

// 构造 V2 Auth
func (s *SignV1) buildAuthStr(ctx context.Context, info SignData) (string, error) {
	str := fmt.Sprintf("merchant_id:%s&appid:%s&nonce:%s&sign_time:%d&sign:%s", info.MerchantId, info.Appid, info.Nonce, info.SignTime, info.Sign)

	str = fmt.Sprintf("%s %s", SignVersionV1, str)

	return str, nil
}

// 解析 V2 Auth
// Crypto-V2 merchant_id:xxx&appid:xxx&nonce:xxx&sign_time:xxx&sign:xxx
func (s *SignV1) ParseAuthStr(ctx context.Context, auth string) (*AuthInfo, error) {

	str, ok := strings.CutPrefix(auth, string(SignVersionV1))
	if !ok {
		s.logger.Errorf(ctx, "ParseAuthStr strings.CutPrefix cut not ok")
		return nil, fmt.Errorf("签名版本异常")
	}

	// 版本2签名
	str = strings.TrimSpace(str)

	strs := strings.Split(str, "&")
	data := AuthInfo{}
	for _, s := range strs {
		kv := strings.Split(s, ":")
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "appid":
			data.Appid = kv[1]
		case "sign":
			data.Sign = kv[1]
		case "nonce":
			data.Nonce = kv[1]
		case "sign_time":
			data.SignTime = cast.ToInt64(kv[1])
		case "merchant_id":
			data.MerchantId = kv[1]
		}
	}

	// 验证参数
	// 有效时间正负5分钟
	if data.SignTime < time.Now().Add(-WindowsDuration).UnixMilli() || data.SignTime > time.Now().Add(WindowsDuration).UnixMilli() {
		s.logger.Errorf(ctx, "ParseAuthStr info.SignTime: %d, time.Now().UnixMilli(): %d", data.SignTime, time.Now().UnixMilli())
		return nil, errors.New("签名时间错误")
	}

	return &data, nil

}
