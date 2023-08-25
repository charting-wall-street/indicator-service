package definition

import (
	"github.com/godoji/candlestick"
)

func Trailing() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"open": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					if len(params) == 0 || params[0] < 0 {
						return 0, true
					}
					c := p.FromLast(params[0])
					return c.Open, !c.Missing
				},
			},
			"close": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					if len(params) == 0 || params[0] < 0 {
						return 0, true
					}
					c := p.FromLast(params[0])
					return c.Close, !c.Missing
				},
			},
			"high": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					if len(params) == 0 || params[0] < 0 {
						return 0, true
					}
					c := p.FromLast(params[0])
					return c.High, !c.Missing
				},
			},
			"low": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					if len(params) == 0 || params[0] < 0 {
						return 0, true
					}
					c := p.FromLast(params[0])
					return c.Low, !c.Missing
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
