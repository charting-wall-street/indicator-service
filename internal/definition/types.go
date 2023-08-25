package definition

import "github.com/godoji/candlestick"

type Moment interface {
	Time() int64
	IsMissing() bool
	IsTransition() bool
	Last() *candlestick.Candle
	FromLast(i int) *candlestick.Candle
	AtTime(ts int64) *candlestick.Candle
	Indicator(name string, params []int) MomentIndicator
}

type MomentIndicator interface {
	AtTime(key string, ts int64) *candlestick.IndicatorValue
	Last(key string) *candlestick.IndicatorValue
	FromLast(key string, i int) *candlestick.IndicatorValue
}

type ValueFunction = func(p Moment, last float64, params []int) (float64, bool)
type FirstFunction = func(p Moment, params []int) float64

type SeriesMeta struct {
	Axis candlestick.AxisType
	Kind candlestick.SeriesType
}

type SeriesDefinition struct {
	Meta SeriesMeta
	Step ValueFunction
	Init FirstFunction
}

type Definition struct {
	Name    string
	Series  map[string]*SeriesDefinition
	Presets [][]float64
}
