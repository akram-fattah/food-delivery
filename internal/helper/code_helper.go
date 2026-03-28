package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateCode() string {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "000000"
	}
	return fmt.Sprintf("%06d", n.Int64())
}
