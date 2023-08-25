package compute

import (
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/godoji/candlestick"
	"inca/internal/definition"
	"log"
	"math"
	"strconv"
	"time"
)

var (
	indicatorCache       *ristretto.Cache
	transitionCache      *ristretto.Cache
	indicatorCacheConfig = ristretto.Config{
		NumCounters: 1 << 12, // 8k items
		MaxCost:     1 << 31, // 2GB
		BufferItems: 64,
	}
	transitionCacheConfig = ristretto.Config{
		NumCounters: 1 << 10, // 2k items
		MaxCost:     1 << 29, // 512MB
		BufferItems: 64,
	}
)

func init() {

	// create cache for indicators and transitions
	var err error
	if indicatorCache, err = ristretto.NewCache(&indicatorCacheConfig); err != nil {
		log.Fatal(err)
	}
	if transitionCache, err = ristretto.NewCache(&transitionCacheConfig); err != nil {
		log.Fatal(err)
	}
}

func NewCalculable(def *definition.Definition, symbol string, block int64, interval int64, resolution int64) *CalculableDefinition {
	return &CalculableDefinition{
		lastFetch:  time.Now().UTC().Unix(),
		def:        def,
		symbol:     symbol,
		data:       nil,
		interval:   interval,
		resolution: resolution,
		block:      block,
		subCache:   map[string]*MomentIndicator{},
	}
}

func (d *CalculableDefinition) Moment(timeStamp int64) *Moment {
	return &Moment{
		calc: d,
		ts:   timeStamp,
		tran: false,
	}
}

func (d *CalculableDefinition) MomentInTransition(timeStamp int64) *Moment {
	return &Moment{
		calc: d,
		ts:   timeStamp,
		tran: true,
	}
}

type CalculableDefinition struct {
	lastFetch  int64
	def        *definition.Definition
	symbol     string
	data       *CandleDataSet
	block      int64
	interval   int64
	resolution int64
	subCache   map[string]*MomentIndicator
	tCandles   *candlestick.CandleSet
	UseCache   bool
}

func concatParams(params []int) string {
	res := ""
	for i := range params {
		if i != 0 {
			res += ","
		}
		res += strconv.Itoa(params[i])
	}
	return res
}

func (d *CalculableDefinition) Key(params []int) string {
	key := d.symbol + ":" + strconv.FormatInt(d.interval, 10) + ":" + strconv.FormatInt(d.block, 10) + ":"
	key += d.def.Name + ":"
	key += concatParams(params)
	return key
}

func (d *CalculableDefinition) TKey(block int64, params []int) string {
	key := d.symbol + ":" + strconv.FormatInt(d.interval, 10) + ":" + strconv.FormatInt(d.block, 10) + ":"
	key += d.def.Name + ":"
	key += concatParams(params) + ":"
	key += strconv.FormatInt(block, 10)
	return key
}

func isSameCandle(t1 int64, t2 int64, interval int64) bool {
	b1 := t1 / interval
	b2 := t2 / interval
	return b1 == b2
}

func printStatus(hit bool, tran bool, symbol string, ind string, block int64, interval int64, params string, nc bool) {
	t := "\033[36mINDICA\u001B[0m"
	if tran {
		t = "\033[33mTRANSI\033[0m"
	}
	msg := fmt.Sprintf("%s | %s | %s | %d | %d | %s", t, symbol, ind, block, interval, params)
	if hit {
		fmt.Println("\033[32mH\033[0m | " + msg)
	} else if nc {
		fmt.Println("\033[31mF\033[0m | " + msg)
	} else {
		fmt.Println("\033[31mM\033[0m | " + msg)
	}
}

