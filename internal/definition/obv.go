package definition

import (
	"github.com/godoji/candlestick"
)

func OnBalanceVolume() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"obv": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.CustomAxis,
				},
				Init: func(p Moment, params []int) float64 {
					alpha := 2.0 / (float64(params[0]) + 1.0)
					obv := float64(0)
					for i := warmUpLength + 1; i >= 0; i-- {
						prevPrice := p.FromLast(i + 1)
						currPrice := p.FromLast(i)
						if prevPrice.Missing || currPrice.Missing {
							continue
						}
						delta := float64(0)
						if currPrice.Close > prevPrice.Close {
							delta = p.Last().Volume
						} else if currPrice.Close < prevPrice.Close {
							delta = -p.Last().Volume
						}
						obv = (1.0-alpha)*obv + alpha*delta
					}
					return obv
				},
				Step: func(p Moment, prev float64, params []int) (float64, bool) {
					prevPrice := p.FromLast(1)
					currPrice := p.FromLast(0)
					delta := float64(0)
					if currPrice.Close > prevPrice.Close {
						delta = p.Last().Volume
					} else if currPrice.Close < prevPrice.Close {
						delta = -p.Last().Volume
					}
					k := 2.0 / (float64(params[0]) + 1.0)
					return prev*(1.0-k) + delta*k, !prevPrice.Missing && !currPrice.Missing
				},
			},
		},
		Presets: [][]float64{
			{100},
			{250},
			{500},
		},
	}
}
