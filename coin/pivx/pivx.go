// Package pivx is a pivx specific implementation of the masternode monitor tool.
package pivx

////////////////////////////////////////////////////////////////////////////////

import (
	"log"
	"os/user"
	"path/filepath"

	"github.com/sabhiram/gomn/coin"
)

////////////////////////////////////////////////////////////////////////////////

var ()

////////////////////////////////////////////////////////////////////////////////

//  basePath      string // path to the coins binaries
// 	dataPath      string // path to data-directory for the coin
// 	daemonBinName string // name of the coin daemon
// 	queryBinName  string // name of the rpc interace binary

////////////////////////////////////////////////////////////////////////////////

func bootstrap(c *coin.Coin, binDir, dataDir string) error {
	log.Printf("Got Bootstrap for PIVX (%s, %s)\n", binDir, dataDir)
	log.Printf("  with coin: %#v\n", c)

	c.DownloadWallet(binDir, dataDir)

	return nil
}

////////////////////////////////////////////////////////////////////////////////

// homeDir gets the user's home directory
func homeDir() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.HomeDir
}

// Automatically register pivx with gomn if it is included.
func init() {
	// Downloader for the wallet.
	walletDownloader := coin.NewWalletDownloader(
		"2.2.1",
		"https://github.com/PIVX-Project/PIVX/releases/download/v2.2.1/pivx-2.2.1-x86_64-linux-gnu.tar.gz",
		"tar.gz",
		"401e238e1989b2efdc6d2ac0af3944f1277b2807f79319ad1366248e870e8fcf")

	bootstrapDownloader := coin.NewBootstrapDownloader(
		"https://github.com/PIVX-Project/PIVX/releases/download/v2.2.1/pivx-chain-684000-bootstrap.dat.zip",
		"zip")

	// Register the coin and any relevant functions.
	coin.RegisterCoin(
		// Register coin constants
		"pivx",                            // Name of the coin
		"pivxd",                           // Daemon binaries
		"pivx-cli",                        // Status binaries
		filepath.Join(homeDir(), "pivx"),  // Default binary path
		filepath.Join(homeDir(), ".pivx"), // Default data path

		// Register wallet / bootstrap fetchers
		walletDownloader,
		bootstrapDownloader,

		// Register coin functions
		bootstrap) // Bootstrap function

}

////////////////////////////////////////////////////////////////////////////////
