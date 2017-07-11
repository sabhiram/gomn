package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"os"
	"os/user"
)

////////////////////////////////////////////////////////////////////////////////

// HomeDir gets the user's home directory
func HomeDir() string {
	if u, err := user.Current(); err == nil {
		return u.HomeDir
	}
	return ""
}

// DirExists returns true if `fp` is a directory and exists.
func DirExists(fp string) bool {
	if st, err := os.Stat(fp); err == nil && st.IsDir() {
		return true
	}
	return false
}

// FileExists returns true if `fp` is a file and exists.
func FileExists(fp string) bool {
	if st, err := os.Stat(fp); err == nil && !st.IsDir() {
		return true
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
