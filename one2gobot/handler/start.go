package handler

import (
	"greenisha/one2gobot/client"
	"greenisha/one2gobot/store"
	"os"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func NewHandler(b *gotgbot.Bot, ctx *ext.Context, s store.Store) Handler {
	return Handler{Bot: b, Context: ctx, Rest: client.Client{RestEndpoint: os.Getenv("RESTAPI")}, S: s}
}

type Handler struct {
	Rest    client.Client
	Bot     *gotgbot.Bot
	Context *ext.Context
	S       store.Store
}

func (h Handler) SendStart() error {

	userId := h.Context.EffectiveMessage.From.Id
	h.S.Set(store.Session{UserId: userId, ChatId: h.Context.EffectiveChat.Id, Mode: store.SelectFrom})
	h.sendStatus(userId)
	return nil
}

func (h Handler) sendStatus(userId int64) error {
	data, _ := h.S.Get(userId)
	var sb strings.Builder
	sb.WriteString("I'll help you order 🚌 🚆 tickets. \n")
	if data.Mode == store.SelectFrom {
		sb.WriteString("Please type departure station, we'll look for it")
	}
	if data.Mode == store.SelectTo {
		sb.WriteString("Departure station is ")
		sb.WriteString(data.StationFrom.Name + " (" + data.StationFrom.Country + ") \n")
		sb.WriteString("Please type destination station, we'll look for it")
	}
	if data.Mode == store.Finished {
		sb.WriteString("Departure station is ")
		sb.WriteString(data.StationFrom.Name + " (" + data.StationFrom.Country + ") \n")
		sb.WriteString("Arrival station is ")
		sb.WriteString(data.StationTo.Name + " (" + data.StationTo.Country + ") \n")

		sb.WriteString("Here is link <a href=\"https://12go.asia/en/travel/'" + data.StationFrom.Slug + "/" + data.StationTo.Slug + "?z=5640093\">Press here</a> \n")
	}
	h.Context.EffectiveMessage.Reply(h.Bot, sb.String(), &gotgbot.SendMessageOpts{ParseMode: "html"})
	return nil
}
