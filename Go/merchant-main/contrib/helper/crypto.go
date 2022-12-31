package helper

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	md5simd "github.com/minio/md5-simd"

	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// md5
func MD5Hash(text string) string {

	server := md5simd.NewServer()
	md5Hash := server.NewHash()
	_, _ = md5Hash.Write([]byte(text))
	digest := md5Hash.Sum([]byte{})
	encrypted := hex.EncodeToString(digest)

	server.Close()
	md5Hash.Close()

	return encrypted
}

// sha1
func Sha1Sum(s string) []byte {

	h := sha1.New()
	h.Write([]byte(s))

	return h.Sum(nil)
}

func HmacSha(source string, key string) string {

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(source))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// sha256
func Sha256sum(param []byte) string {

	h := sha256.New()
	h.Write(param)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func RsaEncrypt(privateKey, origData []byte) []byte {

	//设置私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil
	}

	prkI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}

	priv := prkI.(*rsa.PrivateKey)
	encodeByte, _ := rsa.SignPKCS1v15(crand.Reader, priv, crypto.MD5, origData)

	return encodeByte
}

func rsaSha256Sign(privateKey, origData []byte) []byte {

	//设置私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil
	}

	prkI, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}

	h := sha256.New()
	h.Write(origData)
	data := h.Sum(nil)

	encodeByte, _ := rsa.SignPKCS1v15(crand.Reader, prkI, crypto.SHA256, data)

	return encodeByte
}

func rsaSha256Very(publicKey, data, signData []byte) bool {

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return false
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		return false
	}
	return true
}

// aes ecd
func AesEcbEncrypt(data, key []byte) []byte {

	block, _ := aes.NewCipher(key)

	data = pkcs5Padding(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

// aes ecd
func AesEcbDecrypt(crypted, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	origData := make([]byte, len(crypted))

	size := block.BlockSize()

	for bs, be := 0, size; bs < len(crypted); bs, be = bs+size, be+size {
		block.Decrypt(origData[bs:be], crypted[bs:be])
	}

	origData = pkcs5UnPadding(origData)
	return origData, nil
}

func cryptBlocks(block cipher.Block, origData, crypted []byte) {

	for len(crypted) > 0 {
		block.Decrypt(origData, crypted[:block.BlockSize()])
		crypted = crypted[block.BlockSize():]
		origData = origData[block.BlockSize():]
	}
}

// aes cbc
func AesCbcEncrypt(plaintext []byte, secretKey, iv string) []byte {

	keyBytes := []byte(secretKey)
	aesBlockCipher, _ := aes.NewCipher(keyBytes)
	content := pkcs5Padding(plaintext, aesBlockCipher.BlockSize())
	encrypted := make([]byte, len(content))
	aesBlockMode := cipher.NewCBCEncrypter(aesBlockCipher, []byte(iv))
	aesBlockMode.CryptBlocks(encrypted, content)

	return encrypted
}

// aes cbc
func AesCbcDecrypt(src, secretKey, iv string) ([]byte, error) {

	crypted, err := base64.StdEncoding.DecodeString(src) //将字符串变成[]byte
	if err != nil {
		return nil, err
	}
	keyBytes := []byte(secretKey)
	aesBlockCipher, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(aesBlockCipher, []byte(iv))
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)

	return origData, nil
}

// des
func DesEncrypt(src, key []byte) (string, error) {

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	bs := block.BlockSize()
	src = pkcs5Padding(src, bs)
	if len(src)%bs != 0 {
		return "", errors.New("need a multiple of the block size")
	}

	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}

	return Base64Encode(out), nil
}

func DesDecrypt(src, key []byte) ([]byte, error) {

	if len(src) < 1 {
		return nil, errors.New("src nil")
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}

	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = pkcs5UnPadding(out)

	return out, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {

	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {

	length := len(origData)
	unPadding := int(origData[length-1])

	if length-unPadding < 0 {
		return origData
	}

	return origData[:(length - unPadding)]
}

// 3DES加密
func TripleDesEncrypt(data, key []byte) ([]byte, error) {

	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	data = pkcs5Padding(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted, nil
}

func Base64Encode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
