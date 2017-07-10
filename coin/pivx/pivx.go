// Package pivx is a pivx specific implementation of the masternode monitor tool.
package pivx

////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/sabhiram/gomn/coin"
)

////////////////////////////////////////////////////////////////////////////////

var ()

////////////////////////////////////////////////////////////////////////////////

// PIVX represents a collection of PIVX specific metadata (opaque from the
// coin).
type PIVX struct {
}

// GetPIVX returns a valid instance of the opaque pivx object from the coin.
func GetPIVX(c *coin.Coin) (*PIVX, error) {
	p, ok := c.GetOpaque().(*PIVX)
	if !ok {
		return nil, errors.New("unable to get opaque object for pivx")
	}
	return p, nil
}

////////////////////////////////////////////////////////////////////////////////

func info(c *coin.Coin, binDir, dataDir string, args []string) error {
	c.PrintCoinInfo(binDir, dataDir, "Info for PIVX:")
	p, err := GetPIVX(c)
	if err != nil {
		return err
	}

	_ = p // p is the opaque part of the coin which stores PIVX info
	return nil
}

func download(c *coin.Coin, binDir, dataDir string, args []string) error {
	fmt.Printf("Got download for PIVX (%s, %s)\n", binDir, dataDir)
	fmt.Printf("  with coin: %#v\n", c)
	fmt.Printf("  args: %#v\n", args)
	return nil
}

func bootstrap(c *coin.Coin, binDir, dataDir string, args []string) error {
	fmt.Printf("Got bootstrap for PIVX (%s, %s)\n", binDir, dataDir)
	fmt.Printf("  with coin: %#v\n", c)
	fmt.Printf("  args: %#v\n", args)
	return nil
}

func configure(c *coin.Coin, binDir, dataDir string, args []string) error {
	fmt.Printf("Got configure for PIVX (%s, %s)\n", binDir, dataDir)
	fmt.Printf("  with coin: %#v\n", c)
	fmt.Printf("  args: %#v\n", args)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// Automatically register pivx with gomn if it is included.
func init() {
	// Register the coin and any relevant functions.
	coin.RegisterCoin(
		////////////////////////////////////////////////////////////
		// Register coin constants.
		"pivx",                                 // Name of the coin
		"pivxd",                                // Daemon binaries
		"pivx-cli",                             // Status binaries
		filepath.Join(coin.HomeDir(), "pivx"),  // Default wallet download path
		filepath.Join(coin.HomeDir(), ".pivx"), // Default data path

		////////////////////////////////////////////////////////////
		// Wallet download utility, if version is not set, this is a no-op.
		// Additionally, if the shasum is not set, it will not be checked.
		&coin.WalletDownloader{
			Version:         "2.2.1",
			DownloadURL:     "https://github.com/PIVX-Project/PIVX/releases/download/v2.2.1/pivx-2.2.1-x86_64-linux-gnu.tar.gz",
			CompressionType: "tar.gz",
			Sha256sum:       "401e238e1989b2efdc6d2ac0af3944f1277b2807f79319ad1366248e870e8fcf",
		},

		////////////////////////////////////////////////////////////
		// Bootstrap download utility. If the URL is not set, this is a no-op.
		&coin.BootstrapDownloader{
			DownloadURL:     "https://github.com/PIVX-Project/PIVX/releases/download/v2.2.1/pivx-chain-684000-bootstrap.dat.zip",
			CompressionType: "zip",
		},

		////////////////////////////////////////////////////////////
		// Register coin functions.
		map[string]coin.CoinFunc{
			"info":      info,
			"download":  download,
			"bootstrap": bootstrap,
			"configure": configure,
		},

		////////////////////////////////////////////////////////////
		// Opaque interface for coin
		&PIVX{})
}

////////////////////////////////////////////////////////////////////////////////
