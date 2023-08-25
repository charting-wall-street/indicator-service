package definition

import (
	"github.com/godoji/candlestick"
)

func RelativeStrengthIndex() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"rsi": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.CustomAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					size := params[0]
					avgU := 0.0
					avgD := 0.0
					points := 0.0
					for i := 0; i < size; i++ {
						c := p.FromLast(i)
						cPrev := p.FromLast(i + 1)
						if c.Missing {
							continue
						}
						if cPrev.Missing {
							continue
						}
						curr := c.Close
						prev := cPrev.Close
						if prev <= curr {
							avgU += curr - prev
						} else {
							avgD += prev - curr
						}
						points++
					}
					if points == 0 {
						return 0, false
					}
					if avgD == 0 {
						return 100, true
					}
					rs := (avgU / points) / (avgD / points)
					rsi := 100 - 100/(1+rs)
					return rsi, true
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
