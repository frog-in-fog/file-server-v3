package igf2

import (
	"hash"
	"io"

	"file-server-v3/pkg/ntru/mgf1"
)

// IGF2 implements the IGF2 Index Generation Function defined in the X9.92 spec
// for NTRUEncrypt.
type IGF2 struct {
	maxValue        int16
	bitsPerIndex    int16
	leftoverBits    int
	numLeftoverBits int
	cutoff          int
	source          io.ByteReader
}

// NextIndex derives the next index.
func (g *IGF2) NextIndex() (int16, error) {
	ret := 0
	for {
		// Make sure leftoverBits has at least bitsPerIndex in it.
		for g.numLeftoverBits < int(g.bitsPerIndex) {
			g.leftoverBits <<= 8
			c, err := g.source.ReadByte()
			if err != nil {
				return 0, err
			}
			g.leftoverBits |= int(c)
			g.numLeftoverBits += 8
		}

		// Pull off bitsPerIndex from leftoverBits.  Store in ret.
		shift := g.numLeftoverBits - int(g.bitsPerIndex)
		ret = 0xffff & (g.leftoverBits >> uint(shift))
		g.numLeftoverBits = shift
		g.leftoverBits &= ((1 << uint(g.numLeftoverBits)) - 1)

		if ret < g.cutoff {
			return int16(ret) % g.maxValue, nil
		}
	}
}

func (g *IGF2) Close() error {
	if closer, ok := (g.source).(io.Closer); ok {
		return closer.Close()
	}
	return nil // Oh well...
}

// New creates an IGF2 driven by a MGF1.
func New(maxValue, bitsPerIndex int16, hashFn func() hash.Hash, minNumRuns int, seed []byte, seedOff, seedLen int) *IGF2 {
	mgf := mgf1.New(hashFn, minNumRuns, true, seed, seedOff, seedLen)
	return NewFromReader(maxValue, bitsPerIndex, mgf)
}

// NewFromReader creates an IGF2 driven by a io.ByteReader.
func NewFromReader(maxValue, bitsPerIndex int16, source io.ByteReader) *IGF2 {
	g := &IGF2{}
	g.maxValue = int16(maxValue)
	g.bitsPerIndex = int16(bitsPerIndex)
	g.source = source
	modulus := 1 << uint(bitsPerIndex)
	g.cutoff = modulus - (modulus % int(maxValue))
	return g
}
