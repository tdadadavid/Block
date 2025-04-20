package toolkit

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

// Base58Encode encodes the given data
//
// Process
//   - First it encodes the data with Base58
//   - Then it returns the encoded data
//
// Parameters
//   - data(byte): The data to be encoded
//
// Returns
//   - encoded(byte): The encoded data
//   - err(error): The error during the process of encoding the data
//
// Example
//   - encoded := Base58Encode([]byte("Hello World"))
//   - fmt.Println(encoded)
//   - // Output: 2NEpo7TZRRrLZSi2U
func Base58Encode(data []byte) []byte {
	encoded := base58.Encode(data)
	return []byte(encoded)
}

// Base58Decode decodes the given data
//
// Process
//   - First it decodes the data with Base58
//   - Then it returns the decoded data
//
// Parameters
//   - data(byte): The data to be decoded
//
// Returns
//   - v(byte): The decoded data
//   - err(error): The error during the process of decoding the data
func Base58Decode(data []byte) (v []byte, err error) {
	decoded, err := base58.Decode(string(data))
	if err != nil {
		err = fmt.Errorf("err decoding base58: " + err.Error())
		return v, err
	}
	return decoded, err
}

// CheckSum calculates the checksum of the given data
//
// Process
//   - First it hashes the data with SHA256 twice
//   - Then it returns the last 4 bytes of the SHA256 hash
//
// Parameters
//   - data(byte): The data to be hashed
func CheckSum(data []byte, CheckSumLength int) []byte {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return second[:CheckSumLength]
}

// PublicKeyHash calculates the public key hash of the given wallet
//
// Process
//   - First it hashes the public key with SHA256
//   - Then it passes it through RIPEMD160 to generate a friendly address with 20 bytes
//     then appends prefix "0x" to the generated address and returns it
//
// Parameters
//   - w(Wallet): The wallet to calculate the public key hash
//
// Returns
//   - hash(byte): The public key hash of the given wallet
//   - error(error): The error during the process of calculating the public key hash
func PublicKeyHash(data []byte) (hash []byte, err error) {
	// hash the public key with SHA256
	pubKeyHash := sha256.Sum256(data)

	ripemd160Hasher := ripemd160.New()
	_, err = ripemd160Hasher.Write(pubKeyHash[:])
	if err != nil {
		err = fmt.Errorf("failed to create public key hash: %w", err)
		return hash, err
	}

	hash = ripemd160Hasher.Sum(nil)
	return hash, err
}

// NewKeyPair generates a new key pair for the wallet
//
// Process
//   - First it generates a new ECDSA key pair
//   - Then it returns the private key and the public key
//
// Parameters
//   - None
//
// NOTE
//   - The private key is a 32-byte big-endian integer
//   - The public key is a 65-byte big-endian integer
//   - The public key is the concatenation of the private key's x-coordinate and y-coordinate
//
// Returns
//   - priKey(ecdsa.PrivateKey): The private key for the wallet
//   - pubKey(crypto.PublicKey): The public key for the wallet
//   - err(error): The error during the process of generating the key pair
func NewKeyPair() (priKey *ecdsa.PrivateKey, pubKey []byte, err error) {
	priKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		err = fmt.Errorf("failed to create wallet: %w", err)
		return priKey, pubKey, err
	}

	// a public key is the concatenation of the private key's x-coordinate and y-coordinate
	pubKey = append(priKey.X.Bytes(), priKey.Y.Bytes()...)

	return priKey, pubKey, err
}
