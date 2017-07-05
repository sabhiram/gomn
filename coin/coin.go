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

type BootstrapFn func(c *Coin, binDir, dataDir string) error

////////////////////////////////////////////////////////////////////////////////

// coin represents all things needed to setup / interact-with or monitor
// a given coin's masternode.
type Coin struct {
	// Coin specific constants
	name      string // name of the coin, used as the key for lookup
	daemonBin string // name of the daemon to launch the coin's node
	statusBin string // name of the binary to check status

	// Coin specific downloaders (can be nil0)
	walletDownloader    *WalletDownloader
	bootstrapDownloader *BootstrapDownloader

	// Coin specific functions to invoke!
	bootstrapFn BootstrapFn // Fetch and bootstrap the coin daemon
}

// coins stores the currently registered coins that the system is aware of.
var (
	coinsLock = sync.RWMutex{}
	coins     = map[string]*Coin{}
)

////////////////////////////////////////////////////////////////////////////////

func RegisterCoin(
	name, daemonBin, statusBin string,
	wdl *WalletDownloader, bdl *BootstrapDownloader,
	bootstrapFn BootstrapFn) error {

	coinsLock.Lock()
	defer coinsLock.Unlock()

	if _, ok := coins[name]; ok {
		return fmt.Errorf("coin with name=%s already registered", name)
	}

	coins[name] = &Coin{
		name:      name,
		daemonBin: daemonBin,
		statusBin: statusBin,

		walletDownloader:    wdl,
		bootstrapDownloader: bdl,

		bootstrapFn: bootstrapFn,
	}

	fmt.Printf("Registered %s\n", name)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (c *Coin) DownloadWallet() error {
	if err := c.walletDownloader.FetchWalletToPath("/Users/shaba/Desktop/work/code/go/src/github.com/sabhiram/gomn/test"); err != nil {
		fmt.Printf("Got error: %s\n", err.Error())
	} else {
		fmt.Printf("All good baby!\n")
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func BootstrapCoin(name string) error {
	coinsLock.Lock()
	defer coinsLock.Unlock()

	if coin, ok := coins[name]; ok {
		return coin.bootstrapFn(coin, "pivxd", "~/.pivx/")
	}
	return fmt.Errorf("coin %s is not registered with gomn", name)
}

// TODO:
// HTTP fetch helper methods to download files etc...
//
