package definition

import (
	"github.com/godoji/candlestick"
)

func MovingAverageConvergenceDivergence() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"macd": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					ema12 := p.Indicator("ema", []int{12}).Last("ema")
					ema26 := p.Indicator("ema", []int{26}).Last("ema")
					if ema12.Missing || ema26.Missing {
						return 0, false
					}
					return ema12.Value - ema26.Value, true
				},
			},
			"signal": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.CustomAxis,
				},
				Init: func(p Moment, params []int) float64 {
					ema12 := p.Indicator("ema", []int{12})
					ema26 := p.Indicator("ema", []int{26})
					length := 9.0
					k := 2.0 / (length + 1.0)
					signal := 0.0
					for i := warmUpLength + 1; i >= 0; i-- {
						macd := ema12.FromLast("ema", i).Value - ema26.FromLast("ema", i).Value
						signal = macd*k + signal*(1.0-k)
					}
					return signal
				},
				Step: func(p Moment, prev float64, params []int) (float64, bool) {
					ema12 := p.Indicator("ema", []int{12}).Last("ema")
					ema26 := p.Indicator("ema", []int{26}).Last("ema")
					if ema12.Missing || ema26.Missing {
						return prev, false
					}
					macd := ema12.Value - ema26.Value
					length := 9.0
					k := 2.0 / (length + 1.0)
					return macd*k + prev*(1.0-k), true
				},
			},
		},
		Presets: [][]float64{
			{0},
		},
	}
}
