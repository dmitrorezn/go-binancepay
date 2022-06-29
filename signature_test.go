/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNonce(t *testing.T) {
	nonce := Nonce()
	assert.Len(t, nonce, 32, "nonce length must be 32")
	fmt.Println(nonce)
}

func TestBuildPayload(t *testing.T) {
	payload := BuildPayload("body", "1654943252000", "5f14d51f62136d24c87480faa9009d7d")
	assert.Equal(t, "1654943252000\n5f14d51f62136d24c87480faa9009d7d\nbody\n", payload)
}

func TestSign(t *testing.T) {
	signature, err := Sign([]byte("test secret key"), []byte("content"))
	assert.Nil(t, err, err)
	assert.Equal(t, "AA9DD39D6D4F32068DFD626FFDF019706959B0A811B1A72CD293DB718521A457B1E258461A59CF81FE0F0F32CB9BC1C502D31AAD9938EE73040A9CC6DE04507E", signature)
}

func Test_VerifySignature(t *testing.T) {
	signature := "a/musJjLkWJLbBvm0/Hy11xS5O0dPDNay9d62dIE9O+35PNmkGJz3B189F89gU4lIrOMDaTOLAQv0qqasv/6mjvih5h3MIqaxVdtBCJpHbrd5UZ1jhuNZyFGYTenVNs0x2fTA3AXYi1NWwtCSGAZlf13J6+iVr0jGSH26XL+VjNifUCGiMpPPd0atYEZROL7X3xbzhFsILwGAshvPgtpQ8vEClHgLiFC2tuL5EOZ+pTsB8Rm1QGXOUcrqFYT5LLG0tX8wVSDReaKLVoyzW6wSlxfk/JBK0FdTwZH29zPFWVX5FPh9WjQoquKZkYpXr06xSpzXQ6kpxrcYmsC5JbS0RXqIqS199vEL1VbS8MiDMzjqzkjgCNwfmEZ4Pn4Lmoo2cIsbXUGA7YVMEtyFxYkRwwz8o2Ze9YBdg1ccstrLhRfWA2jpZ5H7lpAYrx6GQCaLwf9fUjeyOJXgndEm5+MMJZk5bAd1W302+ubobM5MCsnNzXonsJLhKcj0O+oCDBCgZrMfL/xGP+1+ifjZSASqWpRnYhBdbaM8jLUlAl7oQSnsAxgGqjg9UrzTOY2xmcxwMhtfPIboVUT6dIqu96LTcIjY3+oIGOj94Y/58iJd09Q0NHXTXLz0WGQAuzJpOwXwiGk+3ShDknn7lAxXL1atNojad6cpf3xrJkncG7UKG8="
	payload := "111"
	public := `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0F0m11B4ACw8+Zt1uRHU
0Yh2qFeLhYA5hsbxjYPA3WJt88tzmV5/lJl8rdxkt6rSrvW2Ad/9Letg7i7lri/8
BxbasG7J1VvTHRMfdsa/0u6MK8W1ZVypkFc0WhTOk0IKwtVrmJPgWvWK4Y0b1cpa
663l1vwp9w8czWFOFlMBSS+scxTuxxV0fqTw12NgoYn94JtVzMnNYw+L8m13V07g
Cv3sPdj1So3Uwj8Ww3HVAntPHrJwl4b5KSTffqEEZNtuCqZyofOlc4YIeb3atBkp
KYTyB77DWDxvUuZMR4Jdy7ae1YeQKeH38R2FishHwi6CkBbW23eclbiyT3eAYjhr
Szhxj9b3F/IciNaxxCASd+AG44IfUOuce9xur5a1vIH/zC06rtNqASt8dGx7QwaL
PjJ5ILAMXgb9EuvGM5eacKunLraB5+CA+yI+/0oe0dupwvK5YYfOKqJ/0VypZlt6
5NGG0FYP/QN+ulcYGFQ/y11yXhtLTEg21VHUGmFjLXk81JLydKBxkM7b+msTNkJL
VDOp0bCekzIks5bYIRYYmmCcvC5uTeFUUlJFnuNWvblygEoRl1MdNk6K9fjSJyrl
5RhAIxJfHbf/K1r0JO42U2kzQdvsfDBMGtYSj61vpZawDJJLi0/NH3TbUDNUFbNQ
3Fzek3Mm5rk99lOHrzar6jkCAwEAAQ==
-----END PUBLIC KEY-----`

	pub, err := ParsePublicKey(public)
	assert.Nil(t, err, err)
	sig, err := base64.StdEncoding.DecodeString(signature)
	assert.Nil(t, err, err)
	err = verifySignature(pub, []byte(payload), sig)
	assert.Nil(t, err, err)
}

func Test_VerifySignature2(t *testing.T) {
	payload := "1654943252000\nabc\n{\"bizType\":\"PAY\",\"bizId\":\"111\",\"data\":\"testdata\",\"bizStatus\":\"SUCCESS\"}\n"
	signature := generateSignature([]byte(payload))

	pub, err := ParsePublicKey(testDataPublicKey)
	assert.Nil(t, err, err)
	sig, err := base64.StdEncoding.DecodeString(signature)
	assert.Nil(t, err, err)
	err = verifySignature(pub, []byte(payload), sig)
	assert.Nil(t, err, err)
}
