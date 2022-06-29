/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

var _ IRequest = &CloseOrderRequest{}

type CloseOrderRequest struct {
	PrepayId        string `json:"prepayId,omitempty"`
	MerchantTradeNo string `json:"merchantTradeNo,omitempty"`
}

func (q *CloseOrderRequest) Validate() error {
	return validate.Struct(q)
}

func (q *CloseOrderRequest) EndPoint() string {
	return "/binancepay/openapi/order/close"
}

// CloseOrderResult
// equals to true when status="SUCCESS"，which is close request is accepted，
//   and successful close result will be notified asynchronously through Order Notification Webhook
type CloseOrderResult bool
