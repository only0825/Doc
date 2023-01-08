const CryptoJS = require("crypto-js");

// 辅助函数
function md5(data) {
    return CryptoJS.MD5(data).toString();
}


// 传入key之前要调用，不然结果不对
function parseKey(key) {
    return CryptoJS.enc.Utf8.parse(key);
}


// 加密过程
function encrypt(mode, plainText, key, iv = null) {
    const uKey = parseKey(key);
    const uIv = parseKey(iv);

    return CryptoJS.AES.encrypt(plainText, uKey,
        {
            iv: uIv,
            mode: mode,
            padding: CryptoJS.pad.Pkcs7
        }
    ).toString();
}

// 解密过程
function decrypt(mode, cipherText, key, iv = null) {
    const uKey = parseKey(key);
    const uIv = parseKey(iv);

    let bytes = CryptoJS.AES.decrypt(cipherText, uKey,
        {
            iv: uIv,
            mode: mode,
            padding: CryptoJS.pad.Pkcs7
        }
    );
    return bytes.toString(CryptoJS.enc.Utf8);
}

function test() {
    // AES的key，并算出偏移量iv
    const key = 'ef14b146e989b922dcb5a00a19f907c5';
    const md5Key = md5(key);
    const iv = md5Key.substr(0, 16);
    // 从/api/live/signurl获取到的加密字符串
    const cipherText = 'oDJZrCJy4/RZIPJC4hZ0hcA8ZDHhpXkJfLH990H+rC6QQRMBcIW70Wm47A0JoTb/MRzYC9jbSoBl9SLHPokWHoU3aHIgPO9a0P+fHXTwCBh4h6jtEUiPsXGcfAZBzjMVbrSyboT5Y41jkd3a7HibZQ==';

    let plainText = decrypt(CryptoJS.mode.CBC, cipherText, key, iv);
    console.log(plainText);
}

test();
