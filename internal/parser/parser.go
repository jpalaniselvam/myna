package parser

import (
	"github.com/pelletier/go-toml/v2"
)

// Unmarshal parses the TOML-encoded data and stores the result in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

// Marshal returns the TOML encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	return toml.Marshal(v)
}
