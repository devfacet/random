// Random
// For the full copyright and license information, please view the LICENSE.txt file.

// Package random provides functions for generating random values.
package random

import (
	"crypto/rand"
	mrand "math/rand"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	mr = mrandReader{}
)

func init() {
	mrand.Seed(time.Now().UnixNano())
}

// MustULID returns a ULID or panics on error.
func MustULID(secure bool) string {
	// Ref: ulid.Monotonic: The returned type isn't safe for concurrent use.
	var u ulid.ULID
	if secure {
		entropy := ulid.Monotonic(rand.Reader, 0) // incremental entropy for same ms (cryptographic)
		u = ulid.MustNew(ulid.Now(), entropy)
	} else {
		// Ref: math/rand: The default Source is safe for concurrent use by multiple goroutines, but Sources created by NewSource are not.
		//entropy := ulid.Monotonic(mrand.New(mrand.NewSource(time.Now().UnixNano())), 0) // incremental entropy for same ms (random)
		entropy := ulid.Monotonic(mr, 0) // incremental entropy for same ms (random)
		u = ulid.MustNew(ulid.Now(), entropy)
	}
	return strings.ToLower(u.String())
}

// MustBytes generates random bytes or panics on error.
func MustBytes(n int, secure bool) []byte {
	b := make([]byte, n)
	if secure {
		if _, err := rand.Read(b); err != nil {
			panic(err)
		}
	} else {
		if _, err := mrand.Read(b); err != nil {
			panic(err)
		}
	}
	return b
}

// Millisecond returns a random amount of millisecond by the given range.
func Millisecond(min, max int64) time.Duration {
	// Uses default source in the math/rand package (safe for concurrent use by multiple goroutines)
	return time.Duration(mrand.Int63n(max-min+1)+min) * time.Millisecond
}

// mrandReader provides reader for math/rand package.
type mrandReader struct {
}

// Read generates random bytes from the math/rand package default Source.
func (r mrandReader) Read(buf []byte) (n int, err error) {
	// Ref: math/rand: The default Source is safe for concurrent use by multiple goroutines, but Sources created by NewSource are not.
	return mrand.Read(buf)
}
