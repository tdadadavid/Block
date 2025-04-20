package wallet

import (
	"context"
	"fmt"
	"github.com/tdadadavid/block/pkg/store"
)

// Wallets store all the available wallets in a chain
type Wallets struct {
	wallets map[string]*Wallet
}

// NewWallets create a new wallets
func NewWallets() (w Wallets) {
	w = Wallets{
		wallets: make(map[string]*Wallet),
	}

	// open the wallet store
	ws, err := store.Open("./data/wallets")
	if err != nil {
		panic(fmt.Errorf("failed to create wallets %v", err))
	}

	// find all the wallets in the store
	data, err := ws.FindAllWallets(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to create wallets %v", err))
	}

	// load all the wallets into the wallets map
	for _, walletData := range data {
		var wallet Wallet
		err = wallet.Deserialize(walletData)
		if err != nil {
			panic(fmt.Errorf("failed to create wallets %v", err))
		}

		w.wallets[string(wallet.PublicKey)] = &wallet
	}

	return w
}

func (w *Wallets) AddWallet(wallet *Wallet) {
	w.wallets[string(wallet.PublicKey)] = wallet
}

func (w *Wallets) GetWallet(pubKey []byte) *Wallet {
	return w.wallets[string(pubKey)]
}
