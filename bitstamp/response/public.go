package response

type Ticker struct {
	Open      float64 `json:"open,string"`
	High      float64 `json:"high,string"`
	Last      float64 `json:"last,string"`
	Timestamp int64   `json:"timestamp,string"`
	Bid       float64 `json:"bid,string"`
	Vwap      float64 `json:"vwap,string"`
	Volume    float64 `json:"volume,string"`
	Low       float64 `json:"low,string"`
	Ask       float64 `json:"ask,string"`
}

type OrderBook struct {
	Timestamp string  `json:"timestamp"`
	Bids      []Order `json:"bids"`
	Asks      []Order `json:"asks"`
}

type Transaction struct {
	Amount float64 `json:"amount,string"`
	Type   int     `json:"type,string"`
	Price  float64 `json:"price,string"`
	Tid    int64   `json:"tid,string"`
	Date   Date    `json:"date"`
}

type TradingPairInfo struct {
	Description     string `json:"description"`
	URLSymbol       string `json:"url_symbol"`
	Trading         string `json:"trading"`
	CounterDecimals int    `json:"counter_decimals"`
	Name            string `json:"name"`
	MinimumOrder    string `json:"minimum_order"`
	BaseDecimals    int    `json:"base_decimals"`
}

type OHLCData struct {
	OHCL []OHLC `json:"ohlc"`
	Pair string `json:"pair"`
}

type OHLC struct {
	Close     float64 `json:"close,string"`
	High      float64 `json:"high,string"`
	Low       float64 `json:"low,string"`
	Open      float64 `json:"open,string"`
	Volume    float64 `json:"volume,string"`
	Timestamp string  `json:"timestamp"`
}
