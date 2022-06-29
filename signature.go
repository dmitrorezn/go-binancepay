/*
 * Created by Du, Chengbin on 2022/6/14.
 */

package binancepay

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"strings"
)

func Nonce() string {
	id, _ := uuid.GenerateUUID()
	hash := md5.Sum([]byte(id))
	return hex.EncodeToString(hash[:])
}

func Sign(secretKey []byte, payload []byte) (string, error) {
	hmac512 := hmac.New(sha512.New, secretKey)
	_, err := hmac512.Write(payload)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(hmac512.Sum(nil))), nil
}

func BuildPayload(body string, timestampMilli string, nonce string) string {
	return fmt.Sprintf("%s\n%s\n%s\n", timestampMilli, nonce, body)
}

func ParsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, fmt.Errorf("decodePublicPEM(): block is nil")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509.ParsePKIXPublicKey(): %w", err)
	}
	return pub.(*rsa.PublicKey), nil
}

func verifySignature(pub *rsa.PublicKey, payload, signature []byte) error {
	h := crypto.SHA256.New()
	h.Write(payload)
	err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, h.Sum(nil), signature)
	if err != nil {
		return fmt.Errorf("rsa.VerifyPKCS1v15(): %w", err)
	}
	return nil
}
