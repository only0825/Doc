
let crypto;
try {
    crypto = require('crypto');
} catch (err) {
    console.log('crypto support is disabled!');
}



function encrypt(plainText, key) {
    const iv = crypto.randomBytes(16);
    const cipher = crypto.createCipheriv("aes-256-cbc", key, iv);
    let cipherText;
    try {
        cipherText = cipher.update(plainText, 'utf8', 'hex');
        cipherText += cipher.final('hex');
        cipherText = iv.toString('hex') + cipherText
    } catch (e) {
        cipherText = null;
    }
    return cipherText;
}

function test() {
    let plainText = new Date().getTime().toString();  // 要加密信息，当前时间戳，必须为13位，精确到毫秒
    let encrypt1 = encrypt(plainText, "d0caea82c21bb744d54cc84bc4d0a430");
    console.log(encrypt1)
}

test()