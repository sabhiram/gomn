// Package types holds a bunch of commonly passed around structures.
package types

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
