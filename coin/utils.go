package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////

// WalletDownloader is a per-coin wallet fetcher.
type WalletDownloader struct {
	version     string // version of the wallet
	downloadURL string // url to fetch the wallet
	sha256sum   string // shasum for the download
}

// NewWalletDownloader returns a new instance of a wallet downloader.
func NewWalletDownloader(v, url, sha string) *WalletDownloader {
	return &WalletDownloader{
		version:     v,
		downloadURL: url,
		sha256sum:   sha,
	}
}

func (w *WalletDownloader) Log() {
	fmt.Printf("WalletDownloader::%#v\n", w)
}

////////////////////////////////////////////////////////////////////////////////

type BootstrapDownloader struct {
	downloadURL string // URL to fetch bootstrap archive
}

// NewBootstrapDownloader returns a new instance of a bootstrap downloader.
func NewBootstrapDownloader(url string) *BootstrapDownloader {
	return &BootstrapDownloader{
		downloadURL: url,
	}
}

func (b *BootstrapDownloader) Log() {
	fmt.Printf("BootstrapDownloader::%#v\n", b)
}

////////////////////////////////////////////////////////////////////////////////
