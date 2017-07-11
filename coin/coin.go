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
)

////////////////////////////////////////////////////////////////////////////////

type CoinFunc func(c *Coin, args []string) error

////////////////////////////////////////////////////////////////////////////////

// CoinState represents a collection of dynamic coin properties that are only
// known at run-time.
type CoinState struct {
	binPath         string // path where the bins will exist
	binPathExists   bool   // true if the above path exists
	daemonBinPath   string // populated if the daemon binary exists at the specified path
	daemonBinExists bool   // true if the above file exists
	statusBinPath   string // populated if the status binary exists at the specified path
	statusBinExists bool   // true if the above file exists

	dataPath         string // path where the data will exist
	dataPathExists   bool   // true if the above path exists
	configFilePath   string // path to config file
	configFileExists bool   // true if the above file exists
}

////////////////////////////////////////////////////////////////////////////////

// Coin represents all things needed to setup / interact-with or monitor
// a given coin's masternode.
type Coin struct {
	// Coin specific constants
	name string // name of the coin, used as the key for lookup
	port int    // port number for the coins peering

	daemonBin       string // name of the daemon to launch the coin's node
	statusBin       string // name of the binary to check status
	configFile      string // name of the config file in the data directory
	defaultBinPath  string // default path where the bins will exist
	defaultDataPath string // default path where the data will exist

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

func (c *Coin) UpdateDynamic(bins, data string) error {
	if c == nil {
		return errors.New("nil coin, cannot update dynamic portions")
	}

	////////////////////////////////////////////////////////////

	c.state.binPath = c.defaultBinPath
	if len(bins) > 0 {
		c.state.binPath = bins
	}

	c.state.binPathExists = DirExists(c.state.binPath)
	c.state.daemonBinPath = filepath.Join(c.state.binPath, c.daemonBin)
	c.state.statusBinPath = filepath.Join(c.state.binPath, c.statusBin)
	c.state.daemonBinExists = FileExists(c.state.daemonBinPath)
	c.state.statusBinExists = FileExists(c.state.statusBinPath)

	////////////////////////////////////////////////////////////

	c.state.dataPath = c.defaultDataPath
	if len(data) > 0 {
		c.state.dataPath = data
	}

	c.state.dataPathExists = DirExists(c.state.dataPath)
	c.state.configFilePath = filepath.Join(c.state.dataPath, c.configFile)
	c.state.configFileExists = FileExists(c.state.configFilePath)

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
  * Binary Directory:   %s
  * Coin daemon binary: %s
  * Coin status binary: %s
  * Data Directory:     %s
  * Config File:        %s
`,
		prefix,
		phelper(c.state.binPath, c.state.binPathExists),
		phelper(c.state.daemonBinPath, c.state.daemonBinExists),
		phelper(c.state.statusBinPath, c.state.statusBinExists),
		phelper(c.state.dataPath, c.state.dataPathExists),
		phelper(c.state.configFilePath, c.state.configFileExists))

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (c *Coin) DownloadWallet() error {
	if c.state.binPathExists &&
		c.state.daemonBinExists &&
		c.state.statusBinExists {
		return errors.New("wallet binary already exists (TODO: Add --force option)")
	}
	return c.walletDownloader.DownloadToPath(c.state.binPath)
}

func (c *Coin) DownloadBootstrap() error {
	if c.state.dataPathExists {
		return errors.New("wallet data already exists (TODO: add --force option)")
	}
	return c.bootstrapDownloader.DownloadToPath(c.state.dataPath)
}

////////////////////////////////////////////////////////////////////////////////
