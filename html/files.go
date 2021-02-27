//go:generate file2go -p html -T html -o files.go

package html

import (
	"github.com/amery/file2go/html"
)

var Files = html.NewCollection()
