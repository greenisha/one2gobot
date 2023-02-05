package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"greenisha/one2gobot/handler"
	"greenisha/one2gobot/store"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("API_KEY")
	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}
	// Create updater and dispatcher.
	updater := ext.NewUpdater(nil)

	dispatcher := updater.Dispatcher
	repository := NewRepository()
	// /start command to introduce the bot
	dispatcher.AddHandler(handlers.NewCommand("start", repository.Start))
	dispatcher.AddHandler(handlers.NewMessage(message.Text, repository.Reply))
	dispatcher.AddHandler(handlers.NewCallback(nil, repository.Callback))
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	fmt.Printf("%s has been started...\n", b.User.Username)

	updater.Idle()

}

type Repository struct {
	s store.Store
}

func NewRepository() Repository {
	return Repository{s: store.NewMemoryStore()}
}

func (s *Repository) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	return handler.NewHandler(b, ctx, s.s).SendStart()
}

func (s *Repository) Reply(b *gotgbot.Bot, ctx *ext.Context) error {
	return handler.NewHandler(b, ctx, s.s).SendReply(ctx.EffectiveMessage.Text)
}

func (s *Repository) Callback(b *gotgbot.Bot, ctx *ext.Context) error {
	return handler.NewHandler(b, ctx, s.s).Callback()
}

// func search(b *gotgbot.Bot, ctx *ext.Context) error {
// 	re := regexp.MustCompile(`/start(@/w+)? (.+)`)

// 	if re.FindStringSubmatch(ctx.EffectiveMessage.Text) != nil {
// 		return handler.NewHandler(b, ctx).SendSearch(re.FindStringSubmatch(ctx.EffectiveMessage.Text)[2])
// 	}
// 	return nil
// }
