package slack

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func VerifySignature(
	secret []byte,
	timestamp string,
	requestBody string,
	signature string,
) bool {
	base := fmt.Sprintf("v0:%s:%s", timestamp, requestBody)

	mac := hmac.New(sha256.New, secret)

	_, _ = mac.Write([]byte(base))

	sig := fmt.Sprintf("v0=%x", mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(sig))
}
