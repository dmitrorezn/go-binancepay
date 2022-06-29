/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

type QueryCertificateRequest struct {
}

func (q QueryCertificateRequest) EndPoint() string {
	return "/binancepay/openapi/certificates"
}
func (q *QueryCertificateRequest) Validate() error {
	return validate.Struct(q)
}

type Certificate struct {
	CertSerial string `json:"certSerial"`
	CertPublic string `json:"certPublic"`
}

type QueryCertificateResult = []Certificate
