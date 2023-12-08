// Package mgftp1 implements the MGF-TP-1 algoritm for converting a byte stream
// into a sequence of trits.  It implements both the forward direction and the
// reverse.
package mgftp1

import (
	"io"

	"file-server-v3/pkg/ntru/polynomial"
)

// GenTrinomial generates a trinomial of degree N using the MGF-TP-1 algorithm
// to convert the io.ByteReader into trits.
func GenTrinomial(n int, r io.ByteReader) (*polynomial.Full, error) {
	p := polynomial.New(n)

	limit := 5 * (n / 5)
	i := 0
	for i < limit {
		o, err := r.ReadByte()
		if err != nil {
			return nil, err
		} else if o >= 243 {
			continue
		}
		for j := 0; j < 5; j++ {
			b := o % 3
			p.P[i+j] = int16(b)
			o = (o - b) / 3
		}
		i += 5
	}
nLoop:
	for i < n {
		o, err := r.ReadByte()
		if err != nil {
			return nil, err
		} else if o >= 243 {
			continue
		}

		for j := 0; j < 5; j++ {
			b := o % 3
			p.P[i+j] = int16(b)
			o = (o - b) / 3
			if i+j+1 == n {
				break nLoop
			}
		}
		i += 5
	}

	// Renormalize from [0..2] to [-1..1]
	for i := range p.P {
		if p.P[i] == 2 {
			p.P[i] = -1
		}
	}
	return p, nil
}

// EncodeTrinomial generates a byte stream that is the encoding of a trinomial.
func EncodeTrinomial(poly *polynomial.Full, w io.ByteWriter) error {
	n := len(poly.P)
	accum := byte(0)

	recenterTritTo0 := func(in int16) byte {
		if in == -1 {
			return 2
		}
		return byte(in)
	}

	// Encode 5 trits per byte, as long as we have >= 5 trits.
	for end := 5; end <= n; end += 5 {
		accum = recenterTritTo0(poly.P[end-1])
		accum = 3*accum + recenterTritTo0(poly.P[end-2])
		accum = 3*accum + recenterTritTo0(poly.P[end-3])
		accum = 3*accum + recenterTritTo0(poly.P[end-4])
		accum = 3*accum + recenterTritTo0(poly.P[end-5])
		if err := w.WriteByte(accum); err != nil {
			return err
		}
	}
	if end := n - (n % 5); end < n {
		n--
		accum = recenterTritTo0(poly.P[n])
		for end < n {
			n--
			accum = 3*accum + recenterTritTo0(poly.P[n])
		}
		if err := w.WriteByte(accum); err != nil {
			return err
		}
	}
	return nil
}
