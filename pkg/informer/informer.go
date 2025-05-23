package informer

import (
	"sync"
	"testgo/pkg/client"
)

type CurrencyInfo struct {
	Price_24h        float64 `json:"price_24h"`
	Volume_24h       float64 `json:"volume_24h"`
	Last_trade_price float64 `json:"last_trade_price"`
}
type CurrencyInfoStr struct {
	Price_24h        string `json:"price_24h"`
	Volume_24h       string `json:"volume_24h"`
	Last_trade_price string `json:"last_trade_price"`
}

type Informer struct {
	CurencyMap map[string]CurrencyInfo
	mu         sync.Mutex
}

func NewInformer() *Informer {
	return &Informer{
		CurencyMap: make(map[string]CurrencyInfo),
	}
}

func (i *Informer) AddInMap(v client.Currency) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.CurencyMap[v.Symbol] = CurrencyInfo{
		Price_24h:        v.Price,
		Volume_24h:       v.Volume,
		Last_trade_price: v.LastPrice,
	}
}

func (i *Informer) GetResultMap() map[string]CurrencyInfo {
	i.mu.Lock()
	defer i.mu.Unlock()

	copyMap := make(map[string]CurrencyInfo)
	for k, v := range i.CurencyMap {
		copyMap[k] = v
	}

	return copyMap
}
