/*
 * Created by Du, Chengbin on 2022/6/15.
 */

package binancepay

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"testing"
)

var logger, _ = zap.NewDevelopment()

type mockTransport struct {
	roundTrip func(request *http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return t.roundTrip(request)
}

func mockHttpClient(roundTrip func(request *http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{
		Transport: &mockTransport{
			roundTrip: roundTrip,
		},
	}
}

func mockHttpClientWithAsserts(t *testing.T, expectedMethod, expectedPath string, expectedReq any, resp any) *http.Client {
	return mockHttpClient(func(request *http.Request) (*http.Response, error) {
		assert.Equal(t, expectedPath, request.URL.Path)
		assert.Equal(t, expectedMethod, request.Method)

		reqBodyBytes, err := ioutil.ReadAll(request.Body)
		assert.Nil(t, err, err)
		expectedReqBytes, err := json.Marshal(expectedReq)
		assert.Nil(t, err, err)
		assert.Equal(t, string(expectedReqBytes), string(reqBodyBytes))

		respBody, err := json.Marshal(resp)
		assert.Nil(t, err, err)

		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil
	})
}

func generateSignature(data []byte) string {
	block, _ := pem.Decode([]byte(testDataPrivateKey))
	if block == nil {
		return ""
	}

	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			panic(err)
		}
		rawkey = rsa
	default:
		panic(fmt.Errorf("ssh: unsupported key type %q", block.Type))
	}

	rsaPrivateKey, ok := rawkey.(*rsa.PrivateKey)
	if !ok {
		panic("not *rsa.PrivateKey")
	}

	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, d)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(sig)
}

var testDataPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

var testDataPublicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAsXqzjaDMzMAABdgntw9b
0qZcNe/oWSlLEAxHYXS9rlCE2HMzFllh57uERPYWvUqz8yJ2m8PvfU2tCUziZedi
EyZvNf6+tXam2bGhZhvz30C2O15RTT31fR1dw8wx1YZk0VXTSf5LM1DCW5gJ1ISB
6zxk/wmvIC67PabTi7NIP7F7IYpgbSYyrkhCsXi8tp8upQIDPxE5XX1Km8eFxo1a
XTjtkHhMMDnvWIGXkrUDQ1kBwbEP/GJYNwAu89M+5nE4IBokLs5XrYgr+bdSU6sP
UZ4FIzUcw9bADtLLrIz/zsnHZiglpPdE5c3SIwYxff9fHh//26IqvQkB8eqOCUo5
u1pPaim8G1sAMlauPlMTbZ6R5HT7jvU0WBQJXojr1ZmgdvzUvBD5rYo8q9fG/2ty
zTWdI7WGcj7FBiS5eg4ssE8nr4xV4uuAhvdlDaqO9TTbAJklAXFdajTGUv84/Kxi
Qg5Q5lB1EjZK+0BlP346nnv0rUog6Tpg91tRVx02BaqJMt2GT+wJZAgLMARcS/Qh
YGtt0gPtF0i49M1tWm5iWwpySEvG5Kyuy8D7E+72gKCVVKcO9bzvOTzZndo0xx4S
AmhaXJZdf0RzjnUYAvvilgGCwSlOzZCX64uO6svrfqprDd3tENMaBGaJDMJbUnlg
Q2aSlx1O9of58Flyzk/F0x0CAwEAAQ==
-----END PUBLIC KEY-----`
