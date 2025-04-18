package bitstamp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ivivanov/crypto-bot/bitstamp/response"
)

const (
	addr = "https://www.bitstamp.net/api"
)

type Conn struct {
	apiKey     string
	apiSecret  string
	customerID string
	verbose    bool
}

// NewConn returns a new Conn.
func NewConn() (*Conn, error) {
	return &Conn{}, nil
}

// NewConn returns a new Conn with authentication parameters.
// Customer ID: https://www.bitstamp.net/account/balance/
// API key & secret: https://www.bitstamp.net/account/security/api/
func NewAuthConn(apiKey, apiSecret, customerID string, verbose bool) (*Conn, error) {
	conn := &Conn{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		customerID: customerID,
		verbose:    verbose,
	}
	return conn, nil
}

// Close closes the connection.
func (c *Conn) Close() error {
	return nil
}

// SetAuth sets authentication parameters.
func (c *Conn) SetAuth(apiKey, apiSecret, customerID string) {
	c.apiKey = apiKey
	c.apiSecret = apiSecret
	c.customerID = customerID
}

func (c *Conn) Request(method string, path string, values url.Values, auth bool) ([]byte, error) {
	if auth {
		if c.customerID == "" || c.apiKey == "" || c.apiSecret == "" {
			return nil, ErrAuthRequired
		}
		setAuthValues(values, c.customerID, c.apiKey, c.apiSecret)
	}

	body := strings.NewReader(values.Encode())

	client := &http.Client{}
	req, err := http.NewRequest(method, addr+path, body)
	if err != nil {
		return nil, err
	}

	if method == http.MethodGet {
		query := req.URL.Query()

		for k, v := range values {
			query.Add(k, v[0])
		}

		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.verbose {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(b))
		fmt.Println(req.URL.String())
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if c.verbose {
		fmt.Println(string(b))
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		res := &response.Error{}
		err = json.Unmarshal(b, res)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("HTTP Status: %s, Error: %s, Reason: %s, Status: %s", resp.Status, res.Code, res.Reason, res.Status)
	}
	return b, nil
}

func setAuthValues(values url.Values, customerID, apiKey, apiSecret string) {
	nonce := genNonce()
	signature := genSignature(nonce, customerID, apiKey, apiSecret)
	values.Set("key", apiKey)
	values.Set("nonce", nonce)
	values.Set("signature", signature)
}

func genNonce() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func genSignature(nonce, customerID, apiKey, apiSecret string) string {
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(nonce + customerID + apiKey))
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}
