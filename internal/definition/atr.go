package definition

import (
	"github.com/godoji/candlestick"
	"math"
)

func AverageTrueRange() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"atr": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.CustomAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					size := params[0]
					acc := 0.0
					points := 0
					for i := 0; i < size; i++ {
						c := p.FromLast(i)
						if c.Missing {
							continue
						}
						acc += math.Abs(c.High - c.Low)
						points++
					}
					if points == 0 {
						return 0, false
					}
					return acc / float64(points), true
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
