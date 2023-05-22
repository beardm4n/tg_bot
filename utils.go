package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func initEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	botToken = os.Getenv("BOT_TOKEN")
	baseUrl = os.Getenv("BASE_URL")

	return nil
}

func getUpdates(baseUrl string, botToken string, method string, offset int) ([]Update, error) {
	resp, err := http.Get(baseUrl + botToken + "/" + method + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		fmt.Println("Something went wrong in requset: ", err)
		return nil, err
	}

	// ответ от сервера получаем в байтах, необходимо обработать его
	body, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		if err != nil {
			fmt.Println("Something went wrong in error handling: ", err)
			return nil, err
		}
	}

	var restResponse RestResponse

	// необходим распарсить json, который получили от сервера, который приведем к структуре RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		fmt.Println("Something went wrong in parse json: ", err)
		return nil, err
	}

	return restResponse.Result, nil
}

func checkSentMessage(message Message) {
	var msg BotMessage

	if message.Text == "/start" {
		msg.ChatId = message.Chat.Id
		msg.Text = "Hi, nice to meet you!"

		sendMessage(baseUrl, botToken, telegramMethods["SEND_MESSAGE"], msg)
	} else if message.Text == "/stop" {
		msg.ChatId = message.Chat.Id
		msg.Text = "Bye bye! see you soon"

		sendMessage(baseUrl, botToken, telegramMethods["SEND_MESSAGE"], msg)
	} else if strings.Contains(message.Text, "/") {
		msg.ChatId = message.Chat.Id
		msg.Text = "Unknow command =/"

		sendMessage(baseUrl, botToken, telegramMethods["SEND_MESSAGE"], msg)
	} else if !strings.Contains(message.Text, "/") {
		msg.ChatId = message.Chat.Id
		msg.Text = message.Text

		sendMessage(baseUrl, botToken, telegramMethods["SEND_MESSAGE"], msg)
	}
}

func sendMessage(baseUrl string, botToken string, method string, message BotMessage) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(message.ChatId))
	params.Set("text", message.Text)

	resp, err := http.PostForm(baseUrl+botToken+"/"+method, params)
	if err != nil {
		fmt.Println("Can't send the message: ", err)
		return err
	}

	resp.Body.Close()
	resp.StatusCode = 150

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error HTTP-request. Status Code: %d", resp.StatusCode)
	}

	return nil
}
