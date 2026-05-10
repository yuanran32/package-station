package random

import (
	crand "crypto/rand"
	"math/big"
	"time"
)

const (
	numericCharset      = "0123456789"
	alphaNumericCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func NumericCode(length int) string {
	return randomFromCharset(numericCharset, length)
}

func AlphaNumeric(length int) string {
	return randomFromCharset(alphaNumericCharset, length)
}

func Serial(prefix string) string {
	return prefix + time.Now().Format("20060102150405") + NumericCode(4)
}

func randomFromCharset(charset string, length int) string {
	if length <= 0 {
		return ""
	}

	out := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := range out {
		n, err := crand.Int(crand.Reader, max)
		if err != nil {
			out[i] = charset[0]
			continue
		}
		out[i] = charset[n.Int64()]
	}
	return string(out)
}
