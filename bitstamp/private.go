package bitstamp

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ivivanov/crypto-bot/bitstamp/response"
)

func (c *Conn) PostBalanceAll() (*response.Balance, error) {
	return c.PostBalance("")
}

func (c *Conn) PostBalance(currencyPair string) (*response.Balance, error) {
	v := url.Values{}
	path := fmt.Sprint("/v2/balance/")
	if currencyPair != "" {
		path += fmt.Sprintf("%s/", currencyPair)
	}
	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}
	res := &response.Balance{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) PostAccountBalances(currencyPair string) (*[]response.AccountBalances, error) {
	v := url.Values{}
	path := fmt.Sprint("/v2/account_balances/")

	if currencyPair != "" {
		path += fmt.Sprintf("%s/", currencyPair)
	}

	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}

	res := &[]response.AccountBalances{}

	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Conn) PostUserTransactionsAll(offset int, limit int, sort string) (*[]response.UserTransaction, error) {
	return c.PostUserTransactions("", offset, limit, sort)
}

func (c *Conn) PostUserTransactions(currencyPair string, offset int, limit int, sort string) (*[]response.UserTransaction, error) {
	v := url.Values{}
	if offset > 0 {
		v.Set("offset", fmt.Sprint(offset))
	}
	if limit > 0 {
		v.Set("limit", fmt.Sprint(limit))
	}
	if sort != "" {
		v.Set("sort", sort)
	}
	path := fmt.Sprint("/v2/user_transactions/")
	if currencyPair != "" {
		path += fmt.Sprintf("%s/", currencyPair)
	}
	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}
	res := &[]response.UserTransaction{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) PostOpenOrdersAll() (*[]response.OpenOrder, error) {
	return c.PostOpenOrders("all")
}

func (c *Conn) PostOpenOrders(currencyPair string) (*[]response.OpenOrder, error) {
	v := url.Values{}
	path := fmt.Sprintf("/v2/open_orders/%s/", currencyPair)
	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}
	res := &[]response.OpenOrder{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) PostOrderStatus(id int64) (*response.OrderStatus, error) {
	v := url.Values{}
	v.Set("id", fmt.Sprint(id))
	path := fmt.Sprint("/order_status/")
	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}
	res := &response.OrderStatus{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) PostCancelOrder(id int64) (*response.CancelOrder, error) {
	v := url.Values{}
	v.Set("id", fmt.Sprint(id))
	path := fmt.Sprint("/v2/cancel_order/")
	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}
	res := &response.CancelOrder{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) PostCancelAllOrders(currencyPair string) (*response.CancelAllOrders, error) {
	v := url.Values{}
	path := fmt.Sprintf("/v2/cancel_all_orders/%s/", currencyPair)

	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}

	res := &response.CancelAllOrders{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Conn) PostBuyLimitOrder(currencyPair, clientOrderID string, amount, price float64) (*response.LimitOrder, error) {
	v := url.Values{}
	v.Set("amount", fmt.Sprint(amount))
	v.Set("price", fmt.Sprintf("%.4f", price))

	if clientOrderID != "" {
		v.Set("client_order_id", fmt.Sprint(clientOrderID))
	}

	path := fmt.Sprintf("/v2/buy/%s/", currencyPair)
	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}
	res := &response.LimitOrder{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) PostBuyMarketOrder(currencyPair string, amount float64) (*response.MarketOrder, error) {
	v := url.Values{}
	v.Set("amount", fmt.Sprint(amount))
	path := fmt.Sprintf("/v2/buy/market/%s/", currencyPair)
	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}
	res := &response.MarketOrder{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) PostSellLimitOrder(currencyPair, clientOrderID string, amount, price float64) (*response.LimitOrder, error) {
	v := url.Values{}
	v.Set("amount", fmt.Sprint(amount))
	v.Set("price", fmt.Sprintf("%.4f", price))

	if clientOrderID != "" {
		v.Set("client_order_id", fmt.Sprint(clientOrderID))
	}

	path := fmt.Sprintf("/v2/sell/%s/", currencyPair)
	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}

	res := &response.LimitOrder{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Conn) PostSellMarketOrder(currencyPair string, amount float64) (*response.MarketOrder, error) {
	v := url.Values{}
	v.Set("amount", fmt.Sprint(amount))
	path := fmt.Sprintf("/v2/sell/market/%s/", currencyPair)
	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}
	res := &response.MarketOrder{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) PostWithdrawalRequests(timedelta int) (*[]response.WithdrawalRequest, error) {
	v := url.Values{}
	if timedelta > 0 {
		v.Set("timedelta", fmt.Sprint(timedelta))
	}
	path := fmt.Sprintf("/v2/withdrawal-requests/")
	b, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}
	res := &[]response.WithdrawalRequest{}
	err = json.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) WebsocketToken() (*response.WebsocketToken, error) {
	path := fmt.Sprint("/v2/websockets_token/")
	v := url.Values{}

	r, err := c.Request("POST", path, v, true)
	if err != nil {
		return nil, err
	}

	res := &response.WebsocketToken{}
	err = json.Unmarshal(r, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
