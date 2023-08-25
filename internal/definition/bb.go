package definition

import (
	"github.com/godoji/candlestick"
)

func BollingerBands() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"lower": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					sma := p.Indicator("sma", []int{params[0]}).Last("sma")
					std := p.Indicator("sd", []int{params[0]}).Last("sd")
					if sma.Missing || std.Missing {
						return 0, false
					}
					return sma.Value - float64(params[1])*std.Value, true
				},
			},
			"upper": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					sma := p.Indicator("sma", []int{params[0]}).Last("sma")
					std := p.Indicator("sd", []int{params[0]}).Last("sd")
					if sma.Missing || std.Missing {
						return 0, false
					}
					return sma.Value + float64(params[1])*std.Value, true
				},
			},
		},
		Presets: [][]float64{
			{20, 2},
			{21, 2},
			{20, 3},
			{21, 3},
			{20, 4},
			{21, 4},
			{20, 5},
			{21, 5},
		},
	}
}
