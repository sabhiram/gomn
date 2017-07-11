package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/sabhiram/gomn/version"
)

////////////////////////////////////////////////////////////////////////////////

const (
	confHeaderFmt = `# Warning: This is an auto-generated file! Do not hand-edit!
#   Generated using gomn version %s on %s
# Warning: This is an auto-generated file! Do not hand-edit!

`
)

////////////////////////////////////////////////////////////////////////////////

// CreateConfFile generates a config file with the specified key-value pairs in
// `m`.
func NewConfFile(fp string, m map[string]string) error {
	data := fmt.Sprintf(confHeaderFmt, version.VersionString, time.Now().String())
	for k, v := range m {
		data += fmt.Sprintf("%s=%s\n", k, v)
	}
	data += "\n"
	return ioutil.WriteFile(fp, []byte(data), 0644)
}

// Returns a random hex string of `size`.
func GetRandomHex(size int) string {
	bs := make([]byte, size)
	rand.Read(bs)
	return hex.EncodeToString(bs)
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	rand.Seed(time.Now().UnixNano())
}

////////////////////////////////////////////////////////////////////////////////
