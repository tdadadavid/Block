package wallet

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
)

const (
	VERSION  = 0x00
	HASH160  = 0x14
	CHECKSUM = 0x04
)

// Wallet represents a wallet
type Wallet struct {
	//SecretKey is the Private key for the wallet used to verify transaction
	SecretKey []byte `json:"secret_key"`

	// PublicKey is the Public key for the wallet used to sign transaction
	PublicKey []byte `json:"public_key"`
}

// Wallets store all the available wallets in a chain
type Wallets struct {
	wallets map[string]*Wallet
}

func New() (w *Wallet, err error) {
	priKey, pubKey, err := NewKeyPair()
	if err != nil {
		err = fmt.Errorf("failed to create wallet: %w", err)
		return w, err
	}
	w = &Wallet{
		SecretKey: priKey.D.Bytes(),
		PublicKey: pubKey.([]byte),
	}
	return w, err
}

// GenAddress generates a new address for the wallet
//
// Behavior
//   - First it hashes the public key with SHA256, then it passes it through RIPEMD160 to generate a friendly address with 20 bytes
//     then appends prefix "0x" to the generated address and returns it
//   - Bitcoin address syntax is  VERSION + HASH160 + CHECKSUM
//   - VERSION is 0x00 for mainnet and 0x6F for testnet
//   - HASH160 is the first 20 bytes of the RIPEMD160 hash of the public key
//   - CHECKSUM is the last 4 bytes of the SHA256 hash of the VERSION + HASH160
//
// Process
//   - If the wallet already has an address, it returns the address
//   - If the wallet doesn't have an address, it generates a new address and returns it
//
// Parameters
//   - None
//
// Returns
//   - address(string): The address generated for the wallet
//   - error(error): The error during the process of generating the address
func (w *Wallet) GenAddress() string {
	// get the public key hash
	pubKeyHash, err := PublicKeyHash(*w)
	if err != nil {
		fmt.Println("err generating public key hash: " + err.Error())
		return ""
	}

	// version + hash + checksum

	//
}

func PublicKeyHash(w Wallet) (hash []byte, err error) {
	// hash the public key with SHA256
	pubKeyHash := sha256.Sum256(w.PublicKey)

	ripemd160Hasher := ripemd160.New()
	_, err = ripemd160Hasher.Write(pubKeyHash[:])
	if err != nil {
		err = fmt.Errorf("failed to create public key hash: %w", err)
		return hash, err
	}

	hash = ripemd160Hasher.Sum(nil)
	return hash, err
}

func NewKeyPair() (priKey *ecdsa.PrivateKey, pubKey crypto.PublicKey, err error) {
	priKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		err = fmt.Errorf("failed to create wallet: %w", err)
		return priKey, pubKey, err
	}
	pubKey = priKey.Public()

	return priKey, pubKey, err
}
