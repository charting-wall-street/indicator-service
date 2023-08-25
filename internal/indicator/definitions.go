package indicator

import (
	"inca/internal/definition"
	"sort"
	"strings"
)

var indicatorMap = make(map[string]*definition.Definition)

func init() {
	indicatorMap["sma"] = definition.SimpleMovingAverage()
	indicatorMap["ema"] = definition.ExponentialMovingAverage()
	indicatorMap["sd"] = definition.StandardDeviation()
	indicatorMap["bb"] = definition.BollingerBands()
	indicatorMap["close"] = definition.ClosePrice()
	indicatorMap["ath"] = definition.AllTimeHigh()
	indicatorMap["atl"] = definition.AllTimeLow()
	indicatorMap["rsi"] = definition.RelativeStrengthIndex()
	indicatorMap["atr"] = definition.AverageTrueRange()
	indicatorMap["fibbol"] = definition.FibonacciBollinger()
	indicatorMap["dema"] = definition.EMADerivative()
	indicatorMap["dsma"] = definition.SMADerivative()
	indicatorMap["fibull"] = definition.FibonacciUptrend()
	indicatorMap["fibear"] = definition.FibonacciDowntrend()
	indicatorMap["trail"] = definition.Trailing()
	indicatorMap["macd"] = definition.MovingAverageConvergenceDivergence()
	indicatorMap["roc"] = definition.RateOfChange()
	indicatorMap["rocma"] = definition.ROCMovingAverage()
	indicatorMap["rocsa"] = definition.ROCSmoothedAverage()
	indicatorMap["rocsd"] = definition.ROCStandardDeviation()
	indicatorMap["volume"] = definition.Volume()
	indicatorMap["obv"] = definition.OnBalanceVolume()
	indicatorMap["obvma"] = definition.OBVMovingAverage()
	indicatorMap["liquid"] = definition.VolumePrice()
}

func ByName(name string) (*definition.Definition, bool) {
	d, ok := indicatorMap[name]
	if !ok {
		return nil, ok
	}
	d.Name = name
	return d, true
}

type ListItem struct {
	Name    string      `json:"name"`
	Presets [][]float64 `json:"presets"`
}

func List() []ListItem {
	keys := make([]ListItem, len(indicatorMap))
	i := 0
	for s, ind := range indicatorMap {
		keys[i].Name = s
		keys[i].Presets = ind.Presets
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i].Name, keys[j].Name) < 0
	})
	return keys
}
