# twitter-ads-golang-sdk

At this moment there is no official SDK for Twitter Ads API for Go programming language [sdk tools and libraries](https://developer.twitter.com/en/docs/twitter-ads-api/tools-and-libraries). 
The purpose of this library is to fetch spend statistics by active twitters campaign. The process of creating a signature for making requests to the Twitter API is a little tricky and can be confusing due to the large number of steps that must be taken to create a signature for each request [creating a signature](https://developer.twitter.com/en/docs/authentication/oauth-1-0a/creating-a-signature). That's why I decided to share this tool with usage example to save your time. Please feel free to use it and improve it. 

## Example

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/civil"
	"github.com/rs/zerolog"

	"github.com/vlkorniienko/twitter_ads"
)

type worker struct {
	c *twitter_ads.API
	logger *zerolog.Logger
}

func main() {
	c, err := twitter_ads.NewAPIClient(makeConfig())
	if err != nil {
		log.Fatal(err)
	}

	w := worker{
		c: c,
		logger: makeLogger(),
	}

	w.process()
}

func (w worker) process() {
	for _, v := range w.c.AdAccounts {
		w.logger.Info().Str("account", v.AdAccountName).Msg("start process account")

		err := w.processAccount(v.AdAccountID, v.AdAccountName)
		if err != nil {
			w.logger.Info().Str("account", v.AdAccountName).Msg("can't process account")
		}
	}
}

type Spend struct {
	Date      civil.Date
	AdAccount string
	Campaign  string
	Spend     float64
	Currency  string
}

func (w worker) processAccount(accountID, accountName string) error {
	timeToWrite := time.Now().AddDate(0, 0, -1)
	startTime := timeToWrite.Format("2006-01-02")
	endTime := time.Now().Format("2006-01-02")

	activeEntities, err := w.c.GetActiveEntities(accountID, startTime, endTime)
	if err != nil {
		return fmt.Errorf("can't get active entities from twitter api: %w", err)
	}

	if len(activeEntities.Data) == 0 {
		w.logger.Info().Str("account", accountName).Msg("no data for selected account, skipping")

		return nil
	}

	for _, v := range activeEntities.Data {
		stats, err := w.c.GetSpendStats(accountID, v.EntityID, startTime, endTime)
		if err != nil {
			return fmt.Errorf("can't get spend stats from twitter api: %w", err)
		}

		if len(stats.Data) == 0 || len(stats.Data[0].IDData) == 0 ||
			len(stats.Data[0].IDData[0].Metrics.BilledChargeLocalMicro) == 0 {
			w.logger.Info().Str("account", accountName).Msg("no stats data for selected account, skipping")

			continue
		}

		campaign, err := w.c.GetCampaignInfo(accountID, v.EntityID)
		if err != nil {
			return fmt.Errorf("can't get campaign info from twitter api: %w", err)
		}

		spend := Spend{
			Campaign:  campaign.Data.Name,
			Date:      civil.DateOf(timeToWrite),
			AdAccount: accountName,
			Spend:     float64(stats.Data[0].IDData[0].Metrics.BilledChargeLocalMicro[0]) / 1000000,
			Currency:  campaign.Data.Currency,
		}

		w.logger.Info().Interface("spend entity", spend).Msg("processed")
	}

	return nil
}

func makeConfig() *twitter_ads.Config {
	conf := &twitter_ads.Config{
		APIKey:       os.Getenv("TWITTER_API_KEY"),
		APISecret:    os.Getenv("TWITTER_API_SECRET"),
		AccessToken:  os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessSecret: os.Getenv("TWITTER_ACCESS_SECRET"),
		AdAccounts: []twitter_ads.TwitterAccount{
			{
				AdAccountName: os.Getenv("TWITTER_AD_ACCOUNT_NAME"),
				AdAccountID:   os.Getenv("TWITTER_AD_ACCOUNT_ID"),
			},
		},
	}

	return conf
}

func makeLogger() *zerolog.Logger {
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	logger := zerolog.New(os.Stdout).
		Level(zerolog.DebugLevel).
		With().Timestamp().
		Logger()

	return &logger
}
```