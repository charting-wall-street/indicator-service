package compute

import (
	"github.com/godoji/candlestick"
	"inca/internal/definition"
	"inca/internal/indicator"
	"log"
	"strconv"
)

type Moment struct {
	calc *CalculableDefinition
	ts   int64
	tran bool
}

func (c *Moment) Time() int64 {
	return c.ts
}

func (c *Moment) IsMissing() bool {
	return c.AtTime(c.ts).Missing
}

func (c *Moment) IsTransition() bool {
	return c.tran
}

func (c *Moment) Last() *candlestick.Candle {
	return c.FromLast(0)
}

func (c *Moment) FromLast(i int) *candlestick.Candle {

	// offset check
	if i < 0 {
		log.Fatalln("candles cannot be index with a negative index")
	}
	if i > 5000 {
		log.Fatalln("history limit reached when accessing candle moment set")
	}

	// get requested timestamp
	ts := c.ts - int64(i)*c.calc.interval

	return c.AtTime(ts)
}

func (c *Moment) AtTime(ts int64) *candlestick.Candle {

	// get dataset containing our data point
	var data *candlestick.CandleSet
	if c.calc.tCandles != nil && c.ts == ts && c.IsTransition() {
		data = c.calc.tCandles
	} else if c.calc.data.currSet == nil || ts < c.calc.data.currSet.UnixFirst() {
		data = c.calc.data.prevSet
	} else {
		data = c.calc.data.currSet
	}

	// return missing candlestick if no data is available
	if data == nil {
		return &candlestick.Candle{
			Missing: true,
			Time:    ts,
		}
	}

	candle := data.AtTime(ts)
	if candle.Time != ts/data.Interval()*data.Interval() {
		log.Fatalln("indexed candle has the wrong timestamp")
	}

	return candle
}

type MomentIndicator struct {
	parent *Moment
	head   *candlestick.Indicator
	tail   *candlestick.Indicator
	tran   *candlestick.Indicator
}

func (m *MomentIndicator) AtTime(key string, ts int64) *candlestick.IndicatorValue {
	var data *candlestick.Indicator
	if m.tran != nil && m.parent.ts == ts && m.parent.IsTransition() {
		data = m.tran
	} else if ts >= m.head.UnixFirst() {
		data = m.head
	} else {
		data = m.tail
	}
	index := (ts - data.UnixFirst()) / data.Interval()
	if index >= 5000 || index < 0 {
		log.Fatalf("invalid index %d\n", index)
	}
	return data.AtTime(key, ts)
}

func (m *MomentIndicator) Last(key string) *candlestick.IndicatorValue {
	return m.FromLast(key, 0)
}

func (m *MomentIndicator) FromLast(key string, i int) *candlestick.IndicatorValue {

	// offset check
	if i < 0 {
		log.Fatalln("indicator cannot be index with a negative index")
	}
	if i > 5000 {
		log.Fatalln("history limit reached when accessing indicator moment set")
	}

	// get requested timestamp
	ts := m.parent.ts - int64(i)*m.parent.calc.interval

	return m.AtTime(key, ts)
}

func (c *Moment) Indicator(name string, params []int) definition.MomentIndicator {

	// create unique sub cache id
	id := name
	for _, param := range params {
		id += strconv.Itoa(param)
	}

	// try to use sub cache
	if c.calc.subCache[id] != nil {
		cached := c.calc.subCache[id]
		cached.parent = c
		return cached
	}

	// get definition
	def, ok := indicator.ByName(name)
	if !ok {
		log.Fatalf("sub definition %s was not found\n", name)
	}

	// compute indicator
	var err error
	mi := new(MomentIndicator)
	mi.parent = c

	// compute previous block first
	calc := NewCalculable(def, c.calc.symbol, c.calc.block-1, c.calc.interval, c.calc.resolution)
	calc.UseCache = c.calc.UseCache
	mi.tail, err = calc.Compute(params)

	// compute current block
	calc = NewCalculable(def, c.calc.symbol, c.calc.block, c.calc.interval, c.calc.resolution)
	calc.UseCache = c.calc.UseCache
	mi.head, err = calc.Compute(params)

	// in case we are using transition data we need to compute the transition interval
	if c.calc.tCandles != nil {
		calc.tCandles = c.calc.tCandles
		mi.tran, err = calc.Transition(c.calc.tCandles.BlockNumber(), c.calc.resolution, params)
		if err != nil {
			log.Fatalf("failed to compute sub indicator transition values %s\n", name)
		}
	}

	// store in sub cache for fast future access
	c.calc.subCache[id] = mi

	return mi
}
