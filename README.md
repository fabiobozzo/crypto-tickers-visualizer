# Cryptocurrencies tickers visualizer

Small app written in Go to monitor a configurable set of crypto-pairs. 
It uses InfluxDB as TSDB and CoinPaprika to scrape the quotes in realtime.

### Requirements
- Docker
- Docker Compose

### Usage
- edit `tickers.json` by adding/removing coins (or tokens) and specifying which quotes you're interested in
- `make start` to start InfluxDB + Grafana + the quotes-fetcher app itself, as a daemon
- `make stop` to destroy everything 

