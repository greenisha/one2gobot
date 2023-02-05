package handler

import (
	"encoding/json"
	"errors"
	"greenisha/one2gobot/model"
	"greenisha/one2gobot/store"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func (h Handler) Callback() error {
	cb := h.Context.Update.CallbackQuery
	var callbackData model.CallbackData
	log.Println(cb.Data)
	err := json.Unmarshal([]byte(cb.Data), &callbackData)
	if err != nil {
		log.Println(err.Error())
	}
	session, found := h.S.Get(cb.From.Id)
	if !found {
		station, err := h.Rest.FindStationBySlug(callbackData.Slug)
		if err != nil {
			return err
		}
		session := store.Session{UserId: cb.From.Id, StationFrom: station, Mode: store.SelectTo}
		h.S.Set(session)

	} else {
		if session.Mode == store.SelectFrom || session.Mode == store.Finished {
			station, err := h.Rest.FindStationBySlug(callbackData.Slug)
			if err != nil {
				return err
			}
			session.StationFrom = station
			session.Mode = store.SelectTo
			h.S.Set(session)

		} else if session.Mode == store.SelectTo {
			if session.StationFrom.Slug == callbackData.Slug {
				cb.Answer(h.Bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Please select station that is not the same as departing"})
				return errors.New("Same station")
			}
			station, err := h.Rest.FindStationBySlug(callbackData.Slug)
			if err != nil {
				return err
			}

			session.StationTo = station
			session.Mode = store.Finished
			h.S.Set(session)

		}
	}
	h.sendStatus(cb.From.Id)

	return nil
}
