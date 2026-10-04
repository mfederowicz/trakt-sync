// Package cfg used for process configuration
package cfg

// BoolText is a "true" / "false" config value that the config file may hold quoted or as a TOML boolean.
type BoolText string

// UnmarshalText keeps the value as text; the TOML decoder passes a boolean as "true" or "false".
func (b *BoolText) UnmarshalText(text []byte) error {
	*b = BoolText(text)
	return nil
}
