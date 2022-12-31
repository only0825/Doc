package helper

import (
	"github.com/spaolacci/murmur3"
	"github.com/xxtea/xxtea-go/xxtea"
)

func MurHash(str string, key uint32) uint64 {
	h64 := murmur3.New64WithSeed(key)
	h64.Write([]byte(str))

	v := h64.Sum64()
	h64.Reset()
	return v
}

func XxTeaEncrypt(str, key string) string {
	enStr := xxtea.EncryptString(str, key)
	return enStr
}

func XxteaDecrypt(str, key string) string {
	deStr, err := xxtea.DecryptString(str, key)
	if err != nil {
		return ""
	}
	return deStr
}
