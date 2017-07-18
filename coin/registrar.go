package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"sync"

	"github.com/sabhiram/gomn/types"
)

////////////////////////////////////////////////////////////////////////////////

type CoinFunc func(c *Coin, args []string) error

type FunctionMap struct {
	InfoFn      CoinFunc
	DownloadFn  CoinFunc
	BootstrapFn CoinFunc
	ConfigureFn CoinFunc
	GetInfoFn   CoinFunc
}

func (fm *FunctionMap) Validate(c *Coin) error {
	if fm.GetInfoFn == nil {
		return fmt.Errorf("Warning: %s does not implement required command: %s", c.name, "getinfo")
	}
	if fm.InfoFn == nil {
		return fmt.Errorf("Warning: %s does not implement required command: %s", c.name, "info")
	}
	if fm.DownloadFn == nil {
		return fmt.Errorf("Warning: %s does not implement required command: %s", c.name, "download")
	}
	if fm.BootstrapFn == nil {
		return fmt.Errorf("Warning: %s does not implement required command: %s", c.name, "bootstrap")
	}
	if fm.ConfigureFn == nil {
		return fmt.Errorf("Warning: %s does not implement required command: %s", c.name, "configure")
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
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
	defWalletPath, defBinSubPath, defDataPath string,
	wdl *WalletDownloader, bdl *BootstrapDownloader,
	fnMap *FunctionMap, opaque interface{}) error {

	coinsLock.Lock()
	defer coinsLock.Unlock()

	if _, ok := coins[name]; ok {
		return fmt.Errorf("coin with name=%s already registered", name)
	}

	coins[name] = &Coin{
		name:    name,
		port:    port,
		rpcPort: rpcPort,

		daemonBin:  daemonBin,
		statusBin:  statusBin,
		configFile: configFile,

		defaultWalletPath: defWalletPath,
		defaultBinSubPath: defBinSubPath,
		defaultDataPath:   defDataPath,

		walletDownloader:    wdl,
		bootstrapDownloader: bdl,

		fnMap:  fnMap,
		opaque: opaque,

		// Computed properties will be set on each command invocation.
		state: &CoinState{},
	}
	return coins[name].fnMap.Validate(coins[name])
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

func GetCoinByName(name string) (*Coin, error) {
	coinsLock.RLock()
	defer coinsLock.RUnlock()

	if c, ok := coins[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("invalid coin specified (%s)", name)
}

////////////////////////////////////////////////////////////////////////////////

// Command executes a given coin's (specified by `name`), `cmd` function
// if one was registered. If the function was nil, then it has no implementation
// and we do nothing.  If the command was not found we return an error.
func Command(cli *types.CLI, cmd string, opts []string) error {
	coinsLock.Lock()
	defer coinsLock.Unlock()

	c, ok := coins[cli.Coin]
	if !ok {
		return fmt.Errorf("invalid coin specified (%s)", cli.Coin)
	}

	// Update the dynamic properties of the coin like the binary directory
	// and the data directory for each coin. This will use the default versions
	// unless we have an override. This allows us to invoke each CoinFunc more
	// tersely and bundles all extra data into the Coin object.
	c.UpdateDynamic(cli.Wallet, cli.BinPath, cli.DataPath)

	// Find and invoke the appropriate coin func (if valid).
	switch cmd {
	case "info":
		return c.fnMap.InfoFn(c, opts)
	case "download":
		return c.fnMap.DownloadFn(c, opts)
	case "bootstrap":
		return c.fnMap.BootstrapFn(c, opts)
	case "configure":
		return c.fnMap.ConfigureFn(c, opts)
	case "getinfo":
		return c.fnMap.GetInfoFn(c, opts)
	default:
		return fmt.Errorf("invalid command specified (%s)", cmd)
	}
}

////////////////////////////////////////////////////////////////////////////////
