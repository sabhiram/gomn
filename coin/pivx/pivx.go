// Package pivx is a pivx specific implementation of the masternode monitor tool.
package pivx

////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/sabhiram/gomn/coin"
	"github.com/sabhiram/gomn/types"
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

	// Parse command arguments
	fs := flag.NewFlagSet("pivx-download", flag.ContinueOnError)
	cargs := &types.Download{}
	fs.StringVar(&cargs.URL, "url", "", "override the wallet download URL")
	fs.StringVar(&cargs.Type, "type", "", "override the wallet download type (compression)")
	fs.StringVar(&cargs.ShaSum, "shasum", "", "override the wallet's sha256sum")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return c.DownloadWallet(args, cargs)
}

func bootstrap(c *coin.Coin, args []string) error {
	fmt.Printf("Attempting to bootstrap PIVX data into %s\n", c.GetDataPath())

	// Parse command arguments
	fs := flag.NewFlagSet("pivx-bootstrap", flag.ContinueOnError)
	cargs := &types.Bootstrap{}
	fs.StringVar(&cargs.URL, "url", "", "override the bootstrap URL")
	fs.StringVar(&cargs.Type, "type", "", "override the bootstrap type (compression)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return c.DownloadBootstrap(args, cargs)
}

func configure(c *coin.Coin, args []string) error {
	fmt.Printf("Attempting to configure %s\n", c.GetConfFilePath())

	// Parse command arguments
	fs := flag.NewFlagSet("pivx-configure", flag.ContinueOnError)
	cargs := &types.Configure{}
	fs.StringVar(&cargs.IP, "ip", "", "masternode's fixed IP (required)")
	fs.StringVar(&cargs.MnPK, "mnpkey", "", "masternode's private key (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: Validate that this is an IPV4/V6 IP
	if len(cargs.IP) == 0 {
		return errors.New("np IP specified (use '--ip A.B.C.D'")
	}
	if len(cargs.MnPK) == 0 {
		return errors.New("no masternode privatekey specified")
	}

	return coin.CreateConfFile(c.GetConfFilePath(), map[string]string{
		"rpcuser":            coin.GetRandomHex(32),
		"rpcpassword":        coin.GetRandomHex(64),
		"rpcallowip":         "127.0.0.1",
		"listen":             "1",
		"server":             "1",
		"daemon":             "1",
		"#masternode":        "1",
		"maxconnections":     "256",
		"bind":               "0.0.0.0",
		"externalip":         cargs.IP,
		"masternodeaddr":     fmt.Sprintf("%s:%d", cargs.IP, c.GetPort()),
		"#masternodeprivkey": cargs.MnPK,
	})
}

func getinfo(c *coin.Coin, args []string) error {
	rsp, err := c.DoJSONRPCCommand("getinfo", nil)
	if err != nil {
		return err
	}

	if rsp.Error.Code != 0 {
		switch rsp.Error.Code {
		case -28:
			fmt.Printf("pivxd starting up -- %s\n", rsp.Error.Message)
		default:
			return fmt.Errorf("RPC error (%d) : %s\n", rsp.Error.Code, rsp.Error.Message)
		}
	}

	// TODO: Return a coin-generic status along with the various info pieces.
	fmt.Printf("GOT RESPONSE: %#v\n", rsp)
	return err
}

////////////////////////////////////////////////////////////////////////////////

// Automatically register pivx with gomn if it is included.
func init() {
	// Register the coin and any relevant functions.
	err := coin.RegisterCoin(
		////////////////////////////////////////////////////////////
		// Register coin constants.
		"pivx",      // Name of the coin
		51472,       // PIVX port
		51473,       // RPC port
		"pivxd",     // Daemon binaries
		"pivx-cli",  // Status binaries
		"pivx.conf", // Coin config file

		filepath.Join(coin.HomeDir(), "pivx"),  // Base path to wallet dl
		filepath.Join("pivx-2.2.1", "bin"),     // Subpath to binaries
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
			DownloadURL:     "https://github.com/PIVX-Project/PIVX/releases/download/v2.2.1/pivx-chain-721000-bootstrap.dat.zip",
			CompressionType: "zip",
		},

		////////////////////////////////////////////////////////////
		// Register required coin functions.
		&coin.FunctionMap{
			InfoFn:      info,
			DownloadFn:  download,
			BootstrapFn: bootstrap,
			ConfigureFn: configure,
			GetInfoFn:   getinfo,
		},

		////////////////////////////////////////////////////////////
		// Opaque interface for coin.
		&PIVX{})
	if err != nil {
		panic(err.Error())
	}
}

////////////////////////////////////////////////////////////////////////////////
