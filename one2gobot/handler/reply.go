package handler

import (
	"encoding/json"
	"fmt"
	"greenisha/one2gobot/model"
	"log"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func (h Handler) SendReply(s string) error {
	input := s

	stations, err := h.Rest.FindStation(input)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("failed to find stations: %w", err)
	}
	h.Context.EffectiveMessage.Reply(h.Bot, "Found this stations", &gotgbot.SendMessageOpts{
		ParseMode:   "html",
		ReplyMarkup: h.constructStations(stations),
	})
	return nil

}

func (h Handler) constructStations(s []model.Station) gotgbot.InlineKeyboardMarkup {
	var buttons [][]gotgbot.InlineKeyboardButton
	rowCount := 0
	var buttonsRow []gotgbot.InlineKeyboardButton
	//maxCount := 7
	for _, value := range s {
		rowCount++
		callbackData := model.CallbackData{Slug: value.Slug}
		jsonData, _ := json.Marshal(callbackData)
		button := gotgbot.InlineKeyboardButton{Text: constructButtonText(value), CallbackData: string(jsonData)}
		buttonsRow = append(buttonsRow, button)
		if rowCount >= 2 {
			buttons = append(buttons, buttonsRow)
			buttonsRow = nil
			rowCount = 0
		}
	}
	if buttonsRow != nil {
		buttons = append(buttons, buttonsRow)
	}
	var numericKeyboard = gotgbot.InlineKeyboardMarkup{InlineKeyboard: buttons}
	return numericKeyboard

}

func constructButtonText(s model.Station) string {
	var sb strings.Builder
	sb.WriteString(s.Name)
	sb.WriteString(" ")
	sb.WriteString(s.Country)
	return sb.String()
}
