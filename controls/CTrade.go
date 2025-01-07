package controls

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL   = "https://api.binance.com"
	apiKey    = "KpCB5G5WQbKPwjy07JDxQo9TEohnmgACS38hVGibyVFZt74AMmzsYz8pzLmmh9Om"
	apiSecret = "3AHZveeV5KCwzCsK4XDHWNuMd7zuM9PSblJVGmpJ5JTB6u2tstX9GDAKHhUJjw2U"
)

func sign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func createOrder(symbol, side, orderType, quantity, price string) (string, error) {
	endpoint := "/api/v3/order"
	timestamp := time.Now().UnixMilli()
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", orderType)
	params.Set("quantity", quantity)
	if price != "" {
		params.Set("price", price)
	}
	params.Set("timeInForce", "GTC") // Условие исполнения ордера (Good-Til-Cancelled)
	params.Set("timestamp", fmt.Sprintf("%d", timestamp))

	query := params.Encode()
	signature := sign(query, apiSecret)
	query += "&signature=" + signature

	req, err := http.NewRequest("POST", baseURL+endpoint+"?"+query, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-MBX-APIKEY", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: %s", string(body))
	}

	return string(body), nil
}

func Trade() {
	symbol := "SOLUSDT"  // Торговая пара
	side := "BUY"        // Покупка
	orderType := "LIMIT" // Тип ордера
	quantity := "0.027"  // Количество
	price := "200"       // Цена (для лимитного ордера)

	response, err := createOrder(symbol, side, orderType, quantity, price)
	if err != nil {
		fmt.Printf("Error creating order: %v\n", err)
		return
	}

	fmt.Printf("Order created: %s\n", response)
}
