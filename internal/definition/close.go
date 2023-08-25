package definition

import (
	"github.com/godoji/candlestick"
)

func ClosePrice() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"close": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					c := p.FromLast(params[0])
					return c.Close, !c.Missing
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
