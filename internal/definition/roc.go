package definition

import (
	"github.com/godoji/candlestick"
)

func RateOfChange() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"roc": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					n := 1
					if len(params) != 0 {
						n = params[0]
					}
					if n <= 0 {
						return 0, false
					}
					prev := p.FromLast(n)
					curr := p.Last()
					if prev.Missing || curr.Missing {
						return 0, false
					}
					return ((curr.Close / prev.Close) - 1.0) * 100, true
				},
			},
		},
		Presets: [][]float64{
			{0},
			{1},
			{2},
			{3},
			{4},
			{5},
			{6},
			{7},
			{14},
			{21},
			{27},
			{30},
		},
	}
}
