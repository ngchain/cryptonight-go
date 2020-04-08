// Package cryptonight implements CryptoNight hash function and some of its
// variant. Original CryptoNight algorithm is defined in CNS008 at
// https://cryptonote.org/cns/cns008.txt
package cryptonight // import "github.com/ngchain/cryptonight-go"

import (
	"sync"
)

type cache struct {
	// DO NOT change the order of these fields in this struct!
	// They are carefully placed in this order to keep at least 16-byte aligned
	// for some fields.
	//
	// In the future the alignment may be set explicitly, see
	// https://github.com/golang/go/issues/19057

	scratchpad [2 * 1024 * 1024 / 8 / 4]uint64 // 512KB scratchpad for memhard loop
	finalState [25]uint64                      // state of keccak1600
	_          [8]byte                         // padded to keep 16-byte align 524288+200=524488=32781*16-8

	blocks [16]uint64 // temporary chunk/pointer of data
	rkeys  [40]uint32 // 10 rounds, instead of 14 as in standard AES-256
}

// cachePool is a pool of cache.
var cachePool = sync.Pool{
	New: func() interface{} {
		return new(cache)
	},
}

// Sum calculate a CryptoNight hash digest. The return value is exactly 32 bytes
// long.
func Sum(data []byte, variant int) []byte {
	cc := cachePool.Get().(*cache)
	sum := cc.sum(data, variant)
	cachePool.Put(cc)

	return sum
}
