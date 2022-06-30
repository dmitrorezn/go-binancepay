/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"
)

type Cache interface {
	GetJSON(ctx context.Context, key string, i interface{}) (ok bool, err error)
	SetJSON(ctx context.Context, key string, data interface{}, dur time.Duration) error
}

const (
	DefaultHost = "https://bpay.binanceapi.com"
)

type Merchant struct {
	host       string
	apiKey     string
	secret     []byte
	logger     *zap.Logger
	httpClient *http.Client

	cert          *Certificate
	certPublicKey *rsa.PublicKey

	requestID uint64
	cache     Cache // store certificate
}

func NewMerchant(apiKey, secret string, cache Cache, logger *zap.Logger) *Merchant {
	return &Merchant{
		host:       DefaultHost,
		apiKey:     apiKey,
		secret:     []byte(secret),
		logger:     logger,
		httpClient: http.DefaultClient,
		requestID:  0,
		cache:      cache,
	}
}

type NotiBizType string

type webhookRawReq struct {
	BizType   NotiBizType `json:"bizType"`
	BizId     json.Number `json:"bizId"`
	RawData   string      `json:"data"`
	BizStatus string      `json:"bizStatus"`
}

type IResponse interface {
	Success() bool
	GetError() error
}

type Response[T any] struct {
	Status string `json:"status"`
	Code   string `json:"code"` // https://developers.binance.com/docs/binance-pay/api-order-create-v2
	Data   T      `json:"data"`
	ErrMsg string `json:"errorMessage"`
}

func (r *Response[T]) Success() bool {
	return r.Status == "SUCCESS"
}

func (r *Response[T]) GetError() error {
	if r.Success() {
		return nil
	}
	return fmt.Errorf("resp status=%s code=%s errorMessage=%s", r.Status, r.Code, r.ErrMsg)
}

func (m *Merchant) VerifyAndParseWebhookRequest(r *http.Request) (*webhookRawReq, error) {
	timestamp := r.Header.Get("Binancepay-Timestamp")
	nonce := r.Header.Get("Binancepay-Nonce")
	signatureStr := r.Header.Get("Binancepay-Signature")

	entityBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("readReqBody(): %w", err)
	}

	m.logger.Debug("verify and parse webhook request",
		zap.String("Binancepay-Timestamp", timestamp),
		zap.String("Binancepay-Nonce", nonce),
		zap.String("Binancepay-Signature", signatureStr),
		zap.ByteString("body", entityBody),
	)

	// Load certificates
	if m.certPublicKey == nil {
		m.logger.Debug("load binance cert")

		cacheKey := "binance-pay:cert:" + m.apiKey
		var cert Certificate
		exists, err := m.cache.GetJSON(context.TODO(), cacheKey, &cert)
		if err != nil {
			m.logger.Error("failed to get binance cert from cache", zap.Error(err))
			return nil, fmt.Errorf("queryCertificatesFromCache(): %w", err)
		}

		if !exists {
			req := &QueryCertificateRequest{}
			var resp Response[QueryCertificateResult]
			if err = m.Do(req, &resp); err != nil {
				m.logger.Error("failed to get binance cert from query certificates API", zap.Error(err))
				return nil, fmt.Errorf("queryCertificates(): %w", err)
			}
			if len(resp.Data) == 0 {
				return nil, fmt.Errorf("got empty certificates response")
			}
			cert.CertSerial = resp.Data[0].CertSerial
			cert.CertPublic = resp.Data[0].CertPublic

			if err = m.cache.SetJSON(context.TODO(), cacheKey, cert, time.Hour*24*365); err != nil {
				return nil, fmt.Errorf("cacheCertificate(): %w", err)
			}
		}

		pub, err := ParsePublicKey(cert.CertPublic)
		if err != nil {
			return nil, fmt.Errorf("ParsePublicKey(): %w", err)
		}
		m.cert = &cert
		m.certPublicKey = pub
	}

	signature, err := base64.StdEncoding.DecodeString(signatureStr)
	if err != nil {
		return nil, fmt.Errorf("decodeSignature(): %w", err)
	}

	payload := BuildPayload(string(entityBody), timestamp, nonce)

	err = verifySignature(m.certPublicKey, []byte(payload), signature)
	if err != nil {
		return nil, fmt.Errorf("verifySignature(): %w", err)
	}

	request := webhookRawReq{}
	if err = json.Unmarshal(entityBody, &request); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %w", err)
	}

	return &request, nil
}

func (m *Merchant) WebhookResponse(w http.ResponseWriter, success bool, message string) error {
	resp := map[string]interface{}{}
	resp["returnCode"] = "FAIL"
	if success {
		resp["returnCode"] = "SUCCESS"
	}
	resp["returnMessage"] = nil
	if message != "" {
		resp["returnMessage"] = message
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	if _, err = w.Write(data); err != nil {
		return err
	}
	return nil
}

type IRequest interface {
	EndPoint() string
	Validate() error
}

type HttpMethodProvider interface {
	HttpMethod() string
}

func (m *Merchant) Do(req IRequest, response IResponse) (err error) {
	logger := m.logger.With(zap.Uint64("id", atomic.AddUint64(&m.requestID, 1)))

	if err = req.Validate(); err != nil {
		return err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("req.MarshalJSON(): %w", err)
	}

	nonce := Nonce()
	timestampMilli := fmt.Sprintf("%d", time.Now().UnixMilli())
	payload := BuildPayload(string(body), timestampMilli, nonce)
	signature, err := Sign(m.secret, []byte(payload))
	if err != nil {
		return fmt.Errorf("sign(): %w", err)
	}

	h := http.Header{}
	h.Set("BinancePay-Certificate-SN", m.apiKey)
	h.Set("BinancePay-Nonce", nonce)
	h.Set("BinancePay-Timestamp", timestampMilli)
	h.Set("BinancePay-Signature", signature)
	h.Set("Content-Type", "application/json;charset=utf-8")

	method := "POST"
	if httpMethodProvider, ok := req.(HttpMethodProvider); ok {
		method = httpMethodProvider.HttpMethod()
	}

	logger.Debug("new request",
		zap.String("method", method),
		zap.String("endpoint", req.EndPoint()),
		zap.ByteString("body", body),
		zap.Strings("header", []string{
			m.apiKey,
			nonce,
			timestampMilli,
			signature,
		}),
	)
	httpReq, err := http.NewRequest(method, m.host+req.EndPoint(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http.NewRequest(): %w", err)
	}

	httpReq.Header = h
	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("httpClient.Do(): %w", err)
	}

	if resp.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAlll(resp.Body): %w", err)
	}

	logger.Debug("got resp", zap.String("status", resp.Status), zap.ByteString("body", respBytes))

	if err = json.Unmarshal(respBytes, response); err != nil {
		return fmt.Errorf("json.Unmarshal(respBytes): %w", err)
	}

	if !response.Success() {
		return response.GetError()
	}

	return nil
}
