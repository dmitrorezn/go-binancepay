/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

type PayoutNoti struct {
	RequestId   string `json:"requestId"`
	BatchStatus string `json:"batchStatus"`
	MerchantId  string `json:"merchantId"`
	Currency    string `json:"currency"`
	TotalAmount string `json:"totalAmount"`
	TotalNumber string `json:"totalNumber"`
}
