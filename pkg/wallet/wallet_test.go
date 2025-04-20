package wallet

import (
	"github.com/stretchr/testify/assert"
	"github.com/tdadadavid/block/pkg/toolkit"
	"os"
	"testing"
)

func cleanUp(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll("./data")
		if err != nil {
			panic(err)
		}
	})
}

func TestWallet_New(t *testing.T) {
	w, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, w)
	assert.NotNil(t, w.GetPrivateKey())
	assert.NotNil(t, w.GetPublicKey())
}

func TestWallet_NewWallets(t *testing.T) {
	defer cleanUp(t)

	wallets := NewWallets()

	assert.NotNil(t, wallets)
}

func TestWallet_AddWallet(t *testing.T) {
	defer cleanUp(t)

	wallets := NewWallets()

	wallet1, err := New()
	wallet2, err := New()
	assert.NoError(t, err)

	wallets.AddWallet(wallet1)
	wallets.AddWallet(wallet2)

	assert.Len(t, wallets.wallets, 2)
}

func TestWallet_GetWallet(t *testing.T) {
	defer cleanUp(t)

	wallets := NewWallets()

	wallet1, err := New()
	wallet2, err := New()
	assert.NoError(t, err)

	wallets.AddWallet(wallet1)
	wallets.AddWallet(wallet2)

	w2 := wallets.GetWallet(wallet2.GetPublicKey())
	assert.NotNil(t, w2)
	assert.Equal(t, w2.GetPublicKey(), wallet2.GetPublicKey())
	assert.Equal(t, w2.GetPrivateKey(), wallet2.GetPrivateKey())
	assert.NotEqual(t, w2.GetPublicKey(), wallet1.GetPublicKey())
}

func TestWallet_GenAddress(t *testing.T) {
	w, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, w)

	address, err := w.GenAddress()
	assert.NoError(t, err)
	assert.NotEmpty(t, address)
}

func TestWallet_Address(t *testing.T) {
	w, err := New()
	w2, err := New()

	address, err := w.GenAddress()
	address2, err := w2.GenAddress()
	assert.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.NotEmpty(t, address2)
	assert.NotEqual(t, address, address2)
}

func TestWallet_Hash(t *testing.T) {
	w, err := New()
	pubKey := w.GetPublicKey()
	hash, err := toolkit.PublicKeyHash(pubKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Equal(t, len(hash), 20)
}
