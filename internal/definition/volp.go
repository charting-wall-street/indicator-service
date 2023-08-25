package definition

import (
	"github.com/godoji/candlestick"
)

func VolumePrice() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"vol": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.CustomAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					return p.Last().Volume * p.Last().Close, !p.Last().Missing
				},
			},
		},
		Presets: [][]float64{
			{0},
		},
	}
}
