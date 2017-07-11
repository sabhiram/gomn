// Package pivx is a pivx specific implementation of the masternode monitor tool.
package pivx

////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/sabhiram/gomn/coin"
)

////////////////////////////////////////////////////////////////////////////////

var (
	CLI = struct {
		ip   string
		mnPK string
	}{}
)

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

func info(c *coin.Coin, args []string) error {
	c.PrintCoinInfo("Info for PIVX:")
	p, err := GetPIVX(c)
	if err != nil {
		return err
	}

	_ = p // p is the opaque part of the coin which stores PIVX info
	return nil
}

func download(c *coin.Coin, args []string) error {
	fmt.Printf("Attempting to download PIVX wallet into %s\n", c.GetBinPath())
	return c.DownloadWallet()
}

func bootstrap(c *coin.Coin, args []string) error {
	fmt.Printf("Attempting to download PIVX data into %s\n", c.GetDataPath())
	return c.DownloadBootstrap()
}

func configure(c *coin.Coin, args []string) error {
	fmt.Printf("Attempting to configure %s\n", c.GetConfFilePath())
	fs := flag.NewFlagSet("pivx-configure", flag.ContinueOnError)
	fs.StringVar(&CLI.ip, "ip", "", "masternode's fixed IP (required)")
	fs.StringVar(&CLI.mnPK, "mnpkey", "", "masternode's private key (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: Validate that this is an IPV4/V6 IP
	if len(CLI.ip) == 0 {
		return errors.New("np IP specified (use '--ip A.B.C.D'")
	}
	if len(CLI.mnPK) == 0 {
		return errors.New("no masternode private-key specified")
	}

	return coin.NewConfFile(c.GetConfFilePath(), map[string]string{
		"rpcuser":            coin.GetRandomHex(32),
		"rpcpassword":        coin.GetRandomHex(64),
		"rpcallowip":         "127.0.0.1",
		"listen":             "1",
		"server":             "1",
		"daemon":             "1",
		"#masternode":        "1",
		"maxconnections":     "256",
		"bind":               "0.0.0.0",
		"externalip":         CLI.ip,
		"masternodeaddr":     fmt.Sprintf("%s:%d", CLI.ip, c.GetPort()),
		"#masternodeprivkey": CLI.mnPK,
	})
}

////////////////////////////////////////////////////////////////////////////////

// Automatically register pivx with gomn if it is included.
func init() {
	// Register the coin and any relevant functions.
	coin.RegisterCoin(
		////////////////////////////////////////////////////////////
		// Register coin constants.
		"pivx",      // Name of the coin
		51472,       // PIVX port
		"pivxd",     // Daemon binaries
		"pivx-cli",  // Status binaries
		"pivx.conf", // Coin config file

		filepath.Join(coin.HomeDir(), "pivx", "pivx-2.2.1", "bin"), // Default binary path
		filepath.Join(coin.HomeDir(), ".pivx"),                     // Default data path

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
			DownloadURL:     "https://github.com/PIVX-Project/PIVX/releases/download/v2.2.1/pivx-chain-721000-bootstrap.dat.zip",
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
