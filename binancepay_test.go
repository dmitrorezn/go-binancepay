/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestUnmarshalRequest(t *testing.T) {
	body := `{
  "bizType": "PAY",
  "data": "{\"merchantTradeNo\":\"9825382937292\",\"totalFee\":0.88000000,\"transactTime\":1619508939664,\"currency\":\"BUSD\",\"openUserId\":\"1211HS10K81f4273ac031\",\"productType\":\"Food\",\"productName\":\"Ice Cream\",\"tradeType\":\"WEB\",\"transactionId\":\"M_R_282737362839373\"}",
  "bizId": 29383937493038367292,
  "bizStatus": "PAY_SUCCESS"
}`
	var req webhookRawReq
	err := json.Unmarshal([]byte(body), &req)
	assert.Nil(t, err, err)
	fmt.Println("bizType", req.BizType)
	fmt.Println("bizId", req.BizId)
}

type fixedPublicKeyTestCache struct {
	exist bool
}

func (c fixedPublicKeyTestCache) GetJSON(ctx context.Context, key string, i interface{}) (ok bool, err error) {
	if !c.exist {
		return false, nil
	}

	cert, _ := i.(*Certificate)
	cert.CertSerial = "abc"
	cert.CertPublic = testDataPublicKey
	return true, nil
}

func (c fixedPublicKeyTestCache) SetJSON(ctx context.Context, key string, data interface{}, dur time.Duration) error {
	return nil
}

func TestWebhookWithCachedPublicKey(t *testing.T) {
	client := NewMerchant("", "", fixedPublicKeyTestCache{exist: true}, logger)
	body := `{"bizType":"PAY","bizId":"111","data":"testdata","bizStatus":"SUCCESS"}`
	httpReq, err := http.NewRequest("POST", "/", strings.NewReader(body))
	assert.Nil(t, err, err)
	timestamp := "1654943252000"
	nonce := "abc"

	h := http.Header{}
	h.Set("BinancePay-Certificate-SN", "")
	h.Set("BinancePay-Nonce", nonce)
	h.Set("BinancePay-Timestamp", timestamp)
	h.Set("BinancePay-Signature", "UVVd2NnJMjJRKHbM+e4544nHi2AgToklZdoNs9Gx4GdQPRsm0llIaWPKMM7JFG5bQV9pId/4z2T/wObPhRlDy8wa7UI2rWHpnCq1M48skeZjSzAauUr0u/jLV5BP9atbr9AKzFoAfMeLdpa0knaX3nOiutXzDZjgrb953fhFd//J7fEUN8pCTSXJXUtXUlXEH6bPePpRaiEeNMBC3tnTIQR93dZWpxHPE8LUkrSFAqxdS6m1Byd1b+6ebuDnUjgNjkBEuN04W7tCbirOYcEXRP67bLirc67DTX/k17mXYhdJrVDlobubDMFg3L9ZkOndDDij99DROtonW9SaNpoBxVsoB4fKr0aRaawxMzlUoDwjkuNb0JKOnj/MHB21N4R5qssCgEtAcs3+PcbzBZ943nF03Rne+wEK99+OEzzyzTXOTlmefErOYVxamn8/gnWkXAgvKILyR1S3ldoHOcb6SPWtociEr75z0JUoHpibjmp/WExVSKrSoRJg+JIB3dw3bNoNhdbR41SBSejuk58UPJSa+qcW1tpGj9F057o4sP5knjBoV1q15pFzZRv+qDa8BUSoaNdKrkR7hpCFW+DmKvjYHzD+8Fkr7XdANcEUEdz/0DpRCuhWEAkKtRDzYlVCnuj0cvQ+ZcIZ7jhR7hPYw55jY6rLWGz9FHEtxJRQjtg=")
	h.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header = h

	rawReq, err := client.VerifyAndParseWebhookRequest(httpReq)
	assert.Nil(t, err, err)
	assert.Equal(t, "testdata", rawReq.RawData)
}

func TestWebhookWithNoCachedPublicKey(t *testing.T) {
	client := NewMerchant("", "", fixedPublicKeyTestCache{exist: false}, logger)
	body := `{"bizType":"PAY","bizId":"111","data":"testdata","bizStatus":"SUCCESS"}`
	httpReq, err := http.NewRequest("POST", "/", strings.NewReader(body))
	assert.Nil(t, err, err)
	timestamp := "1654943252000"
	nonce := "abc"

	h := http.Header{}
	h.Set("BinancePay-Certificate-SN", "")
	h.Set("BinancePay-Nonce", nonce)
	h.Set("BinancePay-Timestamp", timestamp)
	h.Set("BinancePay-Signature", "UVVd2NnJMjJRKHbM+e4544nHi2AgToklZdoNs9Gx4GdQPRsm0llIaWPKMM7JFG5bQV9pId/4z2T/wObPhRlDy8wa7UI2rWHpnCq1M48skeZjSzAauUr0u/jLV5BP9atbr9AKzFoAfMeLdpa0knaX3nOiutXzDZjgrb953fhFd//J7fEUN8pCTSXJXUtXUlXEH6bPePpRaiEeNMBC3tnTIQR93dZWpxHPE8LUkrSFAqxdS6m1Byd1b+6ebuDnUjgNjkBEuN04W7tCbirOYcEXRP67bLirc67DTX/k17mXYhdJrVDlobubDMFg3L9ZkOndDDij99DROtonW9SaNpoBxVsoB4fKr0aRaawxMzlUoDwjkuNb0JKOnj/MHB21N4R5qssCgEtAcs3+PcbzBZ943nF03Rne+wEK99+OEzzyzTXOTlmefErOYVxamn8/gnWkXAgvKILyR1S3ldoHOcb6SPWtociEr75z0JUoHpibjmp/WExVSKrSoRJg+JIB3dw3bNoNhdbR41SBSejuk58UPJSa+qcW1tpGj9F057o4sP5knjBoV1q15pFzZRv+qDa8BUSoaNdKrkR7hpCFW+DmKvjYHzD+8Fkr7XdANcEUEdz/0DpRCuhWEAkKtRDzYlVCnuj0cvQ+ZcIZ7jhR7hPYw55jY6rLWGz9FHEtxJRQjtg=")
	h.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header = h

	{
		req := &QueryCertificateRequest{}
		expectedResp := Response[QueryCertificateResult]{
			Status: "SUCCESS",
			Code:   "1",
			Data: []Certificate{
				{
					CertSerial: "abc",
					CertPublic: testDataPublicKey,
				},
			},
			ErrMsg: "",
		}
		client.httpClient = mockHttpClientWithAsserts(t, "POST", "/binancepay/openapi/certificates", req, expectedResp)
	}

	rawReq, err := client.VerifyAndParseWebhookRequest(httpReq)
	assert.Nil(t, err, err)
	assert.Equal(t, "testdata", rawReq.RawData)
}
