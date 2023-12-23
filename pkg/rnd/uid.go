package rnd

import (
	"strconv"
	"time"
)

const (
	PrefixNone  = byte(0)
	PrefixMixed = byte('*')
)

// GenerateUID returns a unique id with prefix as string.
func GenerateUID(prefix byte) string {
	rnd := GenerateRandomString(25)

	result := make([]byte, 0, 32)
	result = append(result, prefix)
	result = append(result, strconv.FormatInt(time.Now().UTC().Unix(), 36)[0:6]...)
	result = append(result, rnd...)

	return string(result)
}
