package xenvcom

// BinCommand is the command used inside generated shell hooks.
//
// In the embedded Kite CLI it should be "kite xenv". In the standalone xenv
// binary it should be "xenv".
var BinCommand = "kite xenv"

// BinName is the executable name used for availability checks in shell hooks.
var BinName = "kite"

// SetBinCommand sets the command used inside generated shell hooks.
func SetBinCommand(command string) {
	if command != "" {
		BinCommand = command
	}
}

// SetBinName sets the executable name used for availability checks.
func SetBinName(name string) {
	if name != "" {
		BinName = name
	}
}
