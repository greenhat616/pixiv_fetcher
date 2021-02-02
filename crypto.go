package main

import "encoding/base64"

// base64DecodeString 用于解密 base64
func base64DecodeString(input []byte) (string, error) {
	tmp := make([]byte, base64.StdEncoding.DecodedLen(len(input)))
	n, err := base64.StdEncoding.Decode(tmp, input)
	return string(tmp[:n]), err
}
