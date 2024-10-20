package telebot

import (
	"PriceWatcher/internal/entities/bot"
	"PriceWatcher/internal/interfaces/broker"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context,
	wg *sync.WaitGroup,
	broker broker.Worker,
	commands []bot.Command) error {
	msgChan, err := broker.Start(ctx)
	if err != nil {
		return err
	}

	go processeMessages(ctx, wg, broker, msgChan, commands)

	go func() {
		<-ctx.Done()
		if err := broker.Stop(); err != nil {
			logrus.Errorln("An error occurs at stopping the message worker: ", err)
		}
		wg.Done()
	}()

	return nil
}

func processeMessages(ctx context.Context,
	wg *sync.WaitGroup,
	broker broker.Worker,
	msgChan <-chan bot.Message,
	commands []bot.Command) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
		case msg := <-msgChan:
			for _, cmd := range commands {
				if cmd.Name == msg.Command {
					result := cmd.Action(msg)
					broker.SendMessage(ctx, result, msg.ChatID)
				}
			}
		}
	}
}
