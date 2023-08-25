package definition

import "github.com/godoji/candlestick"

func calculateFibDown(p Moment, params []int, level float64) (float64, bool) {
	topLength := 150
	if len(params) > 0 {
		topLength = params[0]
	}
	bottomLength := 150
	if len(params) > 1 {
		bottomLength = params[1]
	}
	lagLength := 10
	if len(params) > 2 {
		lagLength = params[2]
	}
	ath := p.Indicator("ath", []int{topLength, lagLength}).Last("ath")
	atl := p.Indicator("atl", []int{bottomLength, lagLength}).Last("atl")
	if ath.Missing || atl.Missing {
		return 0, false
	}
	delta := ath.Value - atl.Value
	return atl.Value + delta*level, true
}

func FibonacciDowntrend() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"0": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					return calculateFibDown(p, params, 0)
				},
			},
			"0.2": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					return calculateFibDown(p, params, 0.236)
				},
			},
			"0.3": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					return calculateFibDown(p, params, 0.382)
				},
			},
			"0.5": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					return calculateFibDown(p, params, 0.5)
				},
			},
			"0.6": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					return calculateFibDown(p, params, 0.618)
				},
			},
			"0.7": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					return calculateFibDown(p, params, 0.786)
				},
			},
			"1": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					return calculateFibDown(p, params, 1)
				},
			},
			"1.6": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					return calculateFibDown(p, params, 1.618)
				},
			},
			"2.6": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					return calculateFibDown(p, params, 2.618)
				},
			},
		},
		Presets: [][]float64{
			{100, 100 / 5},
			{150, 150 / 5},
			{200, 200 / 5},
			{400, 400 / 5},
			{500, 500 / 5},
			{600, 600 / 5},
			{750, 750 / 5},
			{900, 900 / 5},
			{1000, 1000 / 5},
		},
	}
}
