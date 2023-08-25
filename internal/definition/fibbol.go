package definition

import (
	"github.com/godoji/candlestick"
)

const FIB1 = 1.618
const FIB2 = 2.618
const FIB3 = 4.236

func FibonacciBollinger() *Definition {
	return &Definition{
		Series: map[string]*SeriesDefinition{
			"top1": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					atr := p.Indicator("atr", []int{params[0]}).Last("atr")
					sma := p.Indicator("sma", []int{params[0]}).Last("sma")
					if sma.Missing || atr.Missing {
						return 0, false
					}
					return sma.Value + atr.Value*FIB1, true
				},
			},
			"top2": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					atr := p.Indicator("atr", []int{params[0]}).Last("atr")
					sma := p.Indicator("sma", []int{params[0]}).Last("sma")
					if sma.Missing || atr.Missing {
						return 0, false
					}
					return sma.Value + atr.Value*FIB2, true
				},
			},
			"top3": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					atr := p.Indicator("atr", []int{params[0]}).Last("atr")
					sma := p.Indicator("sma", []int{params[0]}).Last("sma")
					if sma.Missing || atr.Missing {
						return 0, false
					}
					return sma.Value + atr.Value*FIB3, true
				},
			},
			"bottom1": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					atr := p.Indicator("atr", []int{params[0]}).Last("atr")
					sma := p.Indicator("sma", []int{params[0]}).Last("sma")
					if sma.Missing || atr.Missing {
						return 0, false
					}
					return sma.Value - atr.Value*FIB1, true
				},
			},
			"bottom2": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					atr := p.Indicator("atr", []int{params[0]}).Last("atr")
					sma := p.Indicator("sma", []int{params[0]}).Last("sma")
					if sma.Missing || atr.Missing {
						return 0, false
					}
					return sma.Value - atr.Value*FIB2, true
				},
			},
			"bottom3": {
				Meta: SeriesMeta{Kind: candlestick.LineChart, Axis: candlestick.PriceAxis},
				Step: func(p Moment, _ float64, params []int) (float64, bool) {
					atr := p.Indicator("atr", []int{params[0]}).Last("atr")
					sma := p.Indicator("sma", []int{params[0]}).Last("sma")
					if sma.Missing || atr.Missing {
						return 0, false
					}
					return sma.Value - atr.Value*FIB3, true
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
			{400},
		},
	}
}
