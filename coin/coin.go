// Package coin defines functions common to all masternode-based coins.
// Additionally, each coin which is implemented in this program will register
// itself with the coin package so that we can use this as the generic
// communication point for the client application.
// This draws a lot of inspiration from the `image` libraries in golang.
package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////

type CoinFunc func(c *Coin, binDir, dataDir string, args []string) error

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
	fnMap map[string]CoinFunc

	// Opaque interface for the coin
	opaque interface{}
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
	fnMap map[string]CoinFunc, opaque interface{}) error {

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

		fnMap:  fnMap,
		opaque: opaque,
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

// PrintCoinInfo is a common function that can be used by all coin
// implementations to print common info for a given coin.
func (c *Coin) PrintCoinInfo(binDir, dataDir, prefix string) error {
	bd := c.defaultBinPath + " [default]"
	if len(binDir) > 0 {
		bd = binDir
	}
	bb := c.defaultDataPath + " [default]"
	if len(dataDir) > 0 {
		bb = dataDir
	}

	fmt.Printf(`%s
  * Current Binary Directory: %s
  * Current Data Directory:   %s
  * Coin daemon binary:       %s
  * Coin status binary:       %s
`, prefix, bd, bb, c.daemonBin, c.statusBin)

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (c *Coin) DownloadWallet(dstPath string) error {
	if len(dstPath) == 0 {
		return errors.New("unspecified destination path")
	}
	return c.walletDownloader.DownloadToPath(dstPath)
}

func (c *Coin) DownloadBootstrap(dstPath string) error {
	if len(dstPath) == 0 {
		return errors.New("unspecified destination path")
	}
	return c.bootstrapDownloader.DownloadToPath(dstPath)
}

func (c *Coin) GetOpaque() interface{} {
	return c.opaque
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

// IsRegistered returns true if the coin specified by `name` is known to the
// gomn system.
func IsRegistered(name string) bool {
	coinsLock.RLock()
	defer coinsLock.RUnlock()

	_, ok := coins[name]
	return ok
}

////////////////////////////////////////////////////////////////////////////////

// Command executes a given coin's (specified by `name`), `cmd` function
// if one was registered. If the function was nil, then it has no implementation
// and we do nothing.  If the command was not found we return an error.
func Command(name, bins, data, cmd string, args []string) error {
	coinsLock.Lock()
	defer coinsLock.Unlock()

	c, ok := coins[name]
	if !ok {
		return fmt.Errorf("invalid coin specified (%s)", name)
	}

	fn, ok := c.fnMap[cmd]
	if !ok {
		return fmt.Errorf("invalid command specified (%s)", cmd)
	}

	if fn == nil {
		fmt.Printf("[%s] %s is a no-op. Doing nothing!\n", name, cmd)
	}
	return fn(c, bins, data, args)
}

////////////////////////////////////////////////////////////////////////////////
