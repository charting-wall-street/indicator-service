package definition

import (
	"github.com/godoji/candlestick"
)

func OBVMovingAverage() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"obvma": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.PriceAxis,
				},
				Init: func(p Moment, params []int) float64 {
					obv := p.Indicator("obv", []int{params[0]})
					length := params[1]
					alpha := 2.0 / (float64(length) + 1.0)
					ema := p.FromLast(warmUpLength).Close
					for i := warmUpLength + 1; i >= 0; i-- {
						v := obv.FromLast("obv", i)
						if v.Missing {
							continue
						}
						ema = v.Value*alpha + ema*(1.0-alpha)
					}
					return ema
				},
				Step: func(p Moment, prev float64, params []int) (float64, bool) {
					obv := p.Indicator("obv", []int{params[0]})
					length := params[1]
					k := 2.0 / (float64(length) + 1.0)
					v := obv.Last("obv")
					return v.Value*k + prev*(1.0-k), !v.Missing
				},
			},
		},
		Presets: [][]float64{
			{200, 10},
			{200, 50},
		},
	}
}
