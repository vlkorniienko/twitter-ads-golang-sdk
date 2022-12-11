package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/civil"
	"github.com/rs/zerolog"

	"github.com/vlkorniienko/twitter_ads/client"
)

type worker struct {
	c *client.API
	l *zerolog.Logger
}

func main() {
	c, err := client.NewAPIClient(makeConfig())
	if err != nil {
		log.Fatal(err)
	}

	w := worker{
		c: c,
		l: makeLogger(),
	}

	w.process()
}

func (w worker) process() {
	for _, v := range w.c.AdAccounts {
		w.l.Info().Str("account", v.AdAccountName).Msg("start process account")

		err := w.processAccount(v.AdAccountID, v.AdAccountName)
		if err != nil {
			w.l.Info().Str("account", v.AdAccountName).Msg("can't process account")
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
		w.l.Info().Str("account", accountName).Msg("no data for selected account, skipping")

		return nil
	}

	for _, v := range activeEntities.Data {
		stats, err := w.c.GetSpendStats(accountID, v.EntityID, startTime, endTime)
		if err != nil {
			return fmt.Errorf("can't get spend stats from twitter api: %w", err)
		}

		if len(stats.Data) == 0 || len(stats.Data[0].IDData) == 0 ||
			len(stats.Data[0].IDData[0].Metrics.BilledChargeLocalMicro) == 0 {
			w.l.Info().Str("account", accountName).Msg("no stats data for selected account, skipping")

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

		w.l.Info().Interface("spend entity", spend).Msg("processed")
	}

	return nil
}

func makeConfig() *client.Config {
	conf := &client.Config{
		APIKey:       os.Getenv("TWITTER_API_KEY"),
		APISecret:    os.Getenv("TWITTER_API_SECRET"),
		AccessToken:  os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessSecret: os.Getenv("TWITTER_ACCESS_SECRET"),
		AdAccounts: []client.TwitterAcc{
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
