/*
 * Created by Du, Chengbin on 2022/6/15.
 */

package binancepay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloseOrder(t *testing.T) {
	req := &CloseOrderRequest{}
	expectedResp := Response[CloseOrderResult]{
		Status: "SUCCESS",
		Code:   "1",
		Data:   CloseOrderResult(false),
		ErrMsg: "",
	}
	client := NewMerchant("", "", nil, logger)
	client.httpClient = mockHttpClientWithAsserts(t, "POST", "/binancepay/openapi/order/close", req, expectedResp)
	var resp Response[CloseOrderResult]
	err := client.Do(req, &resp)
	assert.Nil(t, err, err)
	assert.Equal(t, expectedResp, resp)
}
