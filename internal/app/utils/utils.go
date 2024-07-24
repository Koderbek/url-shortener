package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func HashGenerator(str string, length int) string {
	hash := md5.New()
	hash.Write([]byte(str))

	return hex.EncodeToString(hash.Sum(nil)[:length])
}
