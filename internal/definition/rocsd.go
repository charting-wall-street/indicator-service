package definition

import (
	"github.com/godoji/candlestick"
	"math"
)

func ROCStandardDeviation() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"rocsd": {
				Meta: SeriesMeta{
					Kind: candlestick.LineChart,
					Axis: candlestick.CustomAxis,
				},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {

					// make sure size is set
					size := params[0]
					if len(params) == 0 {
						return 0, false
					}

					// fetch mean
					mean := p.Indicator("rocma", params).Last("rocma").Value

					// get roc, make sure to pass lagging candles parameter correctly
					n := 1
					if len(params) > 1 {
						n = params[1]
					}
					roc := p.Indicator("roc", []int{n})

					// calculate standard deviation
					dev := 0.0
					points := 0.0
					for i := 0; i < size; i++ {
						c := roc.FromLast("roc", i)
						if c.Missing {
							continue
						}
						dev += math.Pow(c.Value-mean, 2)
						points++
					}

					// discard results when there is less than two points
					if points <= 1 {
						return 0, false
					}

					// divide by n-1 as we are using a sample of a population
					std := math.Sqrt(dev / (points - 1))

					return std, true
				},
			},
		},
		Presets: [][]float64{
			{20},
			{30},
			{50},
			{100},
			{200},
			{500},
			{1000},
			{1500},
		},
	}
}
