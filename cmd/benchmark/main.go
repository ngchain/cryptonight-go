package main // import "github.com/ngin-network/cryptonight-go/cmd/cnhash"

import (
	"fmt"
	"github.com/ngin-network/cryptonight-go"
	"runtime"
	"sync/atomic"
	"time"
)

func main() {
	t := runtime.GOMAXPROCS(0)
	fmt.Println("GOMAXPROCS =", t)

	fmt.Println("Start test hashing")

	hashes := uint64(0)
	lastSnap := hashes
	benchData := []byte("NGIN TESTNET")

	for i := 0; i < t; i++ {
		go func() {
			for j := 0; true; j++ {
				CNgo(benchData)
				atomic.AddUint64(&hashes, 1)
			}
		}()
	}

	i := uint64(0)
	for range time.Tick(5 * time.Second) {
		i++
		snap := atomic.LoadUint64(&hashes)
		fmt.Printf("%.2f H/s  (%.2f H/s)\n", float64(snap-lastSnap)/5, float64(snap)/float64(5*i))
		lastSnap = snap
	}

}

func CNgo(blob []byte) []byte {
	var result []byte
	result = cryptonight.Sum(blob, 0)

	return result
}
