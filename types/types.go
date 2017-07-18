// Package types holds a bunch of commonly passed around structures.
package types

////////////////////////////////////////////////////////////////////////////////

type CLI struct {
	Coin     string   // Name of the coin we are operating on
	Wallet   string   // base path to where wallet binaries will be extracted  (empty => coin default)
	BinPath  string   // sub-Path to coin's binary directory (empty => coin default)
	DataPath string   // Path to coin's data directory (empty => coin default)
	Args     []string // Rest of the command line, args[0] is the command
}

////////////////////////////////////////////////////////////////////////////////

type Download struct {
	URL    string
	Type   string
	ShaSum string
}

type Bootstrap struct {
	URL  string
	Type string
}

type Configure struct {
	IP   string
	MnPK string
}

////////////////////////////////////////////////////////////////////////////////
