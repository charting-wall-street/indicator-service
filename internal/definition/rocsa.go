package definition

import (
	"github.com/godoji/candlestick"
)

func ROCSmoothedAverage() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"rocsa": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.CustomAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					if len(params) < 2 {
						return 0, false
					}
					acc := 0.0
					points := 0
					rocma := p.Indicator("rocma", []int{params[1]})
					for i := 0; i < params[0]; i++ {
						c := rocma.FromLast("rocma", i)
						if c.Missing {
							continue
						}
						acc += c.Value
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
			{1000},
		},
	}
}
