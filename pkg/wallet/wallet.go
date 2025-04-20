package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/tdadadavid/block/pkg/toolkit"
)

const (
	// CheckSumLength is the length of the checksum we are interested in
	CheckSumLength = 4
	// VERSION is the version of the address
	VERSION = 0x00
)

// Wallet represents a wallet
type Wallet struct {
	//SecretKey is the Private key for the wallet used to verify transaction
	SecretKey ecdsa.PrivateKey `json:"secret_key"`

	// PublicKey is the Public key for the wallet used to sign transaction
	PublicKey []byte `json:"public_key"`
}

// Wallets store all the available wallets in a chain
type Wallets struct {
	wallets map[string]*Wallet
}

// New creates a new wallet
//
// Process
//   - First it generates a new key pair for the wallet
//   - Then it returns the wallet
//
// Parameters
//   - None
//
// Example
//   - w, err := wallet.New()
//   - if err != nil {
//     panic(err)
//     }
//   - fmt.Println(w)
//
// Returns
//   - w(Wallet): The newly created wallet
//   - err(error): The error during the process of creating the wallet
func New() (w *Wallet, err error) {
	priKey, pubKey, err := toolkit.NewKeyPair()
	if err != nil {
		err = fmt.Errorf("failed to create wallet: %w", err)
		return w, err
	}

	w = &Wallet{
		SecretKey: *priKey,
		PublicKey: pubKey,
	}
	return w, err
}

func (w *Wallet) GetPrivateKey() ecdsa.PrivateKey {
	return w.SecretKey
}

func (w *Wallet) GetPublicKey() []byte {
	return w.PublicKey
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
// The reason why the address is Base58 encoded is that it is easier to read and write than hexadecimal, it eliminates similar characters eg i and 1, o and 0
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
func (w *Wallet) GenAddress() (address []byte, err error) {
	// get the public key hash
	pubKeyHash, err := toolkit.PublicKeyHash(w.PublicKey)
	if err != nil {
		err = fmt.Errorf("err generating public key hash: %w", err)
		return address, err
	}

	// version + hash + checksum
	version := []byte{VERSION}
	addr := append(version, pubKeyHash...) // version + hash
	checkSum := toolkit.CheckSum(addr, CheckSumLength)
	addr = append(addr, checkSum...) // version + hash + checksum

	address = toolkit.Base58Encode(addr)

	// base58 encode the address
	fmt.Printf("base58 address for wallet 0x%x\n", address)

	return address, err
}
