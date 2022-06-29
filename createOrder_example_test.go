/*
 * Created by Du, Chengbin on 2022/6/15.
 */

package binancepay

import (
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"strings"
)

func Example_createOrder() {
	req := &CreateOrderV2Request{
		Env: Env{
			TerminalType: "WEB",
		},
		MerchantTradeNo: Nonce(),
		Currency:        "USDT",
		OrderAmount:     decimal.NewFromFloat(0.01),
		Goods: Goods{
			GoodsType:        "02",
			GoodsCategory:    "Z000",
			ReferenceGoodsId: "test_goods_id",
			GoodsName:        "test goods",
		},
	}
	client := NewMerchant("", "", nil, logger)
	client.httpClient = mockHttpClient(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(`{"status":"SUCCESS","code":"000000","data":{"prepayId":"1234","terminalType":"WEB","expireTime":1655859345000,"qrcodeLink":"","qrContent":"","checkoutUrl":"","deeplink":"","universalUrl":""}}`)),
		}, nil
	})
	var resp Response[CreateOrderV2Result]
	err := client.Do(req, &resp)
	fmt.Println(err)
	fmt.Println(resp.Status)
	fmt.Println(resp.Data.TerminalType)
	fmt.Println(resp.Data.PrepayId)
	// Output:
	// <nil>
	// SUCCESS
	// WEB
	// 1234
}
