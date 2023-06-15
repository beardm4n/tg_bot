package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	htgotts "github.com/hegedustibor/htgo-tts"
	voices "github.com/hegedustibor/htgo-tts/voices"
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
	basePathToLoadMp3File = os.Getenv("BASE_PATH_TO_LOAD_MP3_FILE")
	basePathToSaveOgaFile = os.Getenv("BASE_PATH_TO_SAVE_OGA_FILE")

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
		textToTalk(update.Message)
		convertMp3ToOga(update.Message)
		sendVoiceMessage(update.Message)
	}
}

func sendMessage(chatId int, text string) error {
	queryParams := url.Values{}
	queryParams.Add("chat_id", strconv.FormatInt(int64(chatId), 10))
	queryParams.Add("text", text)

	urlStr := baseUrl + botToken + "/" + telegramMethods["SEND_MESSAGE"]

	resp, err := http.PostForm(urlStr, queryParams)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return err
	}

	resp.Body.Close()

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

func textToTalk(message Message) {
	speech := htgotts.Speech{
		Folder:   folderAudioName,
		Language: voices.Russian,
	}

	speech.CreateSpeechFile(message.Text, strconv.Itoa(message.MessageId))
}

func convertMp3ToOga(message Message) error {
	// Конвертация в OGG
	cmd := exec.Command("python", "scripts/converter.py", strconv.Itoa(message.MessageId), basePathToLoadMp3File, basePathToSaveOgaFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Python script exited with an error: %v\n%s", exitErr, output)
		} else {
			fmt.Printf("Error while executing python script: %v", err)
		}
	}

	fmt.Println("Python script completed successfully!")
	fmt.Println(string(output))

	return nil
}

func sendVoiceMessage(message Message) error {
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
	apiURL := baseUrl + botToken + "/" + telegramMethods["SEND_VOICE"] + "?chat_id=" + strconv.Itoa(message.Chat.Id)

	fmt.Println("apiURL", apiURL)

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		fmt.Println("Can't request to send voice: ", err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Отправляеме запрос и получаем ответ
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Выводим статус отправки сообщения
	log.Println("Voice message send")

	return nil
}

/*
*	TODO надо сделать переключалку в командах в каком виде хочет пользователь получать ответы - аудио, войс или текст
*	пока метод неиспользуется
 */
func sendAudioMessage(message Message) error {
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
	apiURL := baseUrl + botToken + "/" + telegramMethods["SEND_AUDIO"] + "?chat_id=" + strconv.Itoa(message.Chat.Id)

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		fmt.Println("Can't request to send voice: ", err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Отправляем запрос и получаем ответ
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Can't get response: ", err)
		return err
	}
	defer resp.Body.Close()

	// Читаем ответ
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Can't read response: ", err)
		return err
	}

	// Выводим ответ
	log.Println("RESPONSE", string(responseBody))

	return nil
}
