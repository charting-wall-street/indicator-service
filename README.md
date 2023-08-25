# INCA - Indicator Calculation Application

## Overview

Indicators serve as vital tools for discerning market trends. The Indicator Calculation Application (INCA) is a specialized service designed to generate indicators based on the processed candle data.

This service enables real-time computation of popular indicators such as the Exponential Moving Average (EMA) and the Relative Strength Index (RSI). The primary objective is to consolidate all indicator logic into a singular service, enabling instant output for a requested indicator.

Given that indicators can be requested with a range of parameters, INCA employs advanced caching techniques to streamline the retrieval and processing of candles.

## Indicator Types

Indicators can be categorized as either stateful or stateless, based on their calculations.

- Stateful Indicators: These indicators, such as the EMA, use the previous value in their calculations. EMA is calculated using the following formula:  
  `EMA = Closing price x multiplier + EMA (previous day) x (1-multiplier)`  
  This requires access to the complete history of the dataset for its calculation.
- Stateless Indicators: These indicators, such as the Simple Moving Average (SMA), do not require the previous value for their calculations. SMA is calculated using the following formula:  
  `SMA = (A1+A2+...+An)/n`

Some indicators like Bollinger Bands (BB) are based on other indicators, for instance, the Simple Moving Average (SMA). INCA is capable of calculating these indicators recursively and strives to prevent redundant calculations.

## Usage

You can run INCA using the command line interface.

```shell
$ go run ./cmd/inca --help
Usage of inca:
  -candles-url string
        path to the candle service (default "http://localhost:9702")
  -origins string
        cors origins (default "*")
  -port string
        port from which to run the service (default "9703")
```
