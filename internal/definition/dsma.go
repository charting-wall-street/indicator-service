package definition

import (
	"github.com/godoji/candlestick"
)

func SMADerivative() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"dsma": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Step: func(p Moment, prev float64, params []int) (float64, bool) {
					sma := p.Indicator("sma", []int{params[0]})
					v1 := sma.FromLast("sma", 0).Value
					v2 := sma.FromLast("sma", 1).Value
					v3 := sma.FromLast("sma", 2).Value
					v4 := sma.FromLast("sma", 3).Value
					w1 := v3 - v4
					w2 := v1 - v2
					return (w2 + w1) / 2, true
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
