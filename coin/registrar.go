package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////

// coins stores the currently registered coins that the system is aware of.
var (
	coinsLock = sync.RWMutex{}
	coins     = map[string]*Coin{}
)

////////////////////////////////////////////////////////////////////////////////

func RegisterCoin(
	name string, port, rpcPort int,
	daemonBin, statusBin, configFile string,
	defBinPath, defDataPath string,
	wdl *WalletDownloader, bdl *BootstrapDownloader,
	fnMap map[string]CoinFunc, opaque interface{}) error {

	coinsLock.Lock()
	defer coinsLock.Unlock()

	if _, ok := coins[name]; ok {
		return fmt.Errorf("coin with name=%s already registered", name)
	}

	coins[name] = &Coin{
		name:    name,
		port:    port,
		rpcPort: rpcPort,

		daemonBin:       daemonBin,
		statusBin:       statusBin,
		configFile:      configFile,
		defaultBinPath:  defBinPath,
		defaultDataPath: defDataPath,

		walletDownloader:    wdl,
		bootstrapDownloader: bdl,

		fnMap:  fnMap,
		opaque: opaque,

		// Computed properties will be set on each command invocation.
		state: &CoinState{},
	}

	return nil
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

	// Update the dynamic properties of the coin like the binary directory
	// and the data directory for each coin. This will use the default versions
	// unless we have an override. This allows us to invoke each CoinFunc more
	// tersely and bundles all extra data into the Coin object.
	c.UpdateDynamic(bins, data)

	// Find and invoke the appropriate coin func (if valid).
	fn, ok := c.fnMap[cmd]
	if !ok {
		return fmt.Errorf("invalid command specified (%s)", cmd)
	}
	if fn == nil {
		fmt.Printf("[%s] %s is a no-op. Doing nothing!\n", name, cmd)
	}
	return fn(c, args)
}

////////////////////////////////////////////////////////////////////////////////
