package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	types "tg_bot/helpers"
)

func SendCommandList(chatId int, commands map[string]string, baseUrl, botToken, telegramMethod string) error {
	buttons := [][]types.KeyboardButton{}

	for command := range commands {
		buttonRow := []types.KeyboardButton{
			{Text: command},
		}

		buttons = append(buttons, buttonRow)
	}

	replyKeyboard := types.ReplyKeyboardMarkup{
		Keyboard:        buttons,
		OneTimeKeyboard: true,
	}

	requestBody := types.BotMessage{
		ChatId:      chatId,
		Text:        "Commands list:",
		ReplyMarkup: replyKeyboard,
	}

	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Can't encode request body: ", err)
		return err
	}

	resp, err := http.Post(baseUrl+botToken+"/"+telegramMethod, "application/json", bytes.NewReader(requestBodyJson))
	if err != nil {
		fmt.Println("Can't set menu buttons: ", err)
		return err
	}

	resp.Body.Close()

	return nil
}

func SendTextMessage(message types.Message, text, baseUrl, botToken, telegramMethod string) error {
	queryParams := url.Values{}
	queryParams.Add("chat_id", strconv.FormatInt(int64(message.Chat.Id), 10))
	queryParams.Add("text", text)

	urlStr := baseUrl + botToken + "/" + telegramMethod

	resp, err := http.PostForm(urlStr, queryParams)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return err
	}

	resp.Body.Close()

	return nil
}

func SendVoiceMessage(message types.Message, basePathToLoadMp3File, baseUrl, botToken, telegramMethod string) error {
	voiceFilePath := filepath.FromSlash(fmt.Sprintf("%s/%s.mp3", basePathToLoadMp3File, strconv.Itoa(message.MessageId)))

	// Открываем голосовой файл
	file, err := os.Open(voiceFilePath)
	if err != nil {
		fmt.Println("Can't find path to voice file: ", err)
		return err
	}
	defer file.Close()

	// Создаем multipart/form-data для отправки файла
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("voice", filepath.Base(voiceFilePath))
	if err != nil {
		fmt.Println("Can't create multipart/form-data: ", err)
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Can't copy multipart/form-data: ", err)
		return err
	}
	writer.Close()

	// Создаем POST-запрос к API телеграма для отправки голосового сообщения
	apiURL := baseUrl + botToken + "/" + telegramMethod

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		fmt.Println("Can't request to send voice: ", err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Добавляем параметры к телу запроса
	querys := req.URL.Query()
	querys.Add("chat_id", strconv.Itoa(message.Chat.Id))
	querys.Add("reply_to_message_id", strconv.Itoa(message.MessageId))
	req.URL.RawQuery = querys.Encode()

	// Отправляеме запрос и получаем ответ
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Can't get response: ", err)
		return err
	}
	defer resp.Body.Close()

	// Выводим статус отправки сообщения
	fmt.Println("Voice message send")

	return nil
}

/*
*	TODO надо сделать переключалку в командах в каком виде хочет пользователь получать ответы - аудио, войс или текст
*	пока метод неиспользуется
 */
func SendAudioMessage(message types.Message, basePathToLoadMp3File, baseUrl, botToken, telegramMethod string) error {
	audioFilePath := filepath.FromSlash(fmt.Sprintf("%s/%s.mp3", basePathToLoadMp3File, strconv.Itoa(message.MessageId)))

	// Открываем аудиофайл
	file, err := os.Open(audioFilePath)
	if err != nil {
		fmt.Println("Can't find path to voice file: ", err)
		return err
	}
	defer file.Close()

	// Создаем multipart/form-data для отправки файла
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audio", filepath.Base(audioFilePath))
	if err != nil {
		fmt.Println("Can't create multipart/form-data: ", err)
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Can't copy multipart/form-data: ", err)
		return err
	}
	writer.Close()

	// Создаем POST-запрос к API телеграма для отправки аудиофайла
	apiURL := baseUrl + botToken + "/" + telegramMethod

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		fmt.Println("Can't request to send audio: ", err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Добавляем параметры к телу запроса
	querys := req.URL.Query()
	querys.Add("chat_id", strconv.Itoa(message.Chat.Id))
	querys.Add("reply_to_message_id", strconv.Itoa(message.MessageId))
	req.URL.RawQuery = querys.Encode()

	// Отправляем запрос и получаем ответ
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Can't get response: ", err)
		return err
	}
	defer resp.Body.Close()

	// Выводим статус отправки сообщения
	fmt.Println("Audio message send")

	return nil
}
