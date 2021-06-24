package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/coinpaprika/coinpaprika-api-go-client/coinpaprika"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	panicErr(err)

	tickersConfig, err := readConfigFromFile("tickers.json")
	panicErr(err)

	coinpaprikaClient := coinpaprika.NewClient(nil)

	coinIDsMap, err := mapCoinSymbolToID(coinpaprikaClient)
	panicErr(err)

	influxdbClient := influxdb2.NewClient(os.Getenv("INFLUXDB_HOST"), os.Getenv("INFLUXDB_TOKEN"))
	influxdbWriteAPI := influxdbClient.WriteAPIBlocking(os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))
	defer influxdbClient.Close()

	for {
		for _, t := range tickersConfig {
			coinID, found := coinIDsMap[t.Coin]
			if !found {
				fmt.Printf("Coin ID not found for %s\n", t.Coin)

				continue
			}

			ticker, err := coinpaprikaClient.Tickers.GetByID(coinID, &coinpaprika.TickersOptions{
				Quotes: strings.Join(t.Quotes[:], ","),
			})
			panicErr(err)

			for quoteSymbol, quote := range ticker.Quotes {
				line := fmt.Sprintf("tickers,coin=%s,quote=%s price=%f", t.Coin, quoteSymbol, *quote.Price)

				if err := influxdbWriteAPI.WriteRecord(context.Background(), line); err != nil {
					fmt.Println(err)
				}
			}
		}

		time.Sleep(60 * time.Second)
	}
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func readConfigFromFile(filename string) ([]Ticker, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config Tickers

	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}

	return config.Tickers, nil
}

func mapCoinSymbolToID(client *coinpaprika.Client) (map[string]string, error) {
	coinSymbolToIDs := map[string]string{}

	coins, err := client.Coins.List()
	if err != nil {
		return nil, err
	}

	for _, c := range coins {
		coinSymbolToIDs[*c.Symbol] = *c.ID
	}

	return coinSymbolToIDs, nil
}
