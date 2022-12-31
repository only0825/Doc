package helper


import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base32"
    "encoding/binary"
    "time"
)

func TRUNCATE(hs []byte) int {
    offset := int(hs[len(hs)-1] & 0x0F)
    p := hs[offset : offset+4]
    return (int(binary.BigEndian.Uint32(p)) & 0x7FFFFFFF) % 1000000
}

func HMACSHA1(k []byte, c uint64) []byte {
    cb := make([]byte, 8)
    binary.BigEndian.PutUint64(cb, c)

    mac := hmac.New(sha256.New, k)
    mac.Write(cb)

    return mac.Sum(nil)
}

func TOTP(k string, x uint64) int {
    key, err := base32.StdEncoding.DecodeString(k)
    if err != nil {
        return 0
    }

    return HOTP(key, T(0, x))
}

func HOTP(k []byte, c uint64) int {
    return TRUNCATE(HMACSHA1(k, c))
}

func T(t0, x uint64) uint64 {
    return (uint64(time.Now().Unix()) - t0) / x
}

// func main() {
//  fmt.Println(TOTP("ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEF", 30)) // 30秒
// }
