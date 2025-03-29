//go:build windows && divert_embed && (amd64 || 386 || arm64)

package divert

import (
	"embed"
)

//go:embed WinDivert
var f embed.FS
