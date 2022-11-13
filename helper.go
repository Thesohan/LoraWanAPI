package main

import (
	"crypto/rand"
	"math/big"
)

// GenerateHexString return hexString of length given in the argument or an error
func GenerateHexString(length int) (string, error) {
	max := big.NewInt(int64(len(AllowedChars)))
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = AllowedChars[n.Int64()]
	}
	return string(b), nil
}

// GenerateNRandomHexString returns a list of unique hex strings of lenght n
func GenerateNRandomHexString(n int) []string {
	set := make(map[string]bool)
	haxStrings := []string{}

	for count := 0; count < 100; count += 1 {
		for {
			hexString, err := GenerateHexString(5)
			if err == nil {
				set[hexString] = true
				haxStrings = append(haxStrings, hexString)
				break
			}
		}
	}
	return haxStrings
}
