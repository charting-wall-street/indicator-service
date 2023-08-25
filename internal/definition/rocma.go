package definition

import (
	"github.com/godoji/candlestick"
)

func ROCMovingAverage() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"rocma": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.CustomAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					size := params[0]
					acc := 0.0
					n := 1
					if len(params) > 1 {
						n = params[1]
					}
					roc := p.Indicator("roc", []int{n})
					points := 0
					for i := 0; i < size; i++ {
						c := roc.FromLast("roc", i)
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
