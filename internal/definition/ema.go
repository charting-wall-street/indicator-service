package definition

import (
	"github.com/godoji/candlestick"
)

const warmUpLength = 1500

func ExponentialMovingAverage() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"ema": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Init: func(p Moment, params []int) float64 {
					length := params[0]
					alpha := 2.0 / (float64(length) + 1.0)
					ema := p.FromLast(warmUpLength).Close
					for i := warmUpLength + 1; i >= 0; i-- {
						if p.FromLast(i).Missing {
							continue
						}
						ema = p.FromLast(i).Close*alpha + ema*(1.0-alpha)
					}
					return ema
				},
				Step: func(p Moment, prev float64, params []int) (float64, bool) {
					length := params[0]
					c := p.Last()
					k := 2.0 / (float64(length) + 1.0)
					return c.Close*k + prev*(1.0-k), true
				},
			},
		},
		Presets: [][]float64{
			{5},
			{7},
			{10},
			{12},
			{20},
			{21},
			{26},
			{50},
			{100},
			{200},
		},
	}
}
