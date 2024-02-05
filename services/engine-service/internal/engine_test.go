package internal

import (
	"testing"
)

func TestEngineService_createOrder(t *testing.T) {

}

func TestEngine_Start(t *testing.T) {
	engines := make(map[string]*Engine)
	coinParis := []string{"btc_usdt", "eth_usdt", "eth_btc"}
	for _, coinPair := range coinParis {
		engine := NewEngine(nil, coinPair, nil)
		go engine.Start(coinPair)
		engines[coinPair] = engine
	}
}
