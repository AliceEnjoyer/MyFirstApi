package random_test

import (
	"testing"

	"github.com/AliceEnjoyer/MyFirstApi/internal/lib/random"
)

func TestNewRandomAlias(t *testing.T) {
	const aliasLenght = 6
	set := make(map[string]bool)
	for i := 0; i < 10000; i++ {
		buf := random.NewRandomAlias(aliasLenght)
		if set[buf] {
			t.Error("Two same aliases")
		}
		set[buf] = true
	}
}
