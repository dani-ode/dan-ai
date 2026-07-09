// pkg/ulid/ulid.go
package ulid

import (
	"github.com/oklog/ulid/v2"
)

// New generates a cryptographically secure, random ULID string.
func New() string {
	return ulid.Make().String()
}
