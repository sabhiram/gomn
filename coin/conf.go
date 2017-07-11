package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/sabhiram/gomn/version"
)

////////////////////////////////////////////////////////////////////////////////

// Returns a random hex string of `size`.
func GetRandomHex(size int) string {
	bs := make([]byte, size)
	rand.Read(bs)
	return hex.EncodeToString(bs)
}

////////////////////////////////////////////////////////////////////////////////

// CreateConfFile generates a config file with the specified key-value pairs in
// `m`.
func CreateConfFile(fp string, m map[string]string) error {
	data := fmt.Sprintf(`# Warning: This is an auto-generated file! Do not hand-edit!
#   Generated using gomn version %s on %s
# Warning: This is an auto-generated file! Do not hand-edit!

`, version.VersionString, time.Now().String())
	for k, v := range m {
		data += fmt.Sprintf("%s=%s\n", k, v)
	}
	data += "\n"
	return ioutil.WriteFile(fp, []byte(data), 0644)
}

// LoadConfFile returns a map of key-value pairs found in a `.conf` file pointed
// to by `fp`.
func LoadConfFile(fp string) (map[string]string, error) {
	bs, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	m := map[string]string{}
	for _, line := range strings.Split(string(bs), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case len(line) == 0:
			break
		case line[0] == '#':
			break
		default:
			idx := strings.IndexByte(line, '=')
			if idx >= 0 {
				m[line[:idx]] = line[idx+1:]
			}
		}
	}
	return m, nil
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	rand.Seed(time.Now().UnixNano())
}

////////////////////////////////////////////////////////////////////////////////
