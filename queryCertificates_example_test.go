/*
 * Created by Du, Chengbin on 2022/6/15.
 */

package binancepay

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Example_queryCertificate() {
	req := &QueryCertificateRequest{}
	client := NewMerchant("", "", nil, logger)
	client.httpClient = mockHttpClient(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(`{"status":"SUCCESS","code":"000000","data":[{"certPublic":"-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwuF7PgDsNBuJVZ6HeVoH\nT2Tqj22+eWLQ/kZbYj3CTgdvedFymj0Kqvxtwy3InVbU6t6g5UkjL+dMOAkt5GxL\nNYI1uNy9g3+mifCDRXDArwvcKkB5jYu0R3WTNtf1ODicjpOf1NUqMZ+0t3jNwVAy\nawvxlyxpX8gMa6OAMbzMtH3iskM52nu5mS57Xh4ryibwjIxd0ssb63gD2qH8jy60\nAK/qgkijlysgEQDzYTk6X2x4t9BfVoOL3+yxkIiwnfL/KY9xkvSmWuAFIZqu4pY7\ng+GXFiG50sSCe2BkBcSzIS56L1Qp/tSzDUl1+fQCGhA3BFY42/zTpvdjLLUgbRYZ\npJCu9Z4w0HsM118rKCxBZNveoc12oHXEbDMDy7y/c39KYNyniCH6iKPNzu6Zi8tb\nXXN7KG9mUHUzstnafVd5QpqIumQUgE+JVTSrdx0YJy7OQeqeSoQeBeI7pr1gsRuH\nx0pcniYOIRwmtZ27Ybkbk0zOu1vBmzE8hC8RAkE8Yz06T7quoa547FicUYQBvtkR\nYLJDbSIjFLkfTFNOgV5VU92JfJvFji3F/nDVQ0gI6iuDktKYB0FNe1LZvKbDgPs+\nJ8/Pssd1DOW8XJbQXmJz8VCrubv/SdOYsy0lP0m/ZybEFjSVSWKT3xCpHVSDVJNm\nrTUypediX9eeNMlfs0x/vmkCAwEAAQ==\n-----END PUBLIC KEY-----","certSerial":"abc"}]}`)),
		}, nil
	})
	var resp Response[QueryCertificateResult]
	err := client.Do(req, &resp)
	fmt.Println(err)
	fmt.Println(resp.Status)
	fmt.Println(len(resp.Data))
	fmt.Println(resp.Data[0].CertSerial)
	fmt.Println(resp.Data[0].CertPublic)
	// Output:
	// <nil>
	// SUCCESS
	// 1
	// abc
	// -----BEGIN PUBLIC KEY-----
	// MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwuF7PgDsNBuJVZ6HeVoH
	// T2Tqj22+eWLQ/kZbYj3CTgdvedFymj0Kqvxtwy3InVbU6t6g5UkjL+dMOAkt5GxL
	// NYI1uNy9g3+mifCDRXDArwvcKkB5jYu0R3WTNtf1ODicjpOf1NUqMZ+0t3jNwVAy
	// awvxlyxpX8gMa6OAMbzMtH3iskM52nu5mS57Xh4ryibwjIxd0ssb63gD2qH8jy60
	// AK/qgkijlysgEQDzYTk6X2x4t9BfVoOL3+yxkIiwnfL/KY9xkvSmWuAFIZqu4pY7
	// g+GXFiG50sSCe2BkBcSzIS56L1Qp/tSzDUl1+fQCGhA3BFY42/zTpvdjLLUgbRYZ
	// pJCu9Z4w0HsM118rKCxBZNveoc12oHXEbDMDy7y/c39KYNyniCH6iKPNzu6Zi8tb
	// XXN7KG9mUHUzstnafVd5QpqIumQUgE+JVTSrdx0YJy7OQeqeSoQeBeI7pr1gsRuH
	// x0pcniYOIRwmtZ27Ybkbk0zOu1vBmzE8hC8RAkE8Yz06T7quoa547FicUYQBvtkR
	// YLJDbSIjFLkfTFNOgV5VU92JfJvFji3F/nDVQ0gI6iuDktKYB0FNe1LZvKbDgPs+
	// J8/Pssd1DOW8XJbQXmJz8VCrubv/SdOYsy0lP0m/ZybEFjSVSWKT3xCpHVSDVJNm
	// rTUypediX9eeNMlfs0x/vmkCAwEAAQ==
	// -----END PUBLIC KEY-----
}
