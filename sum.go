package cryptonight

import (
	"github.com/ngin-network/cryptonight-go/internal/aes"
	"github.com/ngin-network/cryptonight-go/internal/sha3"
	//"encoding/binary"
)

func (cc *cache) sum(data []byte, variant int) []byte {
	//////////////////////////////////////////////////
	// these variables never escape to heap
	var (
		// used in memory hard
		a, b, c, d [2]uint64
	)

	//////////////////////////////////////////////////
	// as per CNS008 sec.3 Scratchpad Initialization
	sha3.Keccak1600State(&cc.finalState, data)

	// scratchpad init
	aes.CnExpandKeyGo(cc.finalState[:4], &cc.rkeys)
	copy(cc.blocks[:], cc.finalState[8:24])

	for i := 0; i < 2*1024*1024/8/4; i += 16 {
		for j := 0; j < 16; j += 2 {
			aes.CnRoundsGo(cc.blocks[j:j+2], cc.blocks[j:j+2], &cc.rkeys)
		}
		copy(cc.scratchpad[i:i+16], cc.blocks[:16])
	}

	//////////////////////////////////////////////////
	// as per CNS008 sec.4 Memory-Hard Loop
	a[0] = cc.finalState[0] ^ cc.finalState[4]
	a[1] = cc.finalState[1] ^ cc.finalState[5]
	b[0] = cc.finalState[2] ^ cc.finalState[6]
	b[1] = cc.finalState[3] ^ cc.finalState[7]

	// Turtle changed here // Ngin
	for i := 0; i < 65536; i++ {
		addr := (a[0] >> 2) & ((((1 << 17) >> 2) - 1) << 4) >> 3
		aes.CnSingleRoundGo(c[:2], cc.scratchpad[addr:addr+2], &a)

		cc.scratchpad[addr+0] = b[0] ^ c[0]
		cc.scratchpad[addr+1] = b[1] ^ c[1]

		addr = (a[0] >> 2) & ((((1 << 17) >> 2) - 1) << 4) >> 3
		d[0] = cc.scratchpad[addr]
		d[1] = cc.scratchpad[addr+1]

		// byteMul
		lo, hi := mul128(c[0], d[0])

		// byteAdd
		a[0] += hi
		a[1] += lo

		cc.scratchpad[addr+0] = a[0]
		cc.scratchpad[addr+1] = a[1]

		a[0] ^= d[0]
		a[1] ^= d[1]

		b[0] = c[0]
		b[1] = c[1]
	}

	//////////////////////////////////////////////////
	// as per CNS008 sec.5 Result Calculation
	aes.CnExpandKeyGo(cc.finalState[4:8], &cc.rkeys)
	tmp := cc.finalState[8:24] // a temp pointer

	for i := 0; i < 2*1024*1024/8/4; i += 16 {
		for j := 0; j < 16; j += 2 {
			cc.scratchpad[i+j+0] ^= tmp[j+0]
			cc.scratchpad[i+j+1] ^= tmp[j+1]
			aes.CnRoundsGo(cc.scratchpad[i+j:i+j+2], cc.scratchpad[i+j:i+j+2], &cc.rkeys)
		}
		tmp = cc.scratchpad[i : i+16]
	}

	copy(cc.finalState[8:24], tmp)
	sha3.Keccak1600Permute(&cc.finalState)

	return cc.finalHash()
}
