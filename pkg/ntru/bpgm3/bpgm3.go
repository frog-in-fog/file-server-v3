// Package bpgm3 implements the BPGM3 algorithm defined in the X9.98 spec.
package bpgm3

import (
	"file-server-v3/pkg/ntru/igf2"
	"file-server-v3/pkg/ntru/polynomial"
)

// GenTrinomial generates a trinonial of degree N-1 that has numOnes
// coefficients set to +1 and numNegones coefficients set to -1, and all other
// coefficients set to 0.
func GenTrinomial(n, numOnes, numNegOnes int16, igf *igf2.IGF2) (*polynomial.Full, error) {
	isSet := make([]bool, n)
	p := polynomial.New(int(n))
	for t := int16(0); t < numOnes; {
		i, err := igf.NextIndex()
		if err != nil {
			return nil, err
		}
		if isSet[i] {
			continue
		}
		p.P[i] = 1
		isSet[i] = true
		t++
	}
	for t := int16(0); t < numNegOnes; {
		i, err := igf.NextIndex()
		if err != nil {
			return nil, err
		}
		if isSet[i] {
			continue
		}
		p.P[i] = -1
		isSet[i] = true
		t++
	}
	return p, nil
}
