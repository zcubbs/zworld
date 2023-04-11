package pgen

import (
	"github.com/zcubbs/opensimplex-go"
	"math"
)

type Octave struct {
	Freq, Scale float64
}

type NoiseMap struct {
	seed     int64
	noise    opensimplex.Noise
	exponent float64
	octaves  []Octave
}

func NewNoiseMap(seed int64, octaves []Octave, exponent float64) *NoiseMap {
	return &NoiseMap{
		seed:     seed,
		noise:    opensimplex.NewNormalized(seed),
		exponent: exponent,
		octaves:  octaves,
	}
}

func (n *NoiseMap) Get(x, y int) float64 {
	ret := 0.0
	for i := range n.octaves {
		xNoise := n.octaves[i].Freq * float64(x)
		yNoise := n.octaves[i].Freq * float64(y)
		ret += n.octaves[i].Scale * n.noise.Eval2(xNoise, yNoise)
	}

	ret = math.Pow(ret, n.exponent)
	return ret
}
