package compute

import (
	"errors"
	"fmt"
	"github.com/godoji/candlestick"
	"inca/internal/config"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var queueLock = sync.Mutex{}
var queueCounter = 0

func enterLimit() {
	for {
		queueLock.Lock()
		if queueCounter <= 32 {
			queueCounter += 1
			queueLock.Unlock()
			break
		}
		queueLock.Unlock()
		time.Sleep(5 * time.Millisecond)
	}
}

func exitLimit() {
	queueLock.Lock()
	queueCounter -= 1
	queueLock.Unlock()
}

type CandleDataSet struct {
	interval int64
	prevSet  *candlestick.CandleSet
	currSet  *candlestick.CandleSet
}

func candleSuperSet(block int64, interval int64, symbol string) (*CandleDataSet, error) {

	prevSet, err := fetchCandleSet(block-1, interval, 0, symbol, false)
	if err != nil {
		return nil, err
	}
	currSet, err := fetchCandleSet(block, interval, 0, symbol, false)
	if err != nil {
		return nil, err
	}

	return &CandleDataSet{
		prevSet:  prevSet,
		currSet:  currSet,
		interval: interval,
	}, nil
}

func fetchCandleSet(block int64, interval int64, resolution int64, symbol string, transition bool) (*candlestick.CandleSet, error) {

	enterLimit()
	defer exitLimit()

	// create url
	baseUrl := config.ServiceConfig().CandleServiceURL() + "/market"
	resParam := ""
	if transition {
		resParam = "&resolution=" + strconv.FormatInt(resolution, 10)
		baseUrl += "/t"
	}
	url := fmt.Sprintf("%s/%s?segment=%d%s&interval=%d", baseUrl, symbol, block, resParam, interval)

	// setup request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {

		return nil, errors.New("failed to setup candles request: " + err.Error())
	}

	// set headers
	req.Header.Set("Accept", "application/octet-stream")

	// execute request and handle any connection or url based error
	resp, err := http.DefaultClient.Do(req)
	if err != nil {

		return nil, errors.New("failed to fetch candles: " + err.Error())
	}

	// drain and close body
	defer func() {
		if resp.Body != nil {
			if _, err = io.Copy(io.Discard, resp.Body); err != nil {
				log.Print(err)
			}
			if err = resp.Body.Close(); err != nil {
				log.Print(err)
			}
		}
	}()

	// case when no candle data exists
	if resp.StatusCode == http.StatusNotFound {

		return nil, nil
	}

	// check if response is useful
	if resp.StatusCode != http.StatusOK {

		return nil, errors.New("request for candles failed: " + err.Error())
	}

	// decode json data
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		return nil, err
	}
	result, err := candlestick.DecodeCandleSet(body)
	if err != nil {

		return nil, err
	}

	return result, nil
}
