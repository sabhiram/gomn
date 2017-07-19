// Package types holds a bunch of commonly passed around structures.
package types

////////////////////////////////////////////////////////////////////////////////

// CLI represents all command line arguments that the `gomn` application can
// receive directly. Sub-command specific arguments are encapsulated in the
// structures below.
type CLI struct {
	Coin     string   // Name of the coin we are operating on
	Wallet   string   // base path to where wallet binaries will be extracted  (empty => coin default)
	BinPath  string   // sub-Path to coin's binary directory (empty => coin default)
	DataPath string   // Path to coin's data directory (empty => coin default)
	Args     []string // Rest of the command line, args[0] is the command
}

////////////////////////////////////////////////////////////////////////////////

// Download represents the arguments passed to the "download" command.
type Download struct {
	URL    string
	Type   string
	ShaSum string
}

// Bootstrap represents the arguments passed to the "download" command.
type Bootstrap struct {
	URL  string
	Type string
}

// Configure represents the arguments passed to the "download" command.
type Configure struct {
	IP   string
	MnPK string
}

////////////////////////////////////////////////////////////////////////////////
