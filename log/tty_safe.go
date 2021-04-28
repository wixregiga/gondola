// +build appengine

package log

import (
	"io"

	"gondola/internal"
)

func isatty(w io.Writer) bool {
	if internal.InAppEngineDevServer() {
		return true
	}
	return false
}
