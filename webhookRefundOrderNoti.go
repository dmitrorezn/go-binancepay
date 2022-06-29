/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

type RefundOrderNoti struct {
	MerchantTradeNo string `json:"merchantTradeNo"`
	ProductType     string `json:"productType"`
	ProductName     string `json:"productName"`
	TradeType       string `json:"tradeType"`
	TotalFee        string `json:"totalFee"`
	Currency        string `json:"currency"`
	OpenUserId      string `json:"openUserId"`
	RefundInfo      string `json:"refundInfo"`
}
