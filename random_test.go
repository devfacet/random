// Random
// For the full copyright and license information, please view the LICENSE.txt file.

package random_test

import (
	"sync"
	"testing"

	"github.com/devfacet/random"
)

func TestMustULID(t *testing.T) {
	u := random.MustULID(true)
	if u == "" {
		t.Errorf("failed to create ULID (secure=true)")
	}
	u = random.MustULID(false)
	if u == "" {
		t.Errorf("failed to create ULID (secure=false)")
	}
}

func TestMustULIDDuplicate(t *testing.T) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	limit := 1000

	ul1 := make(map[string]int, limit)
	dups1 := 0
	for i := 0; i < limit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			u := random.MustULID(true)
			mu.Lock()
			ul1[u]++
			if ul1[u] > 1 {
				dups1++
			}
			mu.Unlock()
		}()
	}
	if dups1 > 0 {
		t.Errorf("failed to create unique ULIDs (secure=true): %d duplicates", dups1)
	}

	ul2 := make(map[string]int, limit)
	dups2 := 0
	for i := 0; i < limit; i++ {
		wg.Add(1)
		go func() {
			u := random.MustULID(false)
			mu.Lock()
			ul2[u]++
			if ul2[u] > 1 {
				dups2++
			}
			mu.Unlock()
			wg.Done()
		}()
	}
	if dups2 > 0 {
		t.Errorf("failed to create unique ULIDs (secure=false): %d duplicates", dups2)
	}

	wg.Wait()
}

func BenchmarkMustULID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		random.MustULID(true)
	}
}

func BenchmarkMustULIDInsecure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		random.MustULID(false)
	}
}

func TestMustBytes(t *testing.T) {
	b := random.MustBytes(32, true)
	if l := len(b); l != 32 {
		t.Errorf("failed to create random bytes (secure=true): %d bytes", l)
	}
	b = random.MustBytes(32, false)
	if l := len(b); l != 32 {
		t.Errorf("failed to create random bytes (secure=false): %d bytes", l)
	}
}

func BenchmarkMustBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		random.MustBytes(32, true)
	}
}

func BenchmarkMustBytesInsecure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		random.MustBytes(32, false)
	}
}

func TestMillisecond(t *testing.T) {
	d := random.Millisecond(1, 5000)
	if d.Milliseconds() < 1 || d.Milliseconds() > 5000 {
		t.Errorf("failed to create correct duration: %s)", d)
	}
}
