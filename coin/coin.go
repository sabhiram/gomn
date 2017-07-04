// Package coin defines functions common to all masternode-based coins.
// Additionally, each coin which is implemented in this program will register
// itself with the coin package so that we can use this as the generic
// communication point for the client application.
// This draws a lot of inspiration from the `image` libraries in golang.
package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////

// WalletDownloader is a per-coin wallet fetcher.
type WalletDownloader struct {
	version     string // version of the wallet
	downloadURL string // url to fetch the wallet
	sha256sum   string // shasum for the download
}

// NewWalletDownloader returns a new instance of a wallet downloader.
func NewWalletDownloader(v, url, sha string) *WalletDownloader {
	return &WalletDownloader{
		version:     v,
		downloadURL: url,
		sha256sum:   sha,
	}
}

func (w *WalletDownloader) Log() {
	fmt.Printf("WalletDownloader::%#v\n", w)
}

////////////////////////////////////////////////////////////////////////////////

type BootstrapFn func(binDir, dataDir string) error

////////////////////////////////////////////////////////////////////////////////

// coin represents all things needed to setup / interact-with or monitor
// a given coin's masternode.
type coin struct {
	name        string      // name of the coin, used as the key for lookup
	bootstrapFn BootstrapFn // Fetch and bootstrap the coin daemon
}

// coins stores the currently registered coins that the system is aware of.
var (
	coinsLock = sync.RWMutex{}
	coins     = map[string]*coin{}
)

////////////////////////////////////////////////////////////////////////////////

func RegisterCoin(key string, bstrapFn BootstrapFn) error {
	coinsLock.Lock()
	defer coinsLock.Unlock()

	if _, ok := coins[key]; ok {
		return fmt.Errorf("coin with name=%s already registered", key)
	}

	coins[key] = &coin{
		name:        key,
		bootstrapFn: bstrapFn,
	}

	fmt.Printf("Registered %s\n", key)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func BootstrapCoin(key string) error {
	coinsLock.Lock()
	defer coinsLock.Unlock()

	if coin, ok := coins[key]; ok {
		return coin.bootstrapFn("pivxd", "~/.pivx/")
	}
	return fmt.Errorf("coin %s is not registered with gomn", key)
}

// TODO:
// HTTP fetch helper methods to download files etc...
//
