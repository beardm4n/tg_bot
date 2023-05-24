package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	fileUrl = os.Getenv("FILE_URL")
	basePathToSaveFile = os.Getenv("BASE_PATH_TO_SAVE_FILE")

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

func downloadFile(message Message) {
	// перенести код ниже сюда
	// посмотреть как можно унифицировать слэши в пути для разных ОС
	file := getFile(message.Voice.FileId)

	downloadFileUrl := fileUrl + botToken + "/" + file.FilePath

	// Выполняем запрос для получения содержимого файла
	fileResp, err := http.Get(downloadFileUrl)
	if err != nil {
		fmt.Println("Something went wrong in error handling in download file")
		return
	}

	defer fileResp.Body.Close()

	filePathOnDisk := filepath.FromSlash(fmt.Sprintf("%s/voice_%s.oga", basePathToSaveFile, message.Voice.FileId))

	// Создаем файл на диске для сохранения содержимого
	fileOnDisk, err := os.Create(filePathOnDisk)
	if err != nil {
		fmt.Println("Something went wrong in create file into disk")
		return
	}

	defer fileOnDisk.Close()

	// Копируем полученный файл в файл на диске
	_, err = io.Copy(fileOnDisk, fileResp.Body)
	if err != nil {
		fmt.Println("Something went wrong when save file into disk")
		return
	}
}

func processUpdate(update Update) {
	message := update.Message
	command := strings.TrimSpace(message.Text)

	// получение и обработка голосового сообщения
	if len(message.Voice.FileId) != 0 {
		downloadFile(message)
		return
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
	queryParams := url.Values{}
	queryParams.Add("chat_id", strconv.FormatInt(int64(chatId), 10))
	queryParams.Add("text", text)

	urlStr := baseUrl + botToken + "/" + telegramMethods["SEND_MESSAGE"]

	_, err := http.PostForm(urlStr, queryParams)
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
