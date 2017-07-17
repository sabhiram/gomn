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
	"path/filepath"

	"github.com/sabhiram/gomn/cmdargs"
)

////////////////////////////////////////////////////////////////////////////////

var (
	emptyMap = map[string]string{}
)

////////////////////////////////////////////////////////////////////////////////

type CoinFunc func(c *Coin, args []string) error

////////////////////////////////////////////////////////////////////////////////

// CoinState represents a collection of dynamic coin properties that are only
// known at run-time.
type CoinState struct {
	walletPath       string // path where the wallet will exist
	walletPathExists bool   // true if the above path exists
	binPath          string // path where the bins will exist
	binPathExists    bool   // true if the above path exists
	daemonBinPath    string // populated if the daemon binary exists at the specified path
	daemonBinExists  bool   // true if the above file exists
	statusBinPath    string // populated if the status binary exists at the specified path
	statusBinExists  bool   // true if the above file exists

	dataPath         string            // path where the data will exist
	dataPathExists   bool              // true if the above path exists
	configFilePath   string            // path to config file
	configFileExists bool              // true if the above file exists
	config           map[string]string // k-v map of `coin`.conf file
}

////////////////////////////////////////////////////////////////////////////////

// Coin represents all things needed to setup / interact-with or monitor
// a given coin's masternode.
type Coin struct {
	// Coin specific constants
	name    string // name of the coin, used as the key for lookup
	port    int    // port number for the coins peering
	rpcPort int    // port number for RPC comm

	daemonBin         string // name of the daemon to launch the coin's node
	statusBin         string // name of the binary to check status
	configFile        string // name of the config file in the data directory
	defaultWalletPath string // default path where the wallet is extracted
	defaultBinSubPath string // subpath to the binaries
	defaultDataPath   string // default path where the data will exist

	// Coin specific downloaders (can be nil0)
	walletDownloader    *WalletDownloader
	bootstrapDownloader *BootstrapDownloader

	// Coin specific functions to invoke!
	fnMap map[string]CoinFunc

	// Opaque interface for the coin
	opaque interface{}

	// Computed state for the given coin based on input parameters etc
	state *CoinState
}

////////////////////////////////////////////////////////////////////////////////

func (c *Coin) UpdateDynamic(wallet, subpath, data string) error {
	if c == nil {
		return errors.New("nil coin, cannot update dynamic portions")
	}

	////////////////////////////////////////////////////////////

	c.state.walletPath = c.defaultWalletPath
	if len(subpath) > 0 {
		c.state.walletPath = subpath
	}
	c.state.walletPathExists = DirExists(c.state.walletPath)

	subpathToBins := c.defaultBinSubPath
	if len(subpath) > 0 {
		subpathToBins = subpath
	}
	c.state.binPath = filepath.Join(c.state.walletPath, subpathToBins)
	c.state.binPathExists = DirExists(c.state.binPath)

	c.state.daemonBinPath = filepath.Join(c.state.binPath, c.daemonBin)
	c.state.daemonBinExists = FileExists(c.state.daemonBinPath)

	c.state.statusBinPath = filepath.Join(c.state.binPath, c.statusBin)
	c.state.statusBinExists = FileExists(c.state.statusBinPath)

	////////////////////////////////////////////////////////////

	c.state.dataPath = c.defaultDataPath
	if len(data) > 0 {
		c.state.dataPath = data
	}

	c.state.dataPathExists = DirExists(c.state.dataPath)
	c.state.configFilePath = filepath.Join(c.state.dataPath, c.configFile)
	c.state.configFileExists = FileExists(c.state.configFilePath)
	c.state.config = emptyMap

	////////////////////////////////////////////////////////////

	if c.state.configFileExists {
		m, err := LoadConfFile(c.state.configFilePath)
		if err != nil {
			return err
		}
		c.state.config = m
	}

	////////////////////////////////////////////////////////////

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (c *Coin) GetOpaque() interface{} {
	return c.opaque
}

func (c *Coin) GetBinPath() string {
	if c == nil || c.state == nil {
		return ""
	}
	return c.state.binPath
}

func (c *Coin) GetDataPath() string {
	if c == nil || c.state == nil {
		return ""
	}
	return c.state.dataPath
}

func (c *Coin) GetConfFilePath() string {
	if c == nil || c.state == nil {
		return ""
	}
	return c.state.configFilePath
}

func (c *Coin) GetPort() int {
	if c == nil {
		return -1
	}
	return c.port
}

func (c *Coin) GetRPCPort() int {
	if c == nil {
		return -1
	}
	return c.rpcPort
}

func (c *Coin) GetConfig() map[string]string {
	if c == nil {
		return emptyMap
	}
	return c.state.config
}

// GetConfigValue returns the value for a given key in the config file.  Returns
// an empty string if the key is not found.
func (c *Coin) GetConfigValue(key string) string {
	m := c.GetConfig()
	if v, ok := m[key]; ok {
		return v
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// PrintCoinInfo is a common function that can be used by all coin
// implementations to print common info for a given coin.
func (c *Coin) PrintCoinInfo(prefix string) error {
	phelper := func(s string, ok bool) string {
		st := "MISSING"
		if ok {
			st = "     OK"
		}
		return fmt.Sprintf("[ %s ] %s", st, s)
	}

	fmt.Printf(`%s
  * Base directory:     %s
  * Binary directory:   %s
  * Coin daemon binary: %s
  * Coin status binary: %s
  * Data directory:     %s
  * Config file:        %s
`,
		prefix,
		phelper(c.state.walletPath, c.state.walletPathExists),
		phelper(c.state.binPath, c.state.binPathExists),
		phelper(c.state.daemonBinPath, c.state.daemonBinExists),
		phelper(c.state.statusBinPath, c.state.statusBinExists),
		phelper(c.state.dataPath, c.state.dataPathExists),
		phelper(c.state.configFilePath, c.state.configFileExists))

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (c *Coin) DownloadWallet(args []string, override *cmdargs.Download) error {
	if c.state.walletPathExists &&
		c.state.binPathExists &&
		c.state.daemonBinExists &&
		c.state.statusBinExists {
		return errors.New("wallet binary already exists (TODO: Add --force option)")
	}
	return c.walletDownloader.DownloadToPath(c.state.walletPath, override)
}

func (c *Coin) DownloadBootstrap(args []string, override *cmdargs.Bootstrap) error {
	if c.state.dataPathExists {
		return errors.New("wallet data already exists (TODO: add --force option)")
	}
	return c.bootstrapDownloader.DownloadToPath(c.state.dataPath, override)
}

////////////////////////////////////////////////////////////////////////////////
