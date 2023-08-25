package definition

import (
	"github.com/godoji/candlestick"
)

func AllTimeLow() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"atl": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					low := 0.0
					size := params[0]
					offset := 1
					found := false
					if len(params) > 1 {
						offset = params[1]
					}
					for i := offset; i < offset+size; i++ {
						c := p.FromLast(i)
						if c.Missing {
							continue
						}
						found = true
						if c.Low < low || low == 0 {
							low = c.Low
						}
					}
					if !found {
						return 0, false
					}
					return low, true
				},
			},
		},
		Presets: [][]float64{
			{5},
			{10},
			{20},
			{30},
			{40},
			{50},
			{75},
			{100},
			{150},
			{200},
			{400},
			{500},
			{600},
			{750},
			{800},
			{1000},
			{2000},
		},
	}
}
