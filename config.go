package main

type Tickers struct {
	Tickers []Ticker `json:"tickers"`
}

type Ticker struct {
	Coin   string   `json:"coin"`
	Quotes []string `json:"quotes"`
}
