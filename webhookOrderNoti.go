/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

import "github.com/shopspring/decimal"

type OrderNoti struct {
	MerchantTradeNo string              `json:"merchantTradeNo"`
	TotalFee        decimal.Decimal     `json:"totalFee"`
	TransactTime    int64               `json:"transactTime"`
	Currency        string              `json:"currency"`
	OpenUserId      string              `json:"openUserId"`
	ProductType     string              `json:"productType"`
	ProductName     string              `json:"productName"`
	TradeType       string              `json:"tradeType"`
	TransactionId   string              `json:"transactionId"`
	PayerInfo       *OrderNotiPayerInfo `json:"payerInfo,omitempty"`
}

type OrderNotiPayerInfo struct {
	FirstName      string `json:"firstName"`
	MiddleName     string `json:"middleName"`
	LastName       string `json:"lastName"`
	WalletId       string `json:"walletId"`
	Country        string `json:"country"`
	City           string `json:"city"`
	Address        string `json:"address"`
	IdentityType   string `json:"identityType"`
	IdentityNumber string `json:"identityNumber"`
	DateOfBirth    string `json:"dateOfBirth"`
	PlaceOfBirth   string `json:"placeOfBirth"`
	Nationality    string `json:"nationality"`
}
