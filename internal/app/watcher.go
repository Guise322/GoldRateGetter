package app

import (
	lTime "PriceWatcher/internal/app/time"
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func (s PriceWatcherService) Watch(done <-chan struct{}, cancel context.CancelFunc, clock lTime.Clock) {
	WatchForInterruption(cancel)

	errMes := "An error occurs while serving a price: %v"
	if err := s.serve(clock); err != nil {
		logrus.Errorf(errMes, err)
	}

	logrus.Info("The price is processed")

	config, err := s.conf.GetConfig()
	if err != nil {
		logrus.Errorf("An error occurs while get the config data: %v", err)

		return
	}

	if strings.ToLower(config.PriceType) == "bank" {
		err := waitHourStart(clock.Now())
		if err != nil {
			logrus.Errorf("An error occurs while waiting when the next hour begins: %v", err)

			return
		}

		if err := s.serve(clock); err != nil {
			logrus.Errorf(errMes, err)
		}

		logrus.Info("The price is processed")
	}

	var dur time.Duration

	if strings.ToLower(config.PriceType) == "marketplace" {
		dur = time.Duration(20 * time.Minute)
	} else {
		dur = time.Duration(1 * time.Hour)
	}

	t := time.NewTicker(dur)
	defer t.Stop()

	for {
		select {
		case <-done:
			logrus.Info("Shut down the application")
			return
		case <-t.C:
			if err := s.serve(clock); err != nil {
				logrus.Errorf(errMes, err)
			}

			logrus.Info("The price is processed")

			if strings.ToLower(config.PriceType) == "marketplace" {
				dur = time.Duration(20+rand.Intn(10)) * time.Minute
				t.Reset(dur)
			}
		}
	}
}
