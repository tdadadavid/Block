package wallet

import (
	"github.com/stretchr/testify/assert"
	"github.com/tdadadavid/block/pkg/toolkit"
	"testing"
)

func TestWallet_New(t *testing.T) {
	w, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, w)
	assert.NotNil(t, w.GetPrivateKey())
	assert.NotNil(t, w.GetPublicKey())
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
