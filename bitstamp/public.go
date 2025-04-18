package bitstamp

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ivivanov/crypto-bot/bitstamp/response"
)

func (c *Conn) GetTicker(currencyPair string) (*response.Ticker, error) {
	v := url.Values{}
	path := fmt.Sprintf("/v2/ticker/%s/", currencyPair)
	b, err := c.Request("GET", path, v, false)
	if err != nil {
		return nil, err
	}
	res := &response.Ticker{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) GetTickerHour(currencyPair string) (*response.Ticker, error) {
	v := url.Values{}
	path := fmt.Sprintf("/v2/ticker_hour/%s/", currencyPair)
	b, err := c.Request("GET", path, v, false)
	if err != nil {
		return nil, err
	}
	res := &response.Ticker{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) GetOrderBook(currencyPair string) (*response.OrderBook, error) {
	v := url.Values{}
	path := fmt.Sprintf("/v2/order_book/%s/", currencyPair)
	b, err := c.Request("GET", path, v, false)
	if err != nil {
		return nil, err
	}
	res := &response.OrderBook{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) GetTransactions(currencyPair string, time string) (*[]response.Transaction, error) {
	v := url.Values{}
	if time != "" {
		v.Set("time", time)
	}
	path := fmt.Sprintf("/v2/transactions/%s/", currencyPair)
	b, err := c.Request("GET", path, v, false)
	if err != nil {
		return nil, err
	}
	res := &[]response.Transaction{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) GetTradingPairsInfo() (*[]response.TradingPairInfo, error) {
	v := url.Values{}
	path := fmt.Sprint("/v2/trading-pairs-info/")
	b, err := c.Request("GET", path, v, false)
	if err != nil {
		return nil, err
	}
	res := &[]response.TradingPairInfo{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TODO add response type
func (c *Conn) GetOHLC(currencyPair string, step, limit int) (*[]response.OHLC, error) {
	if currencyPair == "" || step == 0 {
		return nil, fmt.Errorf("all args are required")
	}

	v := url.Values{}
	v.Set("step", fmt.Sprint(step))
	v.Set("limit", fmt.Sprint(limit))
	path := fmt.Sprintf("/v2/ohlc/%s/", currencyPair)

	b, err := c.Request("GET", path, v, false)
	if err != nil {
		return nil, err
	}

	res := map[string]response.OHLCData{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	ohlc := res["data"].OHCL

	return &ohlc, nil
}
