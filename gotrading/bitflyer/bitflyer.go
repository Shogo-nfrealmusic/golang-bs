package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	baseURL             = "https://api.bitflyer.com/v1"
	wsEndpoint          = "wss://ws.lightstream.bitflyer.com/json-rpc"
	maxReconnectBackoff = 30 * time.Second
)

type APIClient struct {
	key string
	secret string
	httpClient *http.Client
}

func New(key, secret string) *APIClient {
	apiClient := &APIClient{key, secret, &http.Client{}}
	return apiClient
}

func (api APIClient) header(method, endpoint string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + method + endpoint + string(body)
	mac := hmac.New(sha256.New, []byte(api.secret))
	mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))
	return map[string]string{
		"ACCESS-KEY": api.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN": sign,
		"Content-Type": "application/json",
	}
}

func (api *APIClient) doRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	endpoint := strings.TrimSuffix(baseURL, "/") + "/" + strings.TrimPrefix(urlPath, "/")
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	for key, value := range api.header(method, req.URL.RequestURI(), data) {
		req.Header.Add(key, value)
	}
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type Balance struct {
	CurrencyCode string `json:"currency_code"`
	Amount float64 `json:"amount"`
	Available float64 `json:"available"`
}

func (api *APIClient) GetBalance() ([]Balance, error) {
	url := "me/getbalance"
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	if err != nil {
		log.Printf("action=GetBalance err=%s", err.Error())
		return nil, err
}
	var balances []Balance
	err = json.Unmarshal(resp, &balances)
	if err != nil {
		log.Printf("action=GetBalance err=%s", err.Error())
		return nil, err
	}
	return balances, nil
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64     `json:"best_bid"`
	BestAsk         float64     `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64     `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64     `json:"total_ask_depth"`
	MarketBidSize   float64     `json:"market_bid_size"`
	MarketAskSize   float64     `json:"market_ask_size"`
	Ltp             float64     `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

func (t *Ticker) GetMidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

func (t *Ticker) DateTime() time.Time {
	layouts := []string{
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05.999999999Z07:00",
		time.RFC3339Nano,
		time.RFC3339,
	}
	for _, layout := range layouts {
		datetime, err := time.ParseInLocation(layout, t.Timestamp, time.UTC)
		if err == nil {
			return datetime
		}
	}
	log.Printf("action=DateTime err=unsupported format timestamp=%s", t.Timestamp)
	return time.Time{}
}

func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().UTC().Truncate(duration)
}

func (api *APIClient) GetTicker(productCode string) (*Ticker, error) {
	url := "ticker"
	resp, err := api.doRequest("GET", url, map[string]string{"product_code": productCode}, nil)
	if err != nil {
		return nil, err
}
	var ticker Ticker
	err = json.Unmarshal(resp, &ticker)
	if err != nil {
		return nil, err
	}
	return &ticker, nil
}

type rpcSubscribeRequest struct {
	JSONRPC string             `json:"jsonrpc"`
	Method  string             `json:"method"`
	Params  rpcSubscribeParams `json:"params"`
}

type rpcSubscribeParams struct {
	Channel string `json:"channel"`
}

type rpcChannelMessage struct {
	JSONRPC string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  rpcMessageParams `json:"params"`
}

type rpcMessageParams struct {
	Channel string          `json:"channel"`
	Message json.RawMessage `json:"message"`
}

// GetRealTimeTicker maintains a WebSocket subscription and pushes ticker updates to ch.
// It reconnects automatically when the connection closes or errors occur.
func (api *APIClient) GetRealTimeTicker(symbol string, ch chan<- Ticker) {
	channel := fmt.Sprintf("lightning_ticker_%s", symbol)
	backoff := time.Second

	for {
		err := api.streamRealTimeTicker(channel, ch)
		if err != nil {
			log.Printf("action=GetRealTimeTicker channel=%s err=%s", channel, err)
		}

		log.Printf("action=GetRealTimeTicker channel=%s reconnecting in %s", channel, backoff)
		time.Sleep(backoff)
		if backoff < maxReconnectBackoff {
			backoff *= 2
			if backoff > maxReconnectBackoff {
				backoff = maxReconnectBackoff
			}
		}
	}
}

func (api *APIClient) streamRealTimeTicker(channel string, ch chan<- Ticker) error {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		return fmt.Errorf("websocket dial: %w", err)
	}
	defer conn.Close()

	subscribe := rpcSubscribeRequest{
		JSONRPC: "2.0",
		Method:  "subscribe",
		Params: rpcSubscribeParams{
			Channel: channel,
		},
	}
	if err := conn.WriteJSON(subscribe); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	log.Printf("action=GetRealTimeTicker channel=%s subscribed", channel)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read message: %w", err)
		}

		var envelope rpcChannelMessage
		if err := json.Unmarshal(message, &envelope); err != nil {
			log.Printf("action=GetRealTimeTicker channel=%s parse envelope err=%s body=%s", channel, err, message)
			continue
		}

		if envelope.Method != "channelMessage" {
			continue
		}

		var ticker Ticker
		if err := json.Unmarshal(envelope.Params.Message, &ticker); err != nil {
			log.Printf("action=GetRealTimeTicker channel=%s parse ticker err=%s body=%s", channel, err, envelope.Params.Message)
			continue
		}

		ch <- ticker
	}
}