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
	name            string // name of the coin, used as the key for lookup
	daemonBin       string // name of the daemon to launch the coin's node
	statusBin       string // name of the binary to check status
	defaultBinPath  string // path to where the bins will exist
	defaultDataPath string // path to where the data will exist

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
	defBinPath, defDataPath string,
	wdl *WalletDownloader, bdl *BootstrapDownloader,
	bootstrapFn BootstrapFn) error {

	coinsLock.Lock()
	defer coinsLock.Unlock()

	if _, ok := coins[name]; ok {
		return fmt.Errorf("coin with name=%s already registered", name)
	}

	coins[name] = &Coin{
		name:            name,
		daemonBin:       daemonBin,
		statusBin:       statusBin,
		defaultBinPath:  defBinPath,
		defaultDataPath: defDataPath,

		walletDownloader:    wdl,
		bootstrapDownloader: bdl,

		bootstrapFn: bootstrapFn,
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (c *Coin) DownloadWallet(binDir, dataDir string) error {
	if len(binDir) > 0 {
		if err := c.walletDownloader.DownloadToPath(binDir); err != nil {
			fmt.Printf("Got error: %s\n", err.Error())
			return err
		}
	}

	if len(dataDir) > 0 {
		if err := c.bootstrapDownloader.DownloadToPath(dataDir); err != nil {
			fmt.Printf("Got error: %s\n", err.Error())
			return err
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func BootstrapCoin(name string) error {
	coinsLock.Lock()
	defer coinsLock.Unlock()

	if coin, ok := coins[name]; ok {
		// TODO: Pass flags into these functions :) Each fn should decide
		// 		 which path to use
		return coin.bootstrapFn(coin, coin.defaultBinPath, coin.defaultDataPath)
	}
	return fmt.Errorf("coin %s is not registered with gomn", name)
}

////////////////////////////////////////////////////////////////////////////////

// RegisteredCoins returns a list of coins that the tool knows how to configure.
func RegisteredCoins() []string {
	ret := []string{}
	coinsLock.RLock()
	defer coinsLock.RUnlock()
	for coin, _ := range coins {
		ret = append(ret, coin)
	}
	return ret
}

////////////////////////////////////////////////////////////////////////////////
