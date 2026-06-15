package pneumatic_test

import (
	"iter"
	"slices"
	"strconv"
	"testing"

	pn "github.com/skubalj/pneumatic"
	"github.com/stretchr/testify/require"
)

func TestPairZipping(t *testing.T) {
	numbers := []string{"10", "10.1", "NaN", "abc"}

	pair, ok := pn.Find(
		func(p pn.Pair[float64, error]) bool { return p.Y != nil },
		pn.Map(
			func(s string) pn.Pair[float64, error] { return pn.ZipPair(strconv.ParseFloat(s, 64)) },
			slices.Values(numbers),
		),
	)
	val, err := pair.Unzip()
	require.True(t, ok)
	require.Equal(t, 0.0, val)
	require.Error(t, err)
}

func TestRange(t *testing.T) {
	require.Empty(t, slices.Collect(pn.Range[int](0, 0)))
	require.Equal(t, []int{1, 2, 3}, slices.Collect(pn.Range(1, 4)))
	require.Equal(t, []rune{'a', 'b', 'c', 'd', 'e'}, slices.Collect(pn.RangeInclusive('a', 'e')))
	require.Equal(t, []byte{0x0a, 0x0b, 0x0c, 0x0d}, slices.Collect(pn.Range[byte](0x0a, 0x0e)))
}

func Test_Primes(t *testing.T) {
	require.Equal(t, []int{2, 3, 5, 7, 11, 13, 17, 19}, slices.Collect(primesUnder(20)))
}

func primesUnder(bound int) iter.Seq[int] {
	return pn.Filter(isPrime, pn.Range(2, bound))
}

func isPrime(candidate int) bool {
	return !pn.Any(func(x int) bool { return x == 0 },
		pn.Map(func(prime int) int { return candidate % prime },
			primesUnder(candidate)))
}
