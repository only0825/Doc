
let crypto;
try {
    crypto = require('crypto');
} catch (err) {
    console.log('crypto support is disabled!');
}

const ALGORITHM = 'aes-256-cbc';
const CIPHER_KEY = "abcdefghijklmnopqrstuvwxyz012345";  // Same key used in Golang
const BLOCK_SIZE = 16;

const plainText = "1234567890";  // This plainText was encrypted to make the cipherText below by Golang
const cipherText = "f17ba46472fa64e40ca496d1b4c91e8fac967926dfbdd7097b4c8f8ebd18f898";  // hexidecimal cipherText created by Golang

const decrypted = decrypt(cipherText);

if (decrypted !== plainText) {
    console.log(`FAILED: expected ${plainText} but got "${decrypted}"`);
} else {
    console.log(`PASSED: ${plainText}`);
}

function test() {
    let encrypt1 = encrypt(plainText);
    let decrypt1 = decrypt(cipherText);
}

// Decrypts cipher text into plain text
function decrypt(cipherText) {
    const contents = Buffer.from(cipherText, 'hex');
    const iv = contents.slice(0, BLOCK_SIZE);
    const textBytes = contents.slice(BLOCK_SIZE);

    const decipher = crypto.createDecipheriv(ALGORITHM, CIPHER_KEY, iv);
    let decrypted = decipher.update(textBytes, 'hex', 'utf8');
    decrypted += decipher.final('utf8');
    return decrypted;
}

// Encrypts plain text into cipher text
function encrypt(plainText) {
    const iv = crypto.randomBytes(BLOCK_SIZE);
    const cipher = crypto.createCipheriv(ALGORITHM, CIPHER_KEY, iv);
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

test()