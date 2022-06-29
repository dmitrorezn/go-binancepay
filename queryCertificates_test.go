/*
 * Created by Du, Chengbin on 2022/6/15.
 */

package binancepay

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestQueryCertificates(t *testing.T) {
	req := &QueryOrderRequest{}
	expectedResp := Response[QueryOrderResult]{
		Status: "SUCCESS",
		Code:   "1",
		Data: QueryOrderResult{
			MerchantId: "123",
		},
		ErrMsg: "",
	}

	httpClient := mockHttpClient(func(request *http.Request) (*http.Response, error) {
		assert.Equal(t, "/binancepay/openapi/v2/order/query", request.URL.Path)
		assert.Equal(t, http.MethodPost, request.Method)

		reqBodyBytes, err := ioutil.ReadAll(request.Body)
		assert.Nil(t, err, err)
		expectedReq, err := json.Marshal(req)
		assert.Nil(t, err, err)
		assert.Equal(t, string(expectedReq), string(reqBodyBytes))

		respBody, err := json.Marshal(expectedResp)
		assert.Nil(t, err, err)

		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(respBody)),
		}, nil
	})

	client := NewMerchant("", "", nil, logger)
	client.httpClient = httpClient
	var resp Response[QueryOrderResult]
	err := client.Do(req, &resp)
	assert.Nil(t, err, err)
	assert.Equal(t, expectedResp, resp)
}
