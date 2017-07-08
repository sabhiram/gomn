package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////

// ProgressTracker implements the Write function so we can have it track the number
// of bytes being written.
type ProgressTracker struct {
	total int64 // Number of total bytes
	count int64 // Number of bytes seen in Reader
}

// NewProgressTracker returns a instance that implements a writer interface so we
// can track progress of a download via a tee reader in the caller.
func NewProgressTracker(total int64) *ProgressTracker {
	return &ProgressTracker{
		total: total,
		count: 0,
	}
}

// Write implements the Write function needed to satisfy the Writer interface.
func (pt *ProgressTracker) Write(bs []byte) (int, error) {
	n := len(bs)
	pt.count += int64(n)
	percent := (pt.count * 10000) / pt.total
	percr := percent / 100
	percf := percent % 100
	fmt.Printf("%02d.%02d%% Done -- %d/%d\r", percr, percf, pt.count, pt.total)
	return n, nil
}

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

	tee := io.TeeReader(resp.Body, NewProgressTracker(resp.ContentLength))
	_, err = io.Copy(out, tee)
	return err
}

////////////////////////////////////////////////////////////////////////////////

// extractTarGzip extracts a given source file path into a destination path
// provided that the input is a valid tar.gz file.
func extractTarGzip(srcfp, dstdp string) error {
	log.Printf("  Extracting .tar.gz file into %s\n", dstdp)

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

	if err := os.MkdirAll(dstdp, 0755); err != nil {
		return err
	}

	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil // Done!
		case err != nil:
			return err // Actual error, bad news
		case hdr == nil:
			continue // edge case, unsure if this is valid
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

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}

// extractZip extracts a given source file path into a destination path
// provided that the input is a valid zip file.
func extractZip(srcfp, dstdp string) error {
	log.Printf("  Extracting .zip file into %s\n", dstdp)

	r, err := zip.OpenReader(srcfp)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dstdp, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dstdp, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 755)
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				return err
			}

			fout, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer fout.Close()

			_, err = io.Copy(fout, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// extractToPath extracts the given type of compressed file specified in `srcfp`
// to `dstdp`.
func extractToPath(ctype, srcfp, dstdp string) error {
	switch strings.ToLower(ctype) {
	case "tar.gz":
		return extractTarGzip(srcfp, dstdp)
	case "zip":
		return extractZip(srcfp, dstdp)
	case "", "none":
		log.Printf("No compression type specified, need to move downloaded file!\n")
		// return os.Rename(srcfp, dstdp)
		return nil
	default:
		return fmt.Errorf("unsupported compression type (%s)", ctype)
	}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
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

// DownloadToPath grabs the underlying wallet file, and checks its sha256sum
// to verify that it is indeed the expected file. If so, it extracts the
// contents to the appropriate
func (w *WalletDownloader) DownloadToPath(walletPath string) error {
	log.Printf("  Fetching wallet from %s into %s\n", w.downloadURL, walletPath)

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
			log.Printf("Warning: Unable to cleanup temp file: %s\n", f)
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
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type BootstrapDownloader struct {
	downloadURL     string // URL to fetch bootstrap archive
	compressionType string // type of compression ["tar.gz", "zip", "none"]
}

// NewBootstrapDownloader returns a new instance of a bootstrap downloader.
func NewBootstrapDownloader(url, ctype string) *BootstrapDownloader {
	return &BootstrapDownloader{
		downloadURL:     url,
		compressionType: ctype,
	}
}

// DownloadToPath grabs a archive from a web url defined in `b` and extracts
// the file if needed into `bootstrapPath`.
func (b *BootstrapDownloader) DownloadToPath(bootstrapPath string) error {
	log.Printf("  Fetching bootstrap from %s into %s\n", b.downloadURL, bootstrapPath)

	// Fetch the file into a temporary file
	tempFile := filepath.Join(os.TempDir(), "bootstrapdl")

	// Try to fetch the wallet to the temporary file
	if err := downloadURLToPath(b.downloadURL, tempFile); err != nil {
		return err
	}

	// Since we have created a file now, make sure the tempfile is removed
	// regardless of if this function succeeds.
	defer func(f string) {
		if err := os.Remove(f); err != nil {
			log.Printf("Warning: Unable to cleanup temp file: %s\n", f)
		}
	}(tempFile)

	// Extract the file to the specified path, we assume that the type of file
	// is specified at the tail end of the URL.
	return extractToPath(b.compressionType, tempFile, bootstrapPath)

}

////////////////////////////////////////////////////////////////////////////////
