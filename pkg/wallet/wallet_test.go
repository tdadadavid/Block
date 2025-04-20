package wallet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWallet_New(t *testing.T) {
	w := New()
	assert.NotNil(t, w)
}
