package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////

// downloadURLToPath fetches a file specified at `url` to `filepath`.
func downloadURLToPath(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

////////////////////////////////////////////////////////////////////////////////

// extractTarGzip extracts a given source file path into a destination path
// provided that the input is a valid tar.gz file.
func extractTarGzip(srcfp, dstdp string) error {
	srcf, err := os.Open(srcfp)
	if err != nil {
		return err
	}
	defer srcf.Close()

	r, err := gzip.NewReader(srcf)
	if err != nil {
		return err
	}
	defer r.Close()

	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		switch {
		case hdr == nil:
			continue // edge case, unsure if this is valid
		case err != nil:
			return err // Actual error, bad news
		case err == io.EOF:
			return nil // Done!
		}

		targetFile := filepath.Join(dstdp, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(targetFile); err != nil {
				if err := os.MkdirAll(targetFile, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(targetFile, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			defer f.Close()
		}
	}
}

// extractToPath extracts the given type of compressed file specified in `srcfp`
// to `dstdp`.
func extractToPath(ctype, srcfp, dstdp string) error {
	switch strings.ToLower(ctype) {
	case "tar.gz":
		return extractTarGzip(srcfp, dstdp)
	case "", "none":
		fmt.Printf("No compression type specified, need to move downloaded file!\n")
		// return os.Rename(srcfp, dstdp)
		return nil
	default:
		return fmt.Errorf("unsupported compression type (%s)", ctype)
	}
}

////////////////////////////////////////////////////////////////////////////////

// WalletDownloader is a per-coin wallet fetcher.
type WalletDownloader struct {
	version         string // version of the wallet
	downloadURL     string // url to fetch the wallet
	compressionType string // type of compression ["tar.gz", "zip", "none"]
	sha256sum       string // shasum for the download
}

// NewWalletDownloader returns a new instance of a wallet downloader.
func NewWalletDownloader(v, url, ct, sha string) *WalletDownloader {
	return &WalletDownloader{
		version:         v,
		downloadURL:     url,
		compressionType: ct,
		sha256sum:       sha,
	}
}

// FetchToPath grabs the underlying wallet file, and checks its sha256sum
// to verify that it is indeed the expected file. If so, it extracts the
// contents to the appropriate
func (w *WalletDownloader) FetchWalletToPath(walletPath string) error {
	// Fetch the file into a temporary file
	tempFile := filepath.Join(os.TempDir(), "walletdl")

	// Try to fetch the wallet to the temporary file
	if err := downloadURLToPath(w.downloadURL, tempFile); err != nil {
		return err
	}

	// Since we have created a file now, make sure the tempfile is removed
	// regardless of if this function succeeds.
	defer func(f string) {
		if err := os.Remove(f); err != nil {
			fmt.Printf("Warning: Unable to cleanup temp file: %s\n", f)
		}
	}(tempFile)

	// Verify sha256 sum
	bs, err := ioutil.ReadFile(tempFile)
	if err != nil {
		return err
	}

	t := []byte{}
	for _, b := range sha256.Sum256(bs) {
		t = append(t, b)
	}
	shasum := hex.EncodeToString(t)
	if strings.ToLower(shasum) != strings.ToLower(w.sha256sum) {
		return fmt.Errorf("shasum for download (%s) does not match expected (%s)", shasum, w.sha256sum)
	}

	// Extract the file to the specified path, we assume that the type of file
	// is specified at the tail end of the URL.
	return extractToPath(w.compressionType, tempFile, walletPath)
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