func (d *CalculableDefinition) Transition(block int64, resolution int64, params []int) (*candlestick.Indicator, error) {

	key := d.TKey(block, params)
	if d.UseCache {
		v, ok := transitionCache.Get(key)
		printStatus(ok, true, d.symbol, d.def.Name, block, d.interval, concatParams(params), false)
		if ok {
			return v.(*candlestick.Indicator), nil
		}
	} else {
		printStatus(false, true, d.symbol, d.def.Name, block, d.interval, concatParams(params), true)
	}

	var err error
	d.data, err = candleSuperSet(d.block, d.interval, d.symbol)
	if err != nil {
		return nil, err
	}

	// get minute transition candles
	tCandles, err := fetchCandleSet(block, resolution, resolution, d.symbol, true)
	if err != nil {
		return nil, err
	}
	d.tCandles = tCandles

	result := &candlestick.Indicator{
		Meta: candlestick.IndicatorMeta{
			UID:          key,
			Block:        block,
			Complete:     d.lastFetch >= ((d.block + 1) * 5000 * d.interval),
			LastUpdate:   d.lastFetch,
			Symbol:       d.symbol,
			Interval:     resolution,
			BaseInterval: d.interval,
			Name:         d.def.Name,
			Parameters:   params,
		},
		Series: make(map[string]*candlestick.IndicatorSeries),
	}

	startTime := 5000 * resolution * block
	prevCandle := (startTime - resolution) / d.interval * d.interval

	for s, def := range d.def.Series {

		values := make([]candlestick.IndicatorValue, 5000)

		lastValue := float64(0)
		lastRelativeValue := float64(0)

		if def.Init != nil {
			lastValue = def.Init(d.Moment(prevCandle), params)
		}

		ts := startTime
		m := d.MomentInTransition(ts)
		for i := 0; i < 5000; i++ {
			m.ts = ts
			if m.IsMissing() {
				values[i].Missing = true
			} else {
				v, ok := def.Step(m, lastValue, params)
				if ok {
					values[i].Value = v
					lastRelativeValue = v
				} else {
					values[i].Missing = true
				}
			}
			if !isSameCandle(ts, ts+d.resolution, d.interval) {
				lastValue = lastRelativeValue
			}
			ts += d.resolution // update last timestamp
		}

		result.Series[s] = &candlestick.IndicatorSeries{
			Values: values,
			Kind:   def.Meta.Kind,
			Axis:   def.Meta.Axis,
		}
	}

	if !d.UseCache {
		return result, nil
	}

	activeBlock := time.Now().UTC().Unix() / d.interval / 5000
	cost := int64(len(result.Series) * (1 << 18))
	if activeBlock != d.block {
		transitionCache.Set(key, result, cost)
	} else {
		transitionCache.SetWithTTL(key, result, cost, time.Minute)
	}

	return result, nil
}

func (d *CalculableDefinition) Compute(params []int) (*candlestick.Indicator, error) {

	key := d.Key(params)
	if d.UseCache {
		v, ok := indicatorCache.Get(key)
		printStatus(ok, false, d.symbol, d.def.Name, d.block, d.interval, concatParams(params), false)
		if ok {
			return v.(*candlestick.Indicator), nil
		}
	} else {
		printStatus(false, false, d.symbol, d.def.Name, d.block, d.interval, concatParams(params), true)
	}

	var err error
	d.data, err = candleSuperSet(d.block, d.interval, d.symbol)
	if err != nil {
		return nil, err
	}

	result := &candlestick.Indicator{
		Meta: candlestick.IndicatorMeta{
			UID:          key,
			Block:        d.block,
			Complete:     d.lastFetch >= ((d.block + 1) * 5000 * d.interval),
			LastUpdate:   d.lastFetch,
			Symbol:       d.symbol,
			Interval:     d.interval,
			BaseInterval: d.interval,
			Name:         d.def.Name,
			Parameters:   params,
		},
		Series: make(map[string]*candlestick.IndicatorSeries),
	}

	startTime := d.block * d.interval * 5000
	hasNaN := false

	for s, def := range d.def.Series {

		values := make([]candlestick.IndicatorValue, 5000)

		lastValue := float64(0)
		if def.Init != nil {
			lastValue = def.Init(d.Moment(startTime-d.interval), params)
		}
		for i := 0; i < len(values); i++ {
			m := d.Moment(startTime + d.interval*int64(i))
			if m.IsMissing() {
				values[i].Missing = true
			} else {
				v, ok := def.Step(m, lastValue, params)
				if !hasNaN && math.IsNaN(v) {
					hasNaN = true
					fmt.Printf("spotted NaN in %s\n", d.def.Name)
				}
				if ok {
					values[i].Value = v
					lastValue = v
				} else {
					values[i].Missing = true
				}
			}
		}

		result.Series[s] = &candlestick.IndicatorSeries{
			Values: values,
			Kind:   def.Meta.Kind,
			Axis:   def.Meta.Axis,
		}
	}

	if !d.UseCache {
		return result, nil
	}

	activeBlock := time.Now().UTC().Unix() / d.interval / 5000
	cost := int64(len(result.Series) * (1 << 18))
	if activeBlock != d.block {
		indicatorCache.Set(key, result, cost)
	} else {
		indicatorCache.SetWithTTL(key, result, cost, time.Minute)
	}

	return result, nil
}
