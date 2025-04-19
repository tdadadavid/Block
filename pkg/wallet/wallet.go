package wallet

type Wallet struct {
	Address string `json:"address"`
}

func New() *Wallet {
	return &Wallet{}
}
