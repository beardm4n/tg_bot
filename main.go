package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	botToken          = "<somekey>"
	botApi            = "https://api.telegram.org/bot"
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func main() {
	requestBorUrl := botApi + botToken
	offset := 0

	for {
		updates, err := getUpdates(requestBorUrl, offset)
		if err != nil {
			log.Println("Something went wrong in getUpdates", err.Error())
		}

		for _, update := range updates {
			err := respond(requestBorUrl, update)
			if err != nil {
				log.Println("Something went wrong in update", err.Error())
			}
			offset = update.UpdateId + 1
		}

		fmt.Println(updates)
	}
}

// получение обновлений
func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/" + getUpdatesMethod + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}

	// закрываем запрос после того как выйдем из функции
	defer resp.Body.Close()

	// ответ от сервера получаем в байтах, необходимо обработать его
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse RestResponse

	// необходим распарсить json, который получили от сервера, который приведем к структуре RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

// ответ на обновления
func respond(botUrl string, update Update) error {
	var botMessage BotMessage

	botMessage.ChatId = update.Message.Chat.Id
	botMessage.Text = update.Message.Text

	// тело запроса передается в байтовом формате
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}

	url := botUrl + "/" + sendMessageMethod

	// для того, чтобы передать тело запроса нужен Reader, при помощи bytes.NewBuffer создаем его
	_, err = http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}
