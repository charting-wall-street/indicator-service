package definition

import (
	"github.com/godoji/candlestick"
	"math"
)

func StandardDeviation() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"sd": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.CustomAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					size := params[0]
					sma := p.Indicator("sma", params).Last("sma").Value
					dev := 0.0
					points := 0.0
					for i := 0; i < size; i++ {
						c := p.FromLast(i)
						if c.Missing {
							continue
						}
						dev += math.Pow(c.Close-sma, 2)
						points++
					}
					if points <= 1 {
						return 0, false
					}
					std := math.Sqrt(dev / (points - 1))
					return std, true
				},
			},
		},
		Presets: [][]float64{
			{20},
			{30},
			{50},
			{100},
			{200},
			{500},
			{1000},
			{1500},
		},
	}
}
