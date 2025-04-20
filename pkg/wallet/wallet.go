package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/tdadadavid/block/pkg/toolkit"
	"math/big"
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

func (w *Wallet) Serialize() (data []byte, err error) {
	// Create a buffer to store the serialized data
	var buf bytes.Buffer

	// Write private key components
	// First write D (private key number)
	d := w.SecretKey.D.Bytes()
	if err := toolkit.SerializeString(&buf, string(d)); err != nil {
		err = fmt.Errorf("failed to serialize private key D: %w", err)
		return data, err
	}

	// Write curve parameters (X and Y of public key point)
	x := w.SecretKey.PublicKey.X.Bytes()
	if err := toolkit.SerializeString(&buf, string(x)); err != nil {
		err = fmt.Errorf("failed to serialize public key X: %w", err)
		return data, err
	}

	y := w.SecretKey.PublicKey.Y.Bytes()
	if err := toolkit.SerializeString(&buf, string(y)); err != nil {
		err = fmt.Errorf("failed to serialize public key Y: %w", err)
		return data, err
	}

	// Write public key bytes
	if err := toolkit.SerializeString(&buf, string(w.PublicKey)); err != nil {
		err = fmt.Errorf("failed to serialize public key bytes: %w", err)
		return data, err
	}

	return buf.Bytes(), err
}

func (w *Wallet) Deserialize(data []byte) (err error) {
	buf := bytes.NewReader(data)

	// Read private key parts
	// Read D
	d, err := toolkit.DeserializeString(buf)
	if err != nil {
		err = fmt.Errorf("failed to deserialize private key D: %w", err)
		return err
	}

	// Read public key point coordinates
	x, err := toolkit.DeserializeString(buf)
	if err != nil {
		err = fmt.Errorf("failed to deserialize public key X: %w", err)
		return err
	}

	y, err := toolkit.DeserializeString(buf)
	if err != nil {
		err = fmt.Errorf("failed to deserialize public key Y: %w", err)
		return err
	}

	// Reconstruct the private key
	curve := elliptic.P256() // Using a P256 curve as it's commonly used
	w.SecretKey = ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     new(big.Int).SetBytes([]byte(x)),
			Y:     new(big.Int).SetBytes([]byte(y)),
		},
		D: new(big.Int).SetBytes([]byte(d)),
	}

	// Read public key bytes
	publicKey, err := toolkit.DeserializeString(buf)
	if err != nil {
		err = fmt.Errorf("failed to deserialize public key bytes: %w", err)
		return err
	}
	w.PublicKey = []byte(publicKey)

	return err
}
