package cryptonight
import (
	"crypto/sha512"
	"hash"
	"sync"
	"unsafe"

	"github.com/ngin-network/cryptonight-go/groestl"
	"github.com/ngin-network/cryptonight-go/jh"

	"github.com/aead/skein"
	"github.com/dchest/blake256"
	"github.com/jzelinskie/whirlpool"
	"github.com/phoreproject/go-x11/bmw"
	"github.com/phoreproject/go-x11/cubed"
	"github.com/phoreproject/go-x11/echo"
	"github.com/phoreproject/go-x11/luffa"
	"github.com/phoreproject/go-x11/shavite"
	"github.com/phoreproject/go-x11/simd"
)

// make ASIC-resistant
var hashPool = [...]*sync.Pool{
	{New: func() interface{} { return blake256.New() }},
	{New: func() interface{} { return bmw.New() }},
	{New: func() interface{} { return cubed.New() }},
	{New: func() interface{} { return echo.New() }},
	{New: func() interface{} { return groestl.New256() }},
	{New: func() interface{} { return jh.New256() }},
	{New: func() interface{} { return luffa.New() }},
	{New: func() interface{} { return skein.New256(nil) }},
	{New: func() interface{} { return simd.New() }},
	{New: func() interface{} { return shavite.New() }},
	{New: func() interface{} { return sha512.New() }},
	{New: func() interface{} { return whirlpool.New() }},
}

func (cc *cache) finalHash() []byte {
	sum := (*[200]byte)(unsafe.Pointer(&cc.finalState))[:]
	for i := uint64(0); i < (cc.finalState[0] & 0x0b); i++ {
		sum = finalPod(sum[:])
	}

	return sum
}

func finalPod(b []byte) []byte {
	var dst [32]byte
	hp := hashPool[b[0]&0x0b]
	h := hp.Get().(hash.Hash)
	h.Reset()
	h.Write(b)
	sum := h.Sum(nil)
	hp.Put(h)

	copy(dst[:], sum)
	return sum
}

