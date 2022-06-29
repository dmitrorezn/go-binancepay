/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

var _ IRequest = &QueryOrderRequest{}

type QueryOrderRequest struct {
	PrepayId        string `json:"prepayId,omitempty"`
	MerchantTradeNo string `json:"merchantTradeNo,omitempty"`
}

func (q *QueryOrderRequest) EndPoint() string {
	return "/binancepay/openapi/v2/order/query"
}
func (q *QueryOrderRequest) Validate() error {
	return validate.Struct(q)
}

type QueryOrderResult struct {
	MerchantId      string `json:"merchantId"`
	PrepayId        string `json:"prepayId"`
	TransactionId   string `json:"transactionId"`
	MerchantTradeNo string `json:"merchantTradeNo"`
	Status          string `json:"status"`
	Currency        string `json:"currency"`
	OrderAmount     string `json:"orderAmount"`
	OpenUserId      string `json:"openUserId"`
	TransactTime    int64  `json:"transactTime"`
	CreateTime      int64  `json:"createTime"`
}
