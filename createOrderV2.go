/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

import (
	"github.com/shopspring/decimal"
)

var _ IRequest = &CreateOrderV2Request{}

type SubMerchant struct {
	SubMerchantId string `json:"subMerchantId"`
}

type Env struct {
	TerminalType  string `json:"terminalType" validate:"oneof=APP WEB WAP MINI_PROGRAM OTHERS"`
	OsType        string `json:"osType,omitempty"`
	OrderClientIp string `json:"orderClientIp,omitempty"`
	CookieId      string `json:"cookieId,omitempty"`
}

type Amount struct {
	Currency string          `json:"currency" validate:"required"`
	Amount   decimal.Decimal `json:"amount" validate:"required"`
}

type Goods struct {
	GoodsType string `json:"goodsType" validate:"required,oneof=01 02"` // "01": Tangible Goods "02": Virtual Goods
	// GoodsCategory
	// 0000: Electronics & Computers
	// 1000: Books, Music & Movies
	// 2000: Home, Garden & Tools
	// 3000: Clothes, Shoes & Bags
	// 4000: Toys, Kids & Baby
	// 5000: Automotive & Accessories
	// 6000: Game & Recharge
	// 7000: Entertainament & Collection
	// 8000: Jewelry
	// 9000: Domestic service
	// A000: Beauty care
	// B000: Pharmacy
	// C000: Sports & Outdoors
	// D000: Food, Grocery & Health products
	// E000: Pet supplies
	// F000: Industry & Science
	// Z000: Others
	GoodsCategory    string  `json:"goodsCategory" validate:"required"`
	ReferenceGoodsId string  `json:"referenceGoodsId" validate:"required"` // The unique ID to identify the goods.
	GoodsName        string  `json:"goodsName" validate:"required"`        // Goods name
	GoodsDetail      string  `json:"goodsDetail,omitempty"`                // Optional.
	GoodsUnitAmount  *Amount `json:"goodsUnitAmount,omitempty"`            // Optional.
	GoodsQuantity    string  `json:"goodsQuantity"`                        // Quantity of goods
}

// CreateOrderV2Request doc https://developers.binance.com/docs/binance-pay/api-order-create-v2#child-attribute
type CreateOrderV2Request struct {
	SubMerchant     SubMerchant     `json:"merchant"`
	Env             Env             `json:"env"`
	MerchantTradeNo string          `json:"merchantTradeNo" validate:"required"` // maximum length 32
	OrderAmount     decimal.Decimal `json:"orderAmount" validate:"required"`     // Range: 0.01 - 20000
	Currency        string          `json:"currency" validate:"required"`        // order currency in upper case. only "BUSD","USDT","MBOX" can be accepted, fiat NOT supported.
	Goods           Goods           `json:"goods" validate:"required"`
	//Shipping           string      `json:"shipping"`
	//Buyer              string      `json:"buyer"`
	ReturnUrl          string `json:"returnUrl"`
	CancelUrl          string `json:"cancelUrl"`
	OrderExpireTime    int64  `json:"orderExpireTime,omitempty"` // milliseconds
	SupportPayCurrency string `json:"supportPayCurrency"`        //  e.g. "BUSD,BNB"
	AppId              string `json:"appId,omitempty"`           // This field is required when terminalType is MINI_PROGRAM
	UniversalUrlAttach string `json:"universalUrlAttach"`
}

func (r *CreateOrderV2Request) EndPoint() string {
	return "/binancepay/openapi/v2/order"
}

func (r *CreateOrderV2Request) Validate() error {
	return validate.Struct(r)
}

type CreateOrderV2Result struct {
	PrepayId     string `json:"prepayId"`
	TerminalType string `json:"terminalType"`
	ExpireTime   int64  `json:"expireTime"` //expire time in milliseconds
	QrcodeLink   string `json:"qrcodeLink"`
	QrContent    string `json:"qrContent"`
	CheckoutUrl  string `json:"checkoutUrl"`
	Deeplink     string `json:"deeplink"`
	UniversalUrl string `json:"universalUrl"`
}
