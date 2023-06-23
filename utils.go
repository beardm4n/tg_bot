package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"tg_bot/file"
	types "tg_bot/helpers"
	messages "tg_bot/message"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/joho/godotenv"
)

func initEnv() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Something went wrong in get env variables")
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

func getUpdates(baseUrl string, botToken string, method string, offset int) ([]types.Update, error) {
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

	var restResponse types.MessageResponse

	// необходим распарсить json, который получили от сервера, который приведем к структуре RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		fmt.Println("Something went wrong in parse json in getUpdates: ", err)
		return nil, err
	}

	return restResponse.Result, nil
}

func processUpdate(update types.Update) {
	message := update.Message
	command := strings.TrimSpace(message.Text)

	// получение и обработка голосового сообщения
	if len(message.Voice.FileId) != 0 {
		file.DownloadFile(baseUrl, botToken, telegramMethods["GET_FILE"], fileUrl, basePathToSaveFile, message)
		fmt.Println("Download file completed")

		return
	}

	if strings.HasPrefix(message.Text, "/") {
		if command == "/commands" {
			messages.SendCommandList(message.Chat.Id, commands, baseUrl, botToken, telegramMethods["SEND_MESSAGE"])
		} else {
			if description, ok := commands[command]; ok {
				sendMessage(message, description)
			} else {
				sendMessage(message, "Send /commands to get a list of commands.")
			}
		}
	} else {
		textToTalk(update.Message)
		convertMp3ToOga(update.Message)
		sendMessage(update.Message, "")
	}
}

func sendMessage(message types.Message, text string) error {
	if text != "" {
		messages.SendTextMessage(message, text, baseUrl, botToken, telegramMethods["SEND_MESSAGE"])
		return nil
	}

	messages.SendVoiceMessage(message, basePathToLoadMp3File, baseUrl, botToken, telegramMethods["SEND_VOICE"])

	return nil
}

func textToTalk(message types.Message) {
	speech := htgotts.Speech{
		Folder:   folderAudioName,
		Language: message.From.LanguageCode,
	}

	speech.CreateSpeechFile(message.Text, strconv.Itoa(message.MessageId))
}

func convertMp3ToOga(message types.Message) error {
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

	return nil
}
