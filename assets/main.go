package assets

import _ "embed"

//go:embed disclaimer.html
var Disclaimer []byte

//go:embed blacklist.html
var Blacklist []byte

//go:embed whitelist.html
var Whitelist []byte
