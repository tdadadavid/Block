package wallet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWallet_New(t *testing.T) {
	w, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, w)
	assert.NotNil(t, w.GetPrivateKey())
	assert.NotNil(t, w.GetPublicKey())
}
