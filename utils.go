package main

import (
	"bytes"
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
			fmt.Println("Something went wrong in error handling in getUpdates: ", err)
			return nil, err
		}
	}

	var restResponse MessageResponse

	// необходим распарсить json, который получили от сервера, который приведем к структуре RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		fmt.Println("Something went wrong in parse json in getUpdates: ", err)
		return nil, err
	}

	return restResponse.Result, nil
}

func getFile(fileId string) File {
	urlGetFile := baseUrl + botToken + "/" + telegramMethods["GET_FILE"] + "?file_id=" + fileId

	resp, err := http.Get(urlGetFile)
	if err != nil {
		fmt.Println("Something went wrong in error handling in getFile: ", err)
	}

	// ответ от сервера получаем в байтах, необходимо обработать его
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if err != nil {
			fmt.Println("Something went wrong in error handling in getUpdates: ", err)
		}
	}

	defer resp.Body.Close()

	var restResponse VoiceResponse

	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		fmt.Println("Something went wrong in parse json in getFile: ", err)
	}

	return restResponse.Result
}

func processUpdate(update Update) {
	message := update.Message
	command := strings.TrimSpace(message.Text)

	if len(message.Voice.FileId) != 0 {
		getFile(message.Voice.FileId)
	}

	if strings.HasPrefix(message.Text, "/") {
		if command == "/commands" {
			sendCommandList(message.Chat.Id)
		} else {
			if description, ok := commands[command]; ok {
				sendMessage(message.Chat.Id, description)
			} else {
				sendMessage(message.Chat.Id, "Unknown coomand.")
			}
		}
	} else {
		sendMessage(update.Message.Chat.Id, update.Message.Text)
	}
}

func sendMessage(chatId int, text string) error {
	data := url.Values{}
	data.Set("chat_id", strconv.FormatInt(int64(chatId), 10))
	data.Set("text", text)

	urlStr := baseUrl + botToken + "/" + telegramMethods["SEND_MESSAGE"]

	_, err := http.PostForm(urlStr, data)
	if err != nil {
		fmt.Println("Failed to send message:", err)
	}

	return nil
}

func sendCommandList(chatId int) error {
	buttons := [][]KeyboardButton{}

	for command := range commands {
		buttonRow := []KeyboardButton{
			{Text: command},
		}

		buttons = append(buttons, buttonRow)
	}

	replyKeyboard := ReplyKeyboardMarkup{
		Keyboard:        buttons,
		OneTimeKeyboard: true,
	}

	requestBody := BotMessage{
		ChatId:      chatId,
		Text:        "Commands list:",
		ReplyMarkup: replyKeyboard,
	}

	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Can't encode request body: ", err)
		return err
	}

	resp, err := http.Post(baseUrl+botToken+"/"+telegramMethods["SEND_MESSAGE"], "application/json", bytes.NewReader(requestBodyJson))
	if err != nil {
		fmt.Println("Can't set menu buttons: ", err)
		return err
	}

	resp.Body.Close()

	return nil
}
