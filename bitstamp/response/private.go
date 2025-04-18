package response

type Balance struct {
	BtcAvailable float64 `json:"btc_available,string"`
	BtcBalance   float64 `json:"btc_balance,string"`
	BtcReserved  float64 `json:"btc_reserved,string"`
	BtceurFee    float64 `json:"btceur_fee,string"`
	BtcusdFee    float64 `json:"btcusd_fee,string"`
	EthAvailable float64 `json:"eth_available,string"`
	EthBalance   float64 `json:"eth_balance,string"`
	EthReserved  float64 `json:"eth_reserved,string"`
	EthbtcFee    float64 `json:"ethbtc_fee,string"`
	EtheurFee    float64 `json:"etheur_fee,string"`
	EthusdFee    float64 `json:"ethusd_fee,string"`
	EurAvailable float64 `json:"eur_available,string"`
	EurBalance   float64 `json:"eur_balance,string"`
	EurReserved  float64 `json:"eur_reserved,string"`
	EurusdFee    float64 `json:"eurusd_fee,string"`
	LtcAvailable float64 `json:"ltc_available,string"`
	LtcBalance   float64 `json:"ltc_balance,string"`
	LtcReserved  float64 `json:"ltc_reserved,string"`
	LtcbtcFee    float64 `json:"ltcbtc_fee,string"`
	LtceurFee    float64 `json:"ltceur_fee,string"`
	LtcusdFee    float64 `json:"ltcusd_fee,string"`
	UsdAvailable float64 `json:"usd_available,string"`
	UsdBalance   float64 `json:"usd_balance,string"`
	UsdReserved  float64 `json:"usd_reserved,string"`
	XrpAvailable float64 `json:"xrp_available,string"`
	XrpBalance   float64 `json:"xrp_balance,string"`
	XrpReserved  float64 `json:"xrp_reserved,string"`
	XrpbtcFee    float64 `json:"xrpbtc_fee,string"`
	XrpeurFee    float64 `json:"xrpeur_fee,string"`
	XrpusdFee    float64 `json:"xrpusd_fee,string"`
}

type AccountBalances struct {
	Currency  string  `json:"currency"`
	Total     float64 `json:"total,string"`
	Available float64 `json:"available,string"`
	Reserved  float64 `json:"reserved,string"`
}

type UserTransaction struct {
	Fee      float64  `json:"fee,string"`
	OrderID  int64    `json:"order_id"`
	ID       int64    `json:"id"`
	Usd      float64  `json:"usd"`
	XrpEur   float64  `json:"xrp_eur"`
	Btc      float64  `json:"btc"`
	DateTime DateTime `json:"datetime"`
	Type     int      `json:"type,string"`
	Xrp      float64  `json:"xrp,string"`
	Eur      float64  `json:"eur,string"`
}

type OpenOrder struct {
	Price        float64  `json:"price,string"`
	CurrencyPair string   `json:"currency_pair"`
	DateTime     DateTime `json:"datetime"`
	Amount       float64  `json:"amount,string"`
	Type         int      `json:"type,string"`
	ID           int64    `json:"id,string"`
}

type OrderStatus struct {
	Status       string        `json:"status"`
	Transactions []Transaction `json:"transactions"`
}

type CancelOrder struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
	Type   int     `json:"type"`
	ID     int64   `json:"id"`
}

type CancelAllOrders struct {
	Success  bool `json:"success"`
	Canceled []struct {
		ID           int64   `json:"id"`
		CurrencyPair string  `json:"currency_pair"`
		Price        float64 `json:"price"`
		Amount       float64 `json:"amount"`
		Type         int     `json:"type"`
	}
}

type MarketOrder struct {
	Price    float64  `json:"price,string"`
	Amount   float64  `json:"amount,string"`
	Type     int      `json:"type,string"`
	ID       int64    `json:"id,string"`
	DateTime DateTime `json:"datetime"`
}

type LimitOrder struct {
	Price         float64             `json:"price,string"`
	Amount        float64             `json:"amount,string"`
	Type          int                 `json:"type,string"`
	ID            int64               `json:"id,string"`
	DateTime      DateTime            `json:"datetime"`
	ClientOrderID string              `json:"client_order_id,omitempty"`
	Status        string              `json:"status,omitempty"`
	Reason        map[string][]string `json:"reason,omitempty"`
}

func (blo *LimitOrder) IsError() bool {
	return blo.Status == "error"
}

type WithdrawalRequest struct {
}

type WebsocketToken struct {
	Token    string `json:"token"`
	ValidSec int    `json:"valid_sec"`
	UserID   int    `json:"user_id"`
}
